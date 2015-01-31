package metrics

import (
  "github.com/magneticio/vamp-loadbalancer/haproxy"
  "time"
  "strings"
  "strconv"
  "reflect"
)


type Streamer struct {

  wantedMetrics []string
  haRuntime *haproxy.Runtime
  pollFrequency int

}

// Just sets the metrics we want for now...
func (s *Streamer) Init(haRuntime *haproxy.Runtime, frequency int) {

  s.wantedMetrics = []string{ "Scur", "Qcur","Smax","Slim","Weight","Qtime","Ctime","Rtime","Ttime","Req_rate","Req_rate_max","Req_tot","Rate","Rate_lim","Rate_max" }
  s.haRuntime = haRuntime
  s.pollFrequency = frequency
}


/* 
  Generates an outgoing stream of discrete Metric struct values.
  This stream can then be consumed by other streams like Kafka or SSE. 
*/ 
func (s *Streamer) Out(c chan Metric) error {

  for  {

      stats, _ := s.haRuntime.GetStats("all")
      localTime := int64(time.Now().Unix())


    // for each proxy in the stats dump, pick out the wanted metrics.
      for _,proxy := range stats {

        // filter out the metrics for haproxy's own stats page
        if (proxy.Pxname != "stats") {

          // loop over all wanted metrics for the current proxy
          for _, metric := range s.wantedMetrics {

            fullMetricName := proxy.Pxname + "." + strings.ToLower(proxy.Svname) + "." + strings.ToLower(metric)
            field  := reflect.ValueOf(proxy).FieldByName(metric).String()
            if (field != "") {

              metricValue,_ := strconv.Atoi(field)
              metric := Metric{fullMetricName, metricValue, localTime}

              c <- metric

            }
          }
        }
      }

    time.Sleep(time.Duration(s.pollFrequency) * time.Millisecond)
  }

}
