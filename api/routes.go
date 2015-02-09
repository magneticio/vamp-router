package api

import (
	"github.com/gin-gonic/gin"
	"github.com/magneticio/vamp-loadbalancer/haproxy"
)

func GetRoutes(c *gin.Context) {

	config := c.MustGet("haConfig").(*haproxy.Config)

	result := config.GetRoutes()
	if result != nil {
		c.JSON(200, result)
	} else {
		c.String(404, "no routes found")
	}

}

func GetRoute(c *gin.Context) {

	route := c.Params.ByName("name")
	config := c.MustGet("haConfig").(*haproxy.Config)

	result := config.GetRoute(route)
	if result != nil {
		c.JSON(200, result)
	} else {
		c.String(404, "no such route")
	}

}

func PostRoute(c *gin.Context) {

	var newRoute haproxy.NewRoute
	config := c.MustGet("haConfig").(*haproxy.Config)

	if c.Bind(&newRoute) {
		config.AddRoute(&newRoute)
		HandleReload(c, config, 201, "created route")
	} else {
		c.String(500, "Invalid JSON")
	}
}

func DeleteRoute(c *gin.Context) {

	routeName := c.Params.ByName("name")
	config := c.MustGet("haConfig").(*haproxy.Config)

	if config.DeleteRoute(routeName) {
		HandleReload(c, config, 200, "deleted route")
	} else {
		c.String(404, "no such route")
	}
}
