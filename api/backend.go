package api

import (
	"github.com/gin-gonic/gin"
	"github.com/magneticio/vamp-loadbalancer/haproxy"
)

func GetBackends(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	result := Config(c).GetBackends()
	if result != nil {
		c.JSON(200, result)
	} else {
		c.String(404, "no backends found")
	}

}

func GetBackend(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	backend := c.Params.ByName("name")

	result := Config(c).GetBackend(backend)
	if result != nil {
		c.JSON(200, result)
	} else {
		c.String(404, "no such backend")
	}

}

func PostBackend(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	var backend haproxy.Backend

	if c.Bind(&backend) {

		if !Config(c).BackendExists(backend.Name){
		Config(c).AddBackend(&backend)
		HandleReload(c, Config(c), 201, "created backend")
		} else {
			c.String(409,"backend already exists")
		}
	} else {
		c.String(500, "Invalid JSON")
	}
}

func DeleteBackend(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	name := c.Params.ByName("name")

	if err := Config(c).DeleteBackend(name); err != nil {
		c.String(404, err.Error())
	} else {
		HandleReload(c, Config(c), 200, "deleted backend")
	}
}

func GetServers(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()
	
	backend := c.Params.ByName("name")

	result := Config(c).GetServers(backend)
	if result != nil {
		c.JSON(200, result)
	} else {
		c.String(404, "no such server")
	}
}

func GetServer(c *gin.Context) {

	Config(c).BeginReadTrans()
	defer Config(c).EndReadTrans()

	backend := c.Params.ByName("name")
	server := c.Params.ByName("server")

	result := Config(c).GetServer(backend, server)
	if result != nil {
		c.JSON(200, result)
	} else {
		c.String(404, "no such server")
	}
}

func PostServer(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	var server haproxy.ServerDetail
	backend := c.Params.ByName("name")

	if c.Bind(&server) {
		if err := Config(c).AddServer(backend, &server); err != nil {
			c.String(404, err.Error())
		} else {
			HandleReload(c, Config(c), 201, "created server")
		} 
	} else {
		c.String(500, "Invalid JSON")
	}
}

func PutServerWeight(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	var json UpdateWeight
	backend := c.Params.ByName("name")
	server := c.Params.ByName("server")

	if c.Bind(&json) {
		status, err := Runtime(c).SetWeight(backend, server, json.Weight)

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

				//update the Config(c) object with the new weight
				err = Config(c).SetWeight(backend, server, json.Weight)
				if err != nil {
					c.String(500, err.Error())
				} else {
					HandleReload(c, Config(c), 200, "updated server weight")
				}
			}
		}
	} else {
		c.String(500, "Invalid JSON")
	}
}

func DeleteServer(c *gin.Context) {

	Config(c).BeginWriteTrans()
	defer Config(c).EndWriteTrans()

	backend := c.Params.ByName("name")
	server := c.Params.ByName("server")

	if err := Config(c).DeleteServer(backend, server); err != nil {
		c.String(404, "no such server")
	} else {
		HandleReload(c, Config(c), 200, "deleted server")
	}
}
