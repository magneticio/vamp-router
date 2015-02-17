package api

import (
	"github.com/gin-gonic/gin"
	"github.com/magneticio/vamp-loadbalancer/haproxy"
)

func GetBackends(c *gin.Context) {

	config := c.MustGet("haConfig").(*haproxy.Config)

	result := config.GetBackends()
	if result != nil {
		c.JSON(200, result)
	} else {
		c.String(404, "no backends found")
	}

}

func GetBackend(c *gin.Context) {

	backend := c.Params.ByName("name")
	config := c.MustGet("haConfig").(*haproxy.Config)

	result := config.GetBackend(backend)
	if result != nil {
		c.JSON(200, result)
	} else {
		c.String(404, "no such backend")
	}

}

func GetServers(c *gin.Context) {

	backend := c.Params.ByName("name")
	config := c.MustGet("haConfig").(*haproxy.Config)

	result := config.GetServers(backend)
	if result != nil {
		c.JSON(200, result)
	} else {
		c.String(404, "no such server")
	}
}

func GetServer(c *gin.Context) {

	backend := c.Params.ByName("name")
	server := c.Params.ByName("server")
	config := c.MustGet("haConfig").(*haproxy.Config)

	result := config.GetServer(backend, server)
	if result != nil {
		c.JSON(200, result)
	} else {
		c.String(404, "no such server")
	}
}

func PostServer(c *gin.Context) {

	var server haproxy.ServerDetail
	backend := c.Params.ByName("name")
	config := c.MustGet("haConfig").(*haproxy.Config)

	if c.Bind(&server) {
		result := config.AddServer(backend, &server)
		if result {
			HandleReload(c, config, 201, "created server")
		} else {
			c.String(404, "no such backend")
		}
	} else {
		c.String(500, "Invalid JSON")
	}
}

func PutServerWeight(c *gin.Context) {

	var json UpdateWeight
	config := c.MustGet("haConfig").(*haproxy.Config)
	runtime := c.MustGet("haRuntime").(*haproxy.Runtime)
	backend := c.Params.ByName("name")
	server := c.Params.ByName("server")

	if c.Bind(&json) {
		status, err := runtime.SetWeight(backend, server, json.Weight)

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
				err = config.SetWeight(backend, server, json.Weight)
				if err != nil {
					c.String(500, err.Error())
				} else {
					HandleReload(c, config, 200, "updated server weight")
				}
			}
		}
	} else {
		c.String(500, "Invalid JSON")
	}
}

func DeleteServer(c *gin.Context) {

	backend := c.Params.ByName("name")
	server := c.Params.ByName("server")

	config := c.MustGet("haConfig").(*haproxy.Config)

	if config.DeleteServer(backend, server) {
		HandleReload(c, config, 200, "deleted server")
	} else {
		c.String(404, "no such server")
	}
}
