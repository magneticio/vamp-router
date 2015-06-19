package metrics

import (
	"github.com/magneticio/vamp-router/haproxy"
	gologger "github.com/op/go-logging"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type Streamer struct {

	// simple counter to give heartbeats in the log how many messages where parsed during a time period
	Counter       int64
	wantedMetrics []string
	haRuntime     *haproxy.Runtime
	pollFrequency int
	Clients       map[chan Metric]bool
	Log           *gologger.Logger
}

// Adds a client to which messages can be multiplexed.
func (s *Streamer) AddClient(c chan Metric) {
	s.Clients[c] = true
}

// Just sets the metrics we want for now...
func (s *Streamer) Init(haRuntime *haproxy.Runtime, frequency int, log *gologger.Logger) {
	s.Log = log
	s.wantedMetrics = []string{"scur", "qcur", "qmax", "smax", "slim", "ereq", "econ", "lastsess", "qtime", "ctime", "rtime", "ttime", "req_rate", "req_rate_max", "req_tot", "rate", "rate_lim", "rate_max", "hrsp_1xx", "hrsp_2xx", "hrsp_3xx", "hrsp_4xx", "hrsp_5xx"}
	s.haRuntime = haRuntime
	s.pollFrequency = frequency
	s.Clients = make(map[chan Metric]bool)
}

/*
  Generates an outgoing stream of discrete Metric struct values.
  This stream can then be consumed by other streams like Kafka or SSE.
*/
func (s *Streamer) Start() error {

	s.Heartbeat()

	// create a channel to send the stats to the parser
	statsChannel := make(chan map[string]map[string]string)

	// start up the parser in a separate routine
	go ParseMetrics(statsChannel, s.Clients, s.wantedMetrics, &s.Counter)

	for {
		// start pumping the stats into the channel
		stats, _ := s.haRuntime.GetStats("all")
		statsChannel <- stats
		time.Sleep(time.Duration(s.pollFrequency) * time.Millisecond)
	}
}

/*
	Parses a []Stats and injects it into each Metric channel in a map of channels
*/

func ParseMetrics(statsChannel chan map[string]map[string]string, c map[chan Metric]bool, wantedMetrics []string, counter *int64) {

	wantedFrontendMetric := make(map[string]bool)
	wantedFrontendMetric["ereq"] = true
	wantedFrontendMetric["rate_lim"] = true
	wantedFrontendMetric["req_rate_max"] = true
	wantedFrontendMetric["req_rate"] = true

	for {

		select {

		case stats := <-statsChannel:

			localTime := time.Now().Format(time.RFC3339)

			// for each proxy in the stats dump, pick out the wanted metrics.
			for _, proxy := range stats {

				// loop over all wanted metrics for the current proxy
				for _, metric := range wantedMetrics {

					// discard all empty metrics
					if proxy[metric] != "" {

						value := proxy[metric]
						svname := proxy["svname"]
						tags := []string{}
						pxnames := strings.Split(proxy["pxname"], "::")

						// allow only some FRONTEND metrics and all non-FRONTEND metrics
						if (svname == "FRONTEND" && wantedFrontendMetric[metric]) || svname != "FRONTEND" {

							// Compile tags
							// we tag the metrics according to the following scheme
							switch {

							//- if pxname has no "." separator, and svname is [BACKEND|FRONTEND] it is the top route or "endpoint"
							case len(pxnames) == 1 && (svname == "BACKEND" || svname == "FRONTEND"):
								tags = append(tags, "routes:"+proxy["pxname"], "route")

								EmitMetric(localTime, tags, metric, value, counter, c)

							//-if pxname has no "."  separator, and svname is not [BACKEND|FRONTEND] it is an "in between"
							// server that routes to the actual service via a socket.
							case len(pxnames) == 1 && (svname != "BACKEND" || svname != "FRONTEND"):
							// sockName := strings.Split(svname, ".")
							// tags = append(tags, "routes:"+proxy["pxname"], "socket_servers:"+sockName[1])

							// we dont emit this metrics currently
							// EmitMetric(localTime, tags, metric, value, counter, c)

							//- if pxname has a separator, and svname is [BACKEND|FRONTEND] it is a service
							case len(pxnames) > 1 && (svname == "BACKEND" || svname == "FRONTEND"):
								tags = append(tags, "routes:"+pxnames[0], "services:"+pxnames[1], "service")
								EmitMetric(localTime, tags, metric, value, counter, c)

							//- if svname is not [BACKEND|FRONTEND] its a SERVER in a SERVICE and we prepend it with "server:"
							case len(pxnames) > 1 && (svname != "BACKEND" && svname != "FRONTEND"):
								tags = append(tags, "routes:"+pxnames[0], "services:"+pxnames[1], "servers:"+svname, "server")
								EmitMetric(localTime, tags, metric, value, counter, c)
							}
						}
					}
				}
			}
		}
	}
}

func EmitMetric(time string, tags []string, metric string, value string, counter *int64, c map[chan Metric]bool) {
	tags = append(tags, "metrics:"+metric)
	_type := "router-metric"
	metricValue, _ := strconv.Atoi(value)
	metricObj := Metric{tags, metricValue, time, _type}
	atomic.AddInt64(counter, 1)

	for s, _ := range c {
		s <- metricObj
	}
}

/*
  Logs a message every ticker interval giving an update on how many messages were parsed
*/
func (s *Streamer) Heartbeat() error {

	ticker := time.NewTicker(60 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				s.Log.Notice("Metrics parsed in last minute: %d \n", s.Counter)
				atomic.StoreInt64(&s.Counter, 0)
			}
		}
	}()
	return nil
}
