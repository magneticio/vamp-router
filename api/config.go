package api

import (
	"github.com/gin-gonic/gin"
	"github.com/magneticio/vamp-router/haproxy"
)

func GetConfig(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	c.JSON(200, Config(c))
}

func PostConfig(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	config := c.MustGet("haConfig").(*haproxy.Config)

	if c.Bind(&config) {
		HandleReload(c, config, 200, "updated config")
	} else {
		c.String(500, "Invalid JSON")
	}
}
