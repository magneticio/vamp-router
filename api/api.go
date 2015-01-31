package api

import (
  "github.com/gin-gonic/gin"
  "strconv"
  "github.com/magneticio/vamp-loadbalancer/haproxy"
  "github.com/magneticio/vamp-loadbalancer/metrics"
  gologger "github.com/op/go-logging"
  "net/http"
  )


func CreateApi(port int, haConfig *haproxy.Config, haRuntime *haproxy.Runtime, log *gologger.Logger, SSEBroker *metrics.SSEBroker) {

  r := gin.New()
  r.Use(CORSMiddleware())
  r.Use(HaproxyMiddleware(haConfig, haRuntime))
  r.Use(LoggerMiddleware(log))
  r.Use(gin.Recovery())
  v1 := r.Group("/v1")

  {
    /*
      Backend Actions
     */
    v1.PUT("/backend/:name/server/:server", func(c *gin.Context){

        var json UpdateWeight
        config := c.MustGet("haConfig").(*haproxy.Config)
        runtime := c.MustGet("haRuntime").(*haproxy.Runtime)

        valid := c.Bind(&json)
        if valid != true {
          c.String(500, "Invalid Json")
        } else {

          backend := c.Params.ByName("name")
          server :=  c.Params.ByName("server")

          status, err := runtime.SetWeight(backend, server, json.Weight)

          // check on Runtime errors
          if err != nil {
            c.String(500, err.Error())
          } else {

            switch status {
            case "No such server.\n\n":
              c.String(404, status)
            case "No such backend.\n\n":
              c.String(404, status)
            default:

              //update the config object with the new weight
              err = config.SetWeight(backend, server, json.Weight)
              if err != nil {
                c.String(500, err.Error())
              } else {
                c.String(200,"Ok")
              }
            }
          }
        }
      })

    /*
      Frontend Actions
     */

    v1.POST("frontend/:name/acl/:acl/:pattern",func(c *gin.Context){

        backend := c.Params.ByName("name")
        acl := c.Params.ByName("acl")
        pattern := c.Params.ByName("pattern")
        runtime := c.MustGet("haRuntime").(*haproxy.Runtime)

        status,err := runtime.SetAcl(backend,acl,pattern)

        // check on Runtime errors
        if err != nil {
          c.String(500, err.Error())
        } else {
          switch status {
          case "No such backend.\n\n":
            c.String(404, status)
          default:

            //update the config object with the new acl
            //err = UpdateWeightInConfig(backend, server, weight, ConfigObj)
            c.String(200,"Ok")
          }
        }
      })

    v1.GET("/frontend/:name/acls",func(c *gin.Context){

        frontend := c.Params.ByName("name")
        config := c.MustGet("haConfig").(*haproxy.Config)

        status := config.GetAcls(frontend)
        c.JSON(200, status)
      
      })


    /*

      Stats Actions

     */

    // get standard stats output from haproxy
    v1.GET("/stats", func(c *gin.Context) {

        runtime := c.MustGet("haRuntime").(*haproxy.Runtime)
        status, err := runtime.GetStats("all")
        if err != nil {
          c.String(500, err.Error())
        } else {
          c.JSON(200, status)
        }

      })
    v1.GET("/stats/backend", func(c *gin.Context) {

        runtime := c.MustGet("haRuntime").(*haproxy.Runtime)
        status, err := runtime.GetStats("backend")
        if err != nil {
          c.String(500, err.Error())
        } else {
          c.JSON(200, status)
        }

      })


    v1.GET("/stats/frontend", func(c *gin.Context) {

        runtime := c.MustGet("haRuntime").(*haproxy.Runtime)
        status, err := runtime.GetStats("frontend")
        if err != nil {
          c.String(500, err.Error())
        } else {
          c.JSON(200, status)
        }

      })
    v1.GET("/stats/server", func(c *gin.Context) {

        runtime := c.MustGet("haRuntime").(*haproxy.Runtime)
        status, err := runtime.GetStats("server")
        if err != nil {
          c.String(500, err.Error())
        } else {
          c.JSON(200, status)
        }

      })
    v1.GET("/stats/stream",SSEMiddleware(SSEBroker), func(c *gin.Context) {

        sseBroker := c.MustGet("sseBroker").(*metrics.SSEBroker)
        sseBroker.ServeHTTP(c.Writer,c.Request)

    })    



    /*
      Full Config Actions
     */

    // get config file
    v1.GET("/config", func(c *gin.Context){
        config := c.MustGet("haConfig").(*haproxy.Config)
        c.JSON(200, config)
    })

    // set config file
    v1.POST("/config", func(c *gin.Context){

        config := c.MustGet("haConfig").(*haproxy.Config)
        runtime := c.MustGet("haRuntime").(*haproxy.Runtime)

        valid := c.Bind(&config)

        if valid != true {
          c.String(500, "Invalid JSON")

        } 

        if valid == true {
          err := config.RenderAndPersist()
          if err != nil {
            c.String(500, "Error rendering config file")
            return
          } else {
            err = runtime.Reload(config)
            if err != nil {
              c.String(500, "Error reloading the HAproxy configuration")
              return
            } else {
              c.String(200, "Ok")
            }

          }
        }
    })

    /*
      Info
     */

    // get info on running process
    v1.GET("/info", func(c *gin.Context) {

        runtime := c.MustGet("haRuntime").(*haproxy.Runtime)
        status, err := runtime.GetInfo()
        if err != nil {
          c.String(500, err.Error())
        } else {
          c.JSON(200, status)
        }

      })
  }

  // Listen and server on port
  // r.Run("0.0.0.0:" + strconv.Itoa(port))

  http.ListenAndServe(":" + strconv.Itoa(port), r)

}


