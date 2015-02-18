package api

import (
	"github.com/gin-gonic/gin"
	"github.com/magneticio/vamp-loadbalancer/haproxy"
)

func GetRoutes(c *gin.Context) {
	
	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	result := Config(c).GetRoutes()
	if result != nil {
		c.JSON(200, result)
	} else {
		c.String(404, "no routes found")
	}

}

func GetRoute(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	routeName := c.Params.ByName("route")

	result, err := Config(c).GetRoute(routeName)
	if err != nil {
		c.String(404, err.Error())
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
		if err := Config(c).UpdateRoute(routeName,&route); err != nil {
			c.String(404, err.Error())
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
		if !Config(c).RouteExists(route.Name) {
			Config(c).AddRoute(&route)
			HandleReload(c, Config(c), 201, "created route")
		} else {
			c.String(409,"route already exists")
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
		c.String(404, err.Error())
	} else {
		HandleReload(c, Config(c), 200, "deleted route")
	}
}

func GetRouteGroups(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	routeName := c.Params.ByName("route")

	result, err := Config(c).GetRouteGroups(routeName)
	if err != nil {
		c.String(404, err.Error())
	} else {
		c.JSON(200, result)
	}
}

func GetRouteGroup(c *gin.Context) {
	
	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	routeName := c.Params.ByName("route")
	groupName := c.Params.ByName("group")

	result, err := Config(c).GetRouteGroup(routeName, groupName)
	if err != nil {
		c.String(404, err.Error())
	} else {
		c.JSON(200, result)
	}

}

func PutRouteGroup(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	var group haproxy.Group
	routeName := c.Params.ByName("route")
	groupName := c.Params.ByName("group")

	if c.Bind(&group) {
		if err := Config(c).UpdateRouteGroup(routeName,groupName,&group); err != nil {
			c.String(404, err.Error())
		} else {
			HandleReload(c, Config(c), 200, "updated group")
		}
	} else {
		c.String(500, "Invalid JSON")
	}
}


func PostRouteGroup(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	var group haproxy.Group
	routeName := c.Params.ByName("route")

	if c.Bind(&group) {
		if !Config(c).GroupExists(routeName,group.Name) {
			err := Config(c).AddRouteGroup(routeName,&group)
			if err != nil {
					c.String(404, err.Error())
				} else {
					HandleReload(c, Config(c), 201, "created group")
				}
		} else {
			c.String(409,"group already exists")
		}
	} else {
		c.String(500, "Invalid JSON")
	}
}

func DeleteRouteGroup(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	routeName := c.Params.ByName("route")
	groupName := c.Params.ByName("group")

	if err := Config(c).DeleteRouteGroup(routeName,groupName); err != nil {
		c.String(404, err.Error())
	} else {
		HandleReload(c, Config(c), 200, "deleted route")
	}
}

func GetGroupServers(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	routeName := c.Params.ByName("route")
	groupName := c.Params.ByName("group")

	result, err := Config(c).GetGroupServers(routeName,groupName)
	if err != nil {
		c.String(404, err.Error())
	} else {
		c.JSON(200, result)
	}
}

func GetGroupServer(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	routeName := c.Params.ByName("route")
	groupName := c.Params.ByName("group")
	serverName := c.Params.ByName("server")


	result, err := Config(c).GetGroupServer(routeName,groupName,serverName)
	if err != nil {
		c.String(404, err.Error())
	} else {
		c.JSON(200, result)
	}
}

func DeleteGroupServer(c *gin.Context) {
	
	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	routeName := c.Params.ByName("route")
	groupName := c.Params.ByName("group")
	serverName := c.Params.ByName("server")

	if err := Config(c).DeleteGroupServer(routeName,groupName,serverName); err != nil {
		c.String(404, err.Error())
	} else {
		HandleReload(c, Config(c), 200, "deleted server")
	}
}


func PostGroupServer(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	var server haproxy.Server
	routeName := c.Params.ByName("route")
	groupName := c.Params.ByName("group")

	if c.Bind(&server) {
		if !Config(c).ServerExists(routeName,groupName,server.Name) {
			err := Config(c).AddGroupServer(routeName,groupName,&server)
			if err != nil {
					c.String(404, err.Error())
				} else {
					HandleReload(c, Config(c), 201, "created server")
			}
		} else {
			c.String(409,"server already exists")
		}
	} else {
		c.String(500, "Invalid JSON")
	}
}

func PutGroupServer(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	var server haproxy.Server
	routeName := c.Params.ByName("route")
	groupName := c.Params.ByName("group")
	serverName := c.Params.ByName("server")

	if c.Bind(&server) {
		if err := Config(c).UpdateGroupServer(routeName,groupName,serverName,&server); err != nil {
			c.String(404, err.Error())
		} else {
			HandleReload(c, Config(c), 200, "updated server")
		}
	} else {
		c.String(500, "Invalid JSON")
	}
}
