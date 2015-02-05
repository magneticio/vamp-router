package api

import (
	"github.com/gin-gonic/gin"
	"github.com/magneticio/vamp-loadbalancer/haproxy"
)

func GetInfo(c *gin.Context) {

	runtime := c.MustGet("haRuntime").(*haproxy.Runtime)
	status, err := runtime.GetInfo()
	if err != nil {
		c.String(500, err.Error())
	} else {
		c.JSON(200, status)
	}

}
