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

	result, err := config.GetRoute(route)
	if err != nil {
		c.String(404, err.Error())
	} else {
		c.JSON(200, result)
	}
}

func PutRoute(c *gin.Context) {

	var route haproxy.Route
	config := c.MustGet("haConfig").(*haproxy.Config)
	name := c.Params.ByName("name")

	if c.Bind(&route) {
		if err := config.DeleteRoute(name); err != nil {
			c.String(404, err.Error())
		} else {
			config.AddRoute(&route)
			HandleReload(c, config, 200, "updated route")
		}
	} else {
		c.String(500, "Invalid JSON")
	}
}

func PostRoute(c *gin.Context) {

	var route haproxy.Route
	config := c.MustGet("haConfig").(*haproxy.Config)

	if c.Bind(&route) {
		if !config.RouteExists(route.Name) {
			config.AddRoute(&route)
			HandleReload(c, config, 201, "created route")
		} else {
			c.String(409,"route already exists")
		}
	} else {
		c.String(500, "Invalid JSON")
	}
}

func DeleteRoute(c *gin.Context) {

	name := c.Params.ByName("name")
	config := c.MustGet("haConfig").(*haproxy.Config)

	if err := config.DeleteRoute(name); err != nil {
		c.String(404, err.Error())
	} else {
		HandleReload(c, config, 200, "deleted route")
	}
}
