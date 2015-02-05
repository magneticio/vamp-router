package api

import (
	"github.com/gin-gonic/gin"
	"github.com/magneticio/vamp-loadbalancer/haproxy"
)

func PutBackendWeight(c *gin.Context) {

	var json UpdateWeight
	config := c.MustGet("haConfig").(*haproxy.Config)
	runtime := c.MustGet("haRuntime").(*haproxy.Runtime)

	valid := c.Bind(&json)
	if valid != true {
		c.String(500, "Invalid Json")
	} else {

		backend := c.Params.ByName("name")
		server := c.Params.ByName("server")

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
					c.String(200, "Ok")
				}
			}
		}
	}
}
