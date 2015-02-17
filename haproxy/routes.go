package haproxy

import (
	"errors"
)

// gets all routes
func (c *Config) GetRoutes() []*Route {

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	return c.Routes
}

// gets a route
func (c *Config) GetRoute(name string) (*Route, error) {

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	var result *Route

	for _, rt := range c.Routes {
		if rt.Name == name {
			return rt, nil
			break
		}
	}
	return result, errors.New("no route found")
}

// add a route to the configuration
func (c *Config) AddRoute(route *Route) error {

	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	// create some slices for all the stuff we are going to create. These are just holders so we can
	// iterate over them once we have created all the basic structures and add them to the configuration.
	feSlice := []*Frontend{}
	beSlice := []*Backend{}


	// When creating a new route, we have to create the stable frontends and backends.
	// 1. Check if the route exists
	// 1. Create stable Backend with empty server slice
	// 2. Create stable Frontend and add the stable Backend to it

	stableBackend := c.backendFactory(route.Name, route.Protocol, true, []*ServerDetail{})
	beSlice = append(beSlice,stableBackend)

	stableFrontend := c.frontendFactory(route.Name, route.Protocol, route.Port, stableBackend)
	feSlice = append(feSlice,stableFrontend)
/*

	for groups

		1. Create socketServer
		2. Add it to the stable Backend
		3. Create Backend (with empty server slice)
		4. Create Frontend (set socket to the socketServer, add Backend)
*/

		for _, group := range route.Groups {
				socketServer := c.socketServerFactory(route.Name + "." + group.Name, group.Weight)
				stableBackend.Servers = append(stableBackend.Servers,socketServer)

				backend := c.backendFactory(route.Name + "." + group.Name, route.Protocol, false, []*ServerDetail{})
				beSlice = append(beSlice,backend)

				frontend := c.socketFrontendFactory(route.Name + "." + group.Name, route.Protocol, socketServer.UnixSock, backend)
				feSlice = append(feSlice, frontend)

/*
		for servers
			1. Create Server
			2. Add Server to Backend Servers slice
*/
			for _, server := range group.Servers {
				srv := c.serverFactory(server.Name, group.Weight, server.Host, server.Port)
				backend.Servers = append(backend.Servers, srv)
			}
		}


	for _, fe := range feSlice {
		c.Frontends = append(c.Frontends, fe)
	}

	for _, be := range beSlice {
		c.Backends = append(c.Backends, be)
	}

	c.Routes = append(c.Routes, route)
	return nil
}

// deletes a route, cascading down the structure and remove all underpinning
// frontends, backends and servers.
func (c *Config) DeleteRoute(name string) error {

	c.Mutex.Lock()
	defer c.Mutex.Unlock()

  for i, route := range c.Routes {

    if route.Name == name {

    	// first remove all the frontends and backends related to the groups
    		for _, group := range route.Groups {

    		for j, backend := range c.Backends {
    			if backend.Name == route.Name + "." + group.Name {
    				c.Backends = append(c.Backends[:j], c.Backends[j+1:]...)
    			} 
    		}
    		for k, frontend := range c.Frontends {
    			if frontend.Name == route.Name + "." + group.Name {
    				c.Frontends = append(c.Frontends[:k], c.Frontends[k+1:]...)
    			} 
    		}				
    	}

    	// then remove the single backend and frontend
    		for l, backend := range c.Backends {
    			if backend.Name == route.Name {
    				c.Backends = append(c.Backends[:l], c.Backends[l+1:]...)
    			} 
    		}

    		for m, frontend := range c.Frontends {
    			if frontend.Name == route.Name {
    				c.Frontends = append(c.Frontends[:m], c.Frontends[m+1:]...)
    			} 
    		}

    	c.Routes = append(c.Routes[:i], c.Routes[i+1:]...)
    	return nil
    	break
    }
  }
  return errors.New("no such route found")
}