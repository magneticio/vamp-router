package api

import (
  "encoding/json"
  "net/http"
  "fmt"
)

// Config

func (rt *Router) getConfig(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  fmt.Println("in http routine:" + rt.haConfig.TemplateFile)
  json.NewEncoder(w).Encode(rt.haConfig)
}

// Stats

func (rt *Router) GetAllStats(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(rt.haConfig)
}
