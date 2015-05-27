package metrics

import (
	"github.com/magneticio/vamp-router/haproxy"
	gologger "github.com/op/go-logging"
	"strconv"
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
	s.wantedMetrics = []string{"scur", "qcur", "smax", "slim", "qtime", "ctime", "rtime", "ttime", "req_rate", "req_rate_max", "req_tot", "rate", "rate_lim", "rate_max"}
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

	localTime := time.Now().Format(time.RFC3339)

	for {

		select {

		case stats := <-statsChannel:

			// for each proxy in the stats dump, pick out the wanted metrics.

			for _, proxy := range stats {

				// loop over all wanted metrics for the current proxy
				for _, metric := range wantedMetrics {

					// compile tags
					// proxies := strings.Split(proxy.Pxname, ".")

					/* we add a specific "route" tag (that doesn't exist in the standard haproxy stats) to
					all "top-level" frontend and backends. We check when to insert this by the length of the
					proxies slice after splitting. One item equals a route.
					*/
					// tags := append(proxies, []string{strings.ToLower(proxy.Svname), strings.ToLower(metric)}...)

					// if len(proxies) == 1 && (proxy.Svname == "BACKEND" || proxy.Svname == "FRONTEND") {
					// 	tags = append(tags, "route")
					// }

					tags := []string{proxy["pxname"]}
					field := proxy[metric]

					if field != "" {

						metricValue, _ := strconv.Atoi(field)
						metric := Metric{tags, metricValue, localTime}
						// fmt.Printf("compiled metric => tags: %s, value: %d, time: %s \n", tags, metricValue, localTime)
						atomic.AddInt64(counter, 1)

						for s, _ := range c {
							s <- metric
						}

					}
				}
			}
		}
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
