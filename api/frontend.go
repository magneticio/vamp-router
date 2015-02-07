package api

import (
	"github.com/gin-gonic/gin"
	"github.com/magneticio/vamp-loadbalancer/haproxy"
)

func GetFrontends(c *gin.Context) {

	config := c.MustGet("haConfig").(*haproxy.Config)

	result := config.GetFrontends()
	if result != nil {
		c.JSON(200, result)
	} else {
		c.String(404, "no frontends found")
	}

}

func GetFrontend(c *gin.Context) {

	frontend := c.Params.ByName("name")
	config := c.MustGet("haConfig").(*haproxy.Config)

	result := config.GetFrontend(frontend)
	if result != nil {
		c.JSON(200, result)
	} else {
		c.String(404, "no such frontend")
	}

}

func PostFrontend(c *gin.Context) {

	var frontend haproxy.Frontend
	config := c.MustGet("haConfig").(*haproxy.Config)

	if c.Bind(&frontend) {
		config.AddFrontend(&frontend)
		HandleReload(c, config, 201, "created frontend")
	} else {
		c.String(500, "Invalid JSON")
	}
}

func DeleteFrontend(c *gin.Context) {

	frontendName := c.Params.ByName("name")
	config := c.MustGet("haConfig").(*haproxy.Config)

	if config.DeleteFrontend(frontendName) {
		HandleReload(c, config, 200, "deleted frontend")
	} else {
		c.String(404, "no such frontend")
	}
}

func GetFrontendACLs(c *gin.Context) {

	frontend := c.Params.ByName("name")
	config := c.MustGet("haConfig").(*haproxy.Config)

	status := config.GetAcls(frontend)
	c.JSON(200, status)

}

func PostFrontendACL(c *gin.Context) {

	var acl haproxy.ACL
	frontend := c.Params.ByName("name")
	config := c.MustGet("haConfig").(*haproxy.Config)

	if c.Bind(&acl) {
		config.AddAcl(frontend, &acl)
		HandleReload(c, config, 201, "created acl")
	} else {
		c.String(500, "Invalid JSON")
	}
}

func DeleteFrontendACL(c *gin.Context) {

	frontendName := c.Params.ByName("name")
	aclName := c.Params.ByName("acl_name")
	config := c.MustGet("haConfig").(*haproxy.Config)

	if config.DeleteAcl(frontendName, aclName) {
		HandleReload(c, config, 200, "deleted acl")
	} else {
		c.String(404, "no such acl")
	}
}
