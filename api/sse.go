package api

import (
  "github.com/gorilla/mux"
  "github.com/magneticio/vamp-loadbalancer/metrics"
  "net/http"
)


func CreateSSE(b *metrics.SSEBroker) {

  r := mux.NewRouter()
  r.HandleFunc("/", b.ServeHTTP)
  http.ListenAndServe(":10002", r)  
}

