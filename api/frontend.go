package api

import (
	"github.com/gin-gonic/gin"
	"github.com/magneticio/vamp-loadbalancer/haproxy"
)

func PostAclPattern(c *gin.Context) {

	backend := c.Params.ByName("name")
	acl := c.Params.ByName("acl")
	pattern := c.Params.ByName("pattern")
	runtime := c.MustGet("haRuntime").(*haproxy.Runtime)

	status, err := runtime.SetAcl(backend, acl, pattern)

	// check on Runtime errors
	if err != nil {
		c.String(500, err.Error())
	} else {
		switch status {
		case "No such backend.\n\n":
			c.String(404, status)
		default:

			//update the config object with the new acl
			//err = UpdateWeightInConfig(backend, server, weight, ConfigObj)
			c.String(200, "Ok")
		}
	}
}

func GetACLs(c *gin.Context) {

	frontend := c.Params.ByName("name")
	config := c.MustGet("haConfig").(*haproxy.Config)

	status := config.GetAcls(frontend)
	c.JSON(200, status)

}
