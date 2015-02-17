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

func GetFrontendFilters(c *gin.Context) {

	frontend := c.Params.ByName("name")
	config := c.MustGet("haConfig").(*haproxy.Config)

	status := config.GetFilters(frontend)
	c.JSON(200, status)

}

func PostFrontendFilter(c *gin.Context) {

	var Filter haproxy.Filter
	frontend := c.Params.ByName("name")
	config := c.MustGet("haConfig").(*haproxy.Config)

	if c.Bind(&Filter) {
		config.AddFilter(frontend, &Filter)
		HandleReload(c, config, 201, "created Filter")
	} else {
		c.String(500, "Invalid JSON")
	}
}

func DeleteFrontendFilter(c *gin.Context) {

	frontendName := c.Params.ByName("name")
	FilterName := c.Params.ByName("Filter_name")
	config := c.MustGet("haConfig").(*haproxy.Config)

	if config.DeleteFilter(frontendName, FilterName) {
		HandleReload(c, config, 200, "deleted Filter")
	} else {
		c.String(404, "no such Filter")
	}
}
