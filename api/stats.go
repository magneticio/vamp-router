package api

import (
	"github.com/gin-gonic/gin"
	"github.com/magneticio/vamp-loadbalancer/haproxy"
	"github.com/magneticio/vamp-loadbalancer/metrics"
)

func GetAllStats(c *gin.Context) {

	runtime := c.MustGet("haRuntime").(*haproxy.Runtime)
	status, err := runtime.GetStats("all")
	if err != nil {
		c.String(500, err.Error())
	} else {
		c.JSON(200, status)
	}

}

func GetBackendStats(c *gin.Context) {

	runtime := c.MustGet("haRuntime").(*haproxy.Runtime)
	status, err := runtime.GetStats("backend")
	if err != nil {
		c.String(500, err.Error())
	} else {
		c.JSON(200, status)
	}

}

func GetFrontendStats(c *gin.Context) {

	runtime := c.MustGet("haRuntime").(*haproxy.Runtime)
	status, err := runtime.GetStats("frontend")
	if err != nil {
		c.String(500, err.Error())
	} else {
		c.JSON(200, status)
	}
}

func GetServerStats(c *gin.Context) {

	runtime := c.MustGet("haRuntime").(*haproxy.Runtime)
	status, err := runtime.GetStats("server")
	if err != nil {
		c.String(500, err.Error())
	} else {
		c.JSON(200, status)
	}

}

func GetSSEStream(c *gin.Context) {
	sseBroker := c.MustGet("sseBroker").(*metrics.SSEBroker)
	sseBroker.ServeHTTP(c.Writer, c.Request)
}
