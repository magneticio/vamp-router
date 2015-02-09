package haproxy

/*

  AddRoutes is a convenience method to create a set of frontends, backends etc.that together form what we
  call a "route". You could create this structure by hand with separate API calls, but this is faster and
  easier in 9 out of 10 cases.

  The structure of a route is as follows. Each

                            -> [server a] -> socket -> [frontend a: backend a] -> [*servers] -> host:port
                          /
  ->[frontend : backend]-
                          \
                            -> [server b] -> socket -> [frontend b: backend b] -> [*servers] -> host:port

  *Note: the servers at the end of the route are not created in this routine.

*/

func (c *Config) AddRoute(newRoute *NewRoute) error {

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	// we create the routes structure by starting "at the back": the a/b servers that use socket
	// to direct traffic to frontends listening on those sockets

	stableServerA := defaultSocketProxyServer(newRoute.Name, "a", 100)
	stableServerB := defaultSocketProxyServer(newRoute.Name, "b", 0)

	backendA := defaultBackend(newRoute.Name, "a", newRoute.Mode, false, []*BackendServer{})
	backendB := defaultBackend(newRoute.Name, "b", newRoute.Mode, false, []*BackendServer{})

	frontendA := defaultSocketProxyFrontend(newRoute.Name, "a", newRoute.Mode, stableServerA.UnixSock, backendA)
	frontendB := defaultSocketProxyFrontend(newRoute.Name, "b", newRoute.Mode, stableServerB.UnixSock, backendB)

	stableBackend := defaultBackend(newRoute.Name, "", newRoute.Mode, true, []*BackendServer{stableServerA, stableServerB})
	stableFrontend := defaultFrontend(newRoute.Name, newRoute.Mode, newRoute.Endpoint, stableBackend)

	// add everything to the config for haproxy
	route := Route{
		Name:           newRoute.Name,
		StableFrontend: stableFrontend,
		StableBackend:  stableBackend,
	}
	c.Routes = append(c.Routes, &route)

	feCollection := []*Frontend{stableFrontend, frontendA, frontendB}
	beCollection := []*Backend{stableBackend, backendA, backendB}

	for _, fe := range feCollection {
		c.Frontends = append(c.Frontends, fe)
	}

	for _, be := range beCollection {
		c.Backends = append(c.Backends, be)
	}
	return nil
}

// deletes a route, cascading down the structure and remove all underpinning
// frontends, backends and servers.
func (c *Config) DeleteRoute(name string) bool {

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	result := false

	// first find the route
	for i, rt := range c.Routes {
		if rt.Name == name {

			// get the servers so we can delete the frontends/backends the Unix sockets are pointing to
			servers := rt.StableBackend.BackendServers

			for _, srv := range servers {
				for _, fe := range c.Frontends {
					if fe.UnixSock == srv.UnixSock {
						c.DeleteFrontend(fe.Name)
						c.DeleteBackend(fe.DefaultBackend)
					}
				}
			}

			c.DeleteFrontend(rt.StableFrontend.Name)
			c.DeleteBackend(rt.StableBackend.Name)
			c.Routes = append(c.Routes[:i], c.Routes[i+1:]...)
			result = true
			break
		}
	}
	return result
}

// gets a route
func (c *Config) GetRoute(name string) *Route {

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	var result *Route

	for _, rt := range c.Routes {
		if rt.Name == name {
			result = rt
			break

		}
	}
	return result
}

func defaultSocketProxyServer(name string, group string, weight int) *BackendServer {
	return &BackendServer{
		Name:          name + "_srv_" + group,
		Host:          "",
		Port:          0,
		UnixSock:      "/tmp/" + name + "_srv_" + group + ".sock",
		Weight:        weight,
		MaxConn:       1000,
		Check:         false,
		CheckInterval: 10,
	}
}

func defaultServer(name string, group string, weight int, host string, port int) *BackendServer {
	return &BackendServer{
		Name:          name + "_srv_" + group,
		Host:          host,
		Port:          port,
		UnixSock:      "",
		Weight:        weight,
		MaxConn:       1000,
		Check:         false,
		CheckInterval: 10,
	}
}

func defaultBackend(name string, group string, mode string, proxy bool, servers []*BackendServer) *Backend {
	var postfix string
	if group == "" {
		postfix = "_be"
	} else {
		postfix = "_be_" + group
	}
	return &Backend{
		Name:           name + postfix,
		Mode:           mode,
		BackendServers: servers,
		Options:        ProxyOptions{},
		ProxyMode:      proxy,
	}
}

func defaultFrontend(name string, mode string, port int, backend *Backend) *Frontend {

	return &Frontend{
		Name:           name + "_fe",
		Mode:           mode,
		BindPort:       port,
		BindIp:         "0.0.0.0",
		Options:        ProxyOptions{},
		DefaultBackend: backend.Name,
		ACLs:           []*ACL{},
		HttpSpikeLimit: SpikeLimit{},
		TcpSpikeLimit:  SpikeLimit{},
	}

}

func defaultSocketProxyFrontend(name string, group string, mode string, socket string, backend *Backend) *Frontend {

	return &Frontend{
		Name:           name + "_fe_" + group,
		Mode:           mode,
		UnixSock:       socket,
		SockProtocol:   "accept-proxy",
		Options:        ProxyOptions{},
		DefaultBackend: backend.Name,
		ACLs:           []*ACL{},
		HttpSpikeLimit: SpikeLimit{},
		TcpSpikeLimit:  SpikeLimit{},
	}

}
