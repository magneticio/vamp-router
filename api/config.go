package api

import (
	"github.com/gin-gonic/gin"
	"github.com/magneticio/vamp-loadbalancer/haproxy"
)

func GetConfig(c *gin.Context) {
	config := c.MustGet("haConfig").(*haproxy.Config)
	c.JSON(200, config)
}

func PostConfig(c *gin.Context) {

	config := c.MustGet("haConfig").(*haproxy.Config)

	if c.Bind(&config) {
		HandleReload(c, config, 200, "updated config")
	} else {
		c.String(500, "Invalid JSON")
	}
}
