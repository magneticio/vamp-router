package metrics

import (
  "github.com/magneticio/vamp-loadbalancer/haproxy"
  "time"
  "encoding/json"
  "strings"
  "strconv"
  "reflect"
)


type Streamer struct {

  wantedMetrics []string
  haRuntime *haproxy.Runtime

}

// Just sets the metrics we want for now...
func (s *Streamer) Init(haRuntime *haproxy.Runtime) {

  s.wantedMetrics = []string{ "Scur", "Qcur","Smax","Slim","Weight","Qtime","Ctime","Rtime","Ttime","Req_rate","Req_rate_max","Req_tot","Rate","Rate_lim","Rate_max" }
  s.haRuntime = haRuntime
}

/* converts the haproxy metrics into a stream of discrete JSON object, like:
  {
   "name": "testbe.test_be_1.rate",   # The rate for server test_be_1 for proxy testbe
   "value": "2",                      # The value of the metric
   "timestamp": 1413546338            # The timestamp in Unix epoch
  }

  This stream can then be consumed by other streams like Kafka or SSE.  
*/ 
func (s *Streamer) ToJson(c chan []byte) error {

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
              json, err := json.MarshalIndent(metric, "", " ")

              if err != nil {
                return err
              }

              c <- json

            }
          }
        }
      }

    time.Sleep(3000 * time.Millisecond)
  }
}