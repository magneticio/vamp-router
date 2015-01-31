package api

import (
  "github.com/gorilla/mux"
  "github.com/magneticio/vamp-loadbalancer/haproxy"
  "net/http"
  gologger "github.com/op/go-logging"
  "fmt"
)

type Router struct {
  log *gologger.Logger
  haConfig *haproxy.Config
  haRuntime *haproxy.Runtime
}

func (rt *Router) Init(port int, haConfig *haproxy.Config, haRuntime *haproxy.Runtime, log *gologger.Logger) {

  // set all the stuff we need
  rt.log = log
  rt.haConfig = haConfig
  rt.haRuntime = haRuntime

  fmt.Println(haConfig.TemplateFile)
    // create all routes based on the routes slice
  router := mux.NewRouter().StrictSlash(true)
  subRouter := router.PathPrefix("/v1").Subrouter()

  subRouter.HandleFunc("/config", rt.getConfig)


  http.ListenAndServe(":10003", router)  

}



