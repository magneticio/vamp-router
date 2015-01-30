package api

import (
  "github.com/gin-gonic/gin"
  "strconv"
  "github.com/magneticio/vamp-loadbalancer/haproxy"
)


func CreateApi(port int, haConfig *haproxy.Config, haRuntime *haproxy.Runtime ) {

  r := gin.New()
  r.Use(CORSMiddleware())
  r.Use(gin.Logger())
  r.Use(gin.Recovery())
  v1 := r.Group("/v1")

  {
    /*
      Backend Actions
     */
    v1.POST("/backend/:name/server/:server/weight/:weight", func(c *gin.Context){

        backend := c.Params.ByName("name")
        server :=  c.Params.ByName("server")
        weight,_  := strconv.Atoi(c.Params.ByName("weight"))

        status, err := haRuntime.SetWeight(backend, server, weight)

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
            err = haConfig.SetWeight(backend, server, weight)
            c.String(200,"Ok")
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

        status,err := haRuntime.SetAcl(backend,acl,pattern)

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

        status := haConfig.GetAcls(frontend)
        c.JSON(200, status)
      
      })


    /*

      Stats Actions

     */

    // get standard stats output from haproxy
    v1.GET("/stats", func(c *gin.Context) {
        status, err := haRuntime.GetStats("all")
        if err != nil {
          c.String(500, err.Error())
        } else {
          c.JSON(200, status)
        }

      })
    v1.GET("/stats/backend", func(c *gin.Context) {
        status, err := haRuntime.GetStats("backend")
        if err != nil {
          c.String(500, err.Error())
        } else {
          c.JSON(200, status)
        }

      })


    v1.GET("/stats/frontend", func(c *gin.Context) {
        status, err := haRuntime.GetStats("frontend")
        if err != nil {
          c.String(500, err.Error())
        } else {
          c.JSON(200, status)
        }

      })
    v1.GET("/stats/server", func(c *gin.Context) {
        status, err := haRuntime.GetStats("server")
        if err != nil {
          c.String(500, err.Error())
        } else {
          c.JSON(200, status)
        }

      })

    /*
      Full Config Actions
     */

    // get config file
    v1.GET("/config", func(c *gin.Context){
        c.JSON(200, haConfig)
    })

    // set config file
    v1.POST("/config", func(c *gin.Context){

        // clear the current config. This is dirty, should be able to clear the ConfigObj in
        // one go.
        haConfig.Frontends = [] *haproxy.Frontend{}
        haConfig.Backends = [] *haproxy.Backend{}

        c.Bind(&haConfig)
        err := haConfig.Render()
        if err != nil {
          c.String(500, "Error rendering config file")
          return
        } else {
          err = haRuntime.Reload(haConfig)
          if err != nil {
            c.String(500, "Error reloading the HAproxy configuration")
            return
          } else {
            c.String(200, "Ok")
          }

        }
    })

    /*
      Info
     */

    // get info on running process
    v1.GET("/info", func(c *gin.Context) {
        status, err := haRuntime.GetInfo()
        if err != nil {
          c.String(500, err.Error())
        } else {
          c.JSON(200, status)
        }

      })
  }

  // Listen and server on port
  r.Run("0.0.0.0:" + strconv.Itoa(port))
}

// override the standard Gin-Gonic middleware to add the CORS headers
func CORSMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {

    c.Writer.Header().Set("Content-Type", "application/json")
    c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
  }
}