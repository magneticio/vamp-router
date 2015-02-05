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
}
