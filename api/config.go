package api

import (
	"github.com/gin-gonic/gin"
	"github.com/magneticio/vamp-router/haproxy"
	"net/http"
)

func GetConfig(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	c.JSON(http.StatusOK, Config(c))
}

func PostConfig(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	config := c.MustGet("haConfig").(*haproxy.Config)

	if c.Bind(&config) {
		HandleReload(c, config, http.StatusOK, gin.H{"status": "updated config"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request"})
	}
}
