package api

import (
	"github.com/gin-gonic/gin"
	"github.com/magneticio/vamp-loadbalancer/haproxy"
	"github.com/magneticio/vamp-loadbalancer/metrics"
	gologger "github.com/op/go-logging"
	"strconv"
)

func CreateApi(port int, haConfig *haproxy.Config, haRuntime *haproxy.Runtime, log *gologger.Logger, SSEBroker *metrics.SSEBroker) {

	gin.SetMode("release")

	r := gin.New()
	r.Use(HaproxyMiddleware(haConfig, haRuntime))
	r.Use(LoggerMiddleware(log))
	r.Use(SSEMiddleware(SSEBroker))
	r.Use(gin.Recovery())
	r.Static("/www", "./www")
	v1 := r.Group("/v1")

	{
		r.GET("/", func(c *gin.Context) {
			c.Redirect(301, "www/index.html")
		})

		/*
		   Frontend
		*/
		v1.GET("/frontends", GetFrontends)
		v1.POST("frontends/:name/acls", PostFrontendACL)
		v1.GET("/frontends/:name/acls", GetFrontendACLs)
		v1.DELETE("/frontends/:name/acls/:acl_name", DeleteFrontendACL)
		v1.GET("/frontends/:name", GetFrontend)
		v1.DELETE("/frontends/:name", DeleteFrontend)
		v1.POST("/frontends", PostFrontend)

		/*
		   Backend
		*/
		v1.GET("/backends", GetBackends)
		v1.GET("/backends/:name", GetBackend)
		v1.GET("/backends/:name/servers", GetServers)
		v1.GET("/backends/:name/servers/:server", GetServer)
		v1.PUT("/backends/:name/servers/:server", PutServerWeight)
		v1.POST("/backends/:name/servers", PostServer)
		v1.DELETE("/backends/:name/servers/:server", DeleteServer)

		/*
		   Stats
		*/
		v1.GET("/stats", GetAllStats)
		v1.GET("/stats/backends", GetBackendStats)
		v1.GET("/stats/frontends", GetFrontendStats)
		v1.GET("/stats/servers", GetServerStats)
		v1.GET("/stats/stream", GetSSEStream)

		/*
		   Config
		*/
		v1.GET("/config", GetConfig)
		v1.POST("/config", PostConfig)

		/*
			Routes
		*/
		v1.GET("/routes", GetRoutes)
		v1.GET("/routes/:name", GetRoute)
		v1.POST("/routes", PostRoute)
		v1.DELETE("/routes/:name", DeleteRoute)

		/*
		   Info
		*/
		v1.GET("/info", GetInfo)
	}

	// Listen and server on port
	r.Run("0.0.0.0:" + strconv.Itoa(port))
}

func HandleReload(c *gin.Context, config *haproxy.Config, status int, message string) {

	runtime := c.MustGet("haRuntime").(*haproxy.Runtime)

	err := config.RenderAndPersist()
	if err != nil {
		c.String(500, "Error rendering config file")
		return
	}

	err = runtime.Reload(config)
	if err != nil {
		c.String(500, "Error reloading the HAproxy configuration")
		return
	}

	c.String(status, message)
}
