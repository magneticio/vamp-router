package api

import (
	"github.com/gin-gonic/gin"
	"github.com/magneticio/vamp-loadbalancer/haproxy"
	"github.com/magneticio/vamp-loadbalancer/metrics"
	gologger "github.com/op/go-logging"
	"strconv"
)

func CreateApi(port int, haConfig *haproxy.Config, haRuntime *haproxy.Runtime, log *gologger.Logger, SSEBroker *metrics.SSEBroker) {

	r := gin.New()
	r.Use(HaproxyMiddleware(haConfig, haRuntime))
	r.Use(LoggerMiddleware(log))
	r.Use(gin.Recovery())
	r.Static("/www", "./www")
	v1 := r.Group("/v1")

	{
		r.GET("/", func(c *gin.Context) {
			c.Redirect(301, "www/index.html")
		})

		/*
		   Backend
		*/
		v1.PUT("/backend/:name/server/:server", PutBackendWeight)

		/*
		   Frontend
		*/
		v1.POST("frontend/:name/acl/:acl/:pattern", PostAclPattern)
		v1.GET("/frontend/:name/acls", GetACLs)

		/*
		   Stats
		*/
		v1.GET("/stats", GetAllStats)
		v1.GET("/stats/backend", GetBackendStats)
		v1.GET("/stats/frontend", GetFrontendStats)
		v1.GET("/stats/server", GetServerStats)
		v1.GET("/stats/stream", SSEMiddleware(SSEBroker), GetSSEStream)

		/*
		   Full Config Actions
		*/
		v1.GET("/config", GetConfig)
		v1.POST("/config", PostConfig)

		/*
		   Info
		*/
		v1.GET("/info", GetInfo)
	}

	// Listen and server on port
	r.Run("0.0.0.0:" + strconv.Itoa(port))

}
