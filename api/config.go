package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/magneticio/vamp-router/haproxy"
)

func GetConfig(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	c.JSON(http.StatusOK, Config(c))
}

func PostConfig(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	var config haproxy.Config

	if err := c.Bind(&config); err == nil {
		if err := Config(c).UpdateConfig(&config); err != nil {
			HandleError(c, err)
		} else {
			HandleReload(c, Config(c), http.StatusCreated, gin.H{"status": "updated config"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"status": "bad request", "error": err.Error()})
	}
}
