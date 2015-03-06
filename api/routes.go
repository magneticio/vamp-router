package api

import (
	"github.com/gin-gonic/gin"
	"github.com/magneticio/vamp-loadbalancer/haproxy"
)

func GetRoutes(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	result := Config(c).GetRoutes()
	if Config(c).GetRoutes() != nil {
		c.JSON(200, result)
	} else {
		c.String(404, "no routes found")
	}

}

func GetRoute(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	routeName := c.Params.ByName("route")

	if result, err := Config(c).GetRoute(routeName); err != nil {
		HandleError(c, err)
	} else {
		c.JSON(200, result)
	}
}

func PutRoute(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	var route haproxy.Route
	routeName := c.Params.ByName("route")

	if c.Bind(&route) {
		if err := Config(c).UpdateRoute(routeName, &route); err != nil {
			HandleError(c, err)
		} else {
			HandleReload(c, Config(c), 200, "updated route")
		}
	} else {
		c.String(500, "Invalid JSON")
	}
}

func PostRoute(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	var route haproxy.Route

	if c.Bind(&route) {
		if err := Config(c).AddRoute(&route); err != nil {
			HandleError(c, err)
		} else {
			HandleReload(c, Config(c), 201, "created route")
		}
	} else {
		c.String(500, "Invalid JSON")
	}
}

func DeleteRoute(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	routeName := c.Params.ByName("route")

	if err := Config(c).DeleteRoute(routeName); err != nil {
		HandleError(c, err)
	} else {
		HandleReload(c, Config(c), 204, "")
	}
}

func GetRouteServices(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	routeName := c.Params.ByName("route")

	result, err := Config(c).GetRouteServices(routeName)
	if err != nil {
		HandleError(c, err)
	} else {
		c.JSON(200, result)
	}
}

func GetRouteService(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	routeName := c.Params.ByName("route")
	serviceName := c.Params.ByName("service")

	result, err := Config(c).GetRouteService(routeName, serviceName)
	if err != nil {
		HandleError(c, err)
	} else {
		c.JSON(200, result)
	}

}

func PutRouteService(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	var service haproxy.Service
	routeName := c.Params.ByName("route")
	serviceName := c.Params.ByName("service")

	if c.Bind(&service) {
		if err := Config(c).UpdateRouteService(routeName, serviceName, &service); err != nil {
			HandleError(c, err)
		} else {
			HandleReload(c, Config(c), 200, "updated service")
		}
	} else {
		c.String(500, "Invalid JSON")
	}
}

func PutRouteServices(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	var services []haproxy.Service
	routeName := c.Params.ByName("route")

	if c.Bind(&services) {
		if err := Config(c).UpdateRouteServices(routeName, &services); err != nil {
			HandleError(c, err)
		} else {
			HandleReload(c, Config(c), 200, "updated services")
		}
	} else {
		c.String(500, "Invalid JSON")
	}
}

func PostRouteService(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	var services []*haproxy.Service
	routeName := c.Params.ByName("route")

	if c.Bind(&services) {
		if err := Config(c).AddRouteServices(routeName, services); err != nil {
			HandleError(c, err)
		} else {
			HandleReload(c, Config(c), 201, "created service(s)")
		}
	} else {
		c.String(500, "Invalid JSON")
	}
}

func DeleteRouteService(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	routeName := c.Params.ByName("route")
	serviceName := c.Params.ByName("service")

	if err := Config(c).DeleteRouteService(routeName, serviceName); err != nil {
		HandleError(c, err)
	} else {
		HandleReload(c, Config(c), 204, "")
	}
}

func GetServiceServers(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	routeName := c.Params.ByName("route")
	serviceName := c.Params.ByName("service")

	result, err := Config(c).GetServiceServers(routeName, serviceName)
	if err != nil {
		HandleError(c, err)
	} else {
		c.JSON(200, result)
	}
}

func GetServiceServer(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	routeName := c.Params.ByName("route")
	serviceName := c.Params.ByName("service")
	serverName := c.Params.ByName("server")

	result, err := Config(c).GetServiceServer(routeName, serviceName, serverName)
	if err != nil {
		HandleError(c, err)
	} else {
		c.JSON(200, result)
	}
}

func DeleteServiceServer(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	routeName := c.Params.ByName("route")
	serviceName := c.Params.ByName("service")
	serverName := c.Params.ByName("server")

	if err := Config(c).DeleteServiceServer(routeName, serviceName, serverName); err != nil {
		HandleError(c, err)
	} else {
		HandleReload(c, Config(c), 204, "")
	}
}

func PostServiceServer(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	var server haproxy.Server
	routeName := c.Params.ByName("route")
	serviceName := c.Params.ByName("service")

	if c.Bind(&server) {
		if err := Config(c).AddServiceServer(routeName, serviceName, &server); err != nil {
			HandleError(c, err)
		} else {
			HandleReload(c, Config(c), 201, "created server")
		}
	} else {
		c.String(500, "Invalid JSON")
	}
}

func PutServiceServer(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	var server haproxy.Server
	routeName := c.Params.ByName("route")
	serviceName := c.Params.ByName("service")
	serverName := c.Params.ByName("server")

	if c.Bind(&server) {
		if err := Config(c).UpdateServiceServer(routeName, serviceName, serverName, &server); err != nil {
			HandleError(c, err)
		} else {
			HandleReload(c, Config(c), 200, "updated server")
		}
	} else {
		c.String(500, "Invalid JSON")
	}
}
