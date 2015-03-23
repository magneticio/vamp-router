package metrics

import (
	"github.com/magneticio/vamp-router/haproxy"
	gologger "github.com/op/go-logging"
	"reflect"
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
	s.wantedMetrics = []string{"Scur", "Qcur", "Smax", "Slim", "Weight", "Qtime", "Ctime", "Rtime", "Ttime", "Req_rate", "Req_rate_max", "Req_tot", "Rate", "Rate_lim", "Rate_max"}
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

	for {

		stats, _ := s.haRuntime.GetStats("all")
		localTime := time.Now().Format(time.RFC3339)

		// for each proxy in the stats dump, pick out the wanted metrics.
		for _, proxy := range stats {

			// filter out the metrics for haproxy's own stats page
			if proxy.Pxname != "stats" {

				// loop over all wanted metrics for the current proxy
				for _, metric := range s.wantedMetrics {

					// compile tags
					proxies := strings.Split(proxy.Pxname, ".")
					tags := append(proxies, []string{strings.ToLower(proxy.Svname), strings.ToLower(metric)}...)
					field := reflect.ValueOf(proxy).FieldByName(metric).String()
					if field != "" {

						metricValue, _ := strconv.Atoi(field)
						metric := Metric{tags, metricValue, localTime}
						atomic.AddInt64(&s.Counter, 1)

						for s, _ := range s.Clients {
							s <- metric
						}

					}
				}
			}
		}
		time.Sleep(time.Duration(s.pollFrequency) * time.Millisecond)
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
