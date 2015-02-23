package haproxy

import (
	"errors"
)

// gets all routes
func (c *Config) GetRoutes() []*Route {
	return c.Routes
}

// gets a route
func (c *Config) GetRoute(name string) (*Route, *Error) {

	var route *Route

	for _, rt := range c.Routes {
		if rt.Name == name {
			return rt, nil
			break
		}
	}
	return route, &Error{404, errors.New("no route found")}
}

// add a route to the configuration
func (c *Config) AddRoute(route *Route) *Error {

	if c.RouteExists(route.Name) {
		return &Error{409, errors.New("route already exists")}
	}

	// create some slices for all the stuff we are going to create. These are just holders so we can
	// iterate over them once we have created all the basic structures and add them to the configuration.
	feSlice := []*Frontend{}
	beSlice := []*Backend{}

	// When creating a new route, we have to create the stable frontends and backends.
	// 1. Check if the route exists
	// 1. Create stable Backend with empty server slice
	// 2. Create stable Frontend and add the stable Backend to it

	stableBackend := c.backendFactory(route.Name, route.Protocol, true, []*ServerDetail{})
	beSlice = append(beSlice, stableBackend)

	stableFrontend := c.frontendFactory(route.Name, route.Protocol, route.Port, stableBackend)
	feSlice = append(feSlice, stableFrontend)
	/*

		for groups

			1. Create socketServer
			2. Add it to the stable Backend
			3. Create Backend (with empty server slice)
			4. Create Frontend (set socket to the socketServer, add Backend)
	*/

	for _, group := range route.Groups {
		socketServer := c.socketServerFactory(ServerName(route.Name, group.Name), group.Weight)
		stableBackend.Servers = append(stableBackend.Servers, socketServer)

		backend := c.backendFactory(BackendName(route.Name, group.Name), route.Protocol, false, []*ServerDetail{})
		beSlice = append(beSlice, backend)

		frontend := c.socketFrontendFactory(FrontendName(route.Name, group.Name), route.Protocol, socketServer.UnixSock, backend)
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
func (c *Config) DeleteRoute(name string) *Error {

	for i, route := range c.Routes {

		if route.Name == name {

			// first remove all the frontends and backends related to the groups
			for _, group := range route.Groups {
				c.DeleteFrontend(FrontendName(route.Name, group.Name))
				c.DeleteBackend(BackendName(route.Name, group.Name))
			}

			// then remove the single backend and frontend
			c.DeleteFrontend(route.Name)
			c.DeleteBackend(route.Name)

			c.Routes = append(c.Routes[:i], c.Routes[i+1:]...)
			return nil
		}
	}
	return &Error{404, errors.New("no  route found")}
}

// just a convenience functions for a delete and a create
func (c *Config) UpdateRoute(name string, route *Route) *Error {

	if err := c.DeleteRoute(name); err != nil {
		return &Error{err.Code, err}
	}

	if err := c.AddRoute(route); err != nil {
		return &Error{err.Code, err}
	}
	return nil
}

func (c *Config) GetRouteGroups(name string) ([]*Group, *Error) {

	var groups []*Group

	for _, rt := range c.Routes {
		if rt.Name == name {
			return rt.Groups, nil
		}
	}
	return groups, &Error{404, errors.New("no groups found")}
}

func (c *Config) GetRouteGroup(routeName string, groupName string) (*Group, *Error) {

	var group *Group

	for _, rt := range c.Routes {
		if rt.Name == routeName {
			for _, grp := range rt.Groups {
				if grp.Name == groupName {
					return grp, nil
				}
			}
		}
	}
	return group, &Error{404, errors.New("no  group found")}
}

func (c *Config) AddRouteGroup(routeName string, group *Group) *Error {

	if c.GroupExists(routeName, group.Name) {
		return &Error{409, errors.New("group already exists")}
	}

	for _, route := range c.Routes {
		if route.Name == routeName {

			socketServer := c.socketServerFactory(ServerName(routeName, group.Name), group.Weight)
			backend := c.backendFactory(BackendName(route.Name, group.Name), route.Protocol, false, []*ServerDetail{})
			frontend := c.socketFrontendFactory(FrontendName(route.Name, group.Name), route.Protocol, socketServer.UnixSock, backend)

			for _, server := range group.Servers {
				srv := c.serverFactory(server.Name, group.Weight, server.Host, server.Port)
				backend.Servers = append(backend.Servers, srv)
			}

			if err := c.AddBackend(backend); err != nil {
				return &Error{500, errors.New("something went wrong adding backend: " + backend.Name)}
			}

			if err := c.AddFrontend(frontend); err != nil {
				return &Error{500, errors.New("something went wrong adding frontend: " + frontend.Name)}
			}

			route.Groups = append(route.Groups, group)
			return nil
		}
	}
	return &Error{404, errors.New("no  route found")}
}

func (c *Config) DeleteRouteGroup(routeName string, groupName string) *Error {

	for _, rt := range c.Routes {
		if rt.Name == routeName {
			for j, grp := range rt.Groups {
				if grp.Name == groupName {

					// order is important here. Always delete frontends first because they hold references to
					// backends. Deleting a backend that is still referenced first will fail.
					if err := c.DeleteFrontend(FrontendName(routeName, groupName)); err != nil {
						return &Error{500, errors.New("Something went wrong deleting frontend: " + FrontendName(routeName, groupName))}
					}

					if err := c.DeleteBackend(BackendName(routeName, groupName)); err != nil {
						return &Error{500, errors.New("Something went wrong deleting backend: " + BackendName(routeName, groupName))}
					}

					rt.Groups = append(rt.Groups[:j], rt.Groups[j+1:]...)
					return nil
				}
			}
		}
	}
	return &Error{404, errors.New("no  route found")}
}

// just a convenience functions for a delete and a create
func (c *Config) UpdateRouteGroup(routeName string, groupName string, group *Group) *Error {

	if err := c.DeleteRouteGroup(routeName, groupName); err != nil {
		return err
	}

	if err := c.AddRouteGroup(routeName, group); err != nil {
		return err
	}
	return nil
}

func (c *Config) GetGroupServers(routeName string, groupName string) ([]*Server, *Error) {

	var servers []*Server

	for _, rt := range c.Routes {
		if rt.Name == routeName {
			for _, grp := range rt.Groups {
				if grp.Name == groupName {
					return grp.Servers, nil
				}
			}
		}
	}
	return servers, &Error{404, errors.New("no servers found")}
}

func (c *Config) GetGroupServer(routeName string, groupName string, serverName string) (*Server, *Error) {

	var server *Server

	for _, rt := range c.Routes {
		if rt.Name == routeName {
			for _, grp := range rt.Groups {
				if grp.Name == groupName {
					for _, srv := range grp.Servers {
						if srv.Name == serverName {
							return srv, nil
						}
					}
				}
			}
		}
	}
	return server, &Error{404, errors.New("no server found")}
}

func (c *Config) DeleteGroupServer(routeName string, groupName string, serverName string) *Error {

	for _, rt := range c.Routes {
		if rt.Name == routeName {
			for _, grp := range rt.Groups {
				if grp.Name == groupName {
					for i, srv := range grp.Servers {
						if srv.Name == serverName {
							if err := c.DeleteServer(BackendName(routeName, groupName), serverName); err != nil {
								return &Error{500, err}
							}
							grp.Servers = append(grp.Servers[:i], grp.Servers[i+1:]...)
							return nil
						}
					}
				}
			}
		}
	}
	return &Error{404, errors.New("no server found")}
}

func (c *Config) AddGroupServer(routeName string, groupName string, server *Server) *Error {

	if c.ServerExists(routeName, groupName, server.Name) {
		return &Error{409, errors.New("server already exists")}
	}

	for _, route := range c.Routes {
		if route.Name == routeName {
			for _, group := range route.Groups {
				if group.Name == groupName {
					srvDetail := c.serverFactory(server.Name, group.Weight, server.Host, server.Port)
					c.AddServer(BackendName(routeName, groupName), srvDetail)
					group.Servers = append(group.Servers, server)
					return nil
				}
			}
		}
	}
	return &Error{404, errors.New("no group found")}
}

// just a convenience functions for a delete and a create
func (c *Config) UpdateGroupServer(routeName string, groupName string, serverName string, server *Server) *Error {

	if err := c.DeleteGroupServer(routeName, groupName, serverName); err != nil {
		return err
	}

	if err := c.AddGroupServer(routeName, groupName, server); err != nil {
		return err
	}
	return nil
}
