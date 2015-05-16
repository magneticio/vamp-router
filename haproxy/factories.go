package haproxy

import (
	"github.com/magneticio/vamp-router/tools"
)

const (
	MAX_SOCKET_LENGTH = 103
)

// creates a Frontend object
func (c *Config) frontendFactory(name string, mode string, port int, ssl bool, filter []*Filter, backend *Backend) *Frontend {

	return &Frontend{
		Name:           name,
		Mode:           mode,
		BindPort:       port,
		BindIp:         "0.0.0.0",
		BindSsl:		ssl,
		Options:        ProxyOptions{},
		DefaultBackend: backend.Name,
		Filters:        filter,
		HttpQuota:      Quota{},
		TcpQuota:       Quota{},
	}
}

// creates a Backend object
func (c *Config) backendFactory(name string, mode string, proxy bool, servers []*ServerDetail) *Backend {

	return &Backend{
		Name:      name,
		Mode:      mode,
		Servers:   servers,
		Options:   ProxyOptions{},
		ProxyMode: proxy,
	}
}

// creates a ServerDetail object
func (c *Config) serverFactory(name string, weight int, host string, port int) *ServerDetail {
	return &ServerDetail{
		Name:          name,
		Host:          host,
		Port:          port,
		UnixSock:      "",
		Weight:        weight,
		MaxConn:       1000,
		Check:         false,
		CheckInterval: 10,
	}
}

// creates a Frontend object that listen on a Unix socket
func (c *Config) socketFrontendFactory(name string, mode string, socket string, backend *Backend) *Frontend {

	return &Frontend{
		Name:           name,
		Mode:           mode,
		UnixSock:       socket,
		SockProtocol:   "accept-proxy",
		Options:        ProxyOptions{},
		DefaultBackend: backend.Name,
		Filters:        []*Filter{},
		HttpQuota:      Quota{},
		TcpQuota:       Quota{},
	}
}

// creates a ServerDetail object that sends traffic to a Unix socket
func (c *Config) socketServerFactory(name string, weight int) *ServerDetail {

	return &ServerDetail{
		Name:          name,
		Host:          "",
		Port:          0,
		UnixSock:      compileSocketName(c.WorkingDir+"/sockets/", name, ".sock"),
		Weight:        weight,
		MaxConn:       1000,
		Check:         false,
		CheckInterval: 10,
	}
}

func compileSocketName(prefix string, base string, postfix string) string {

	if len(base) == 0 {
		return (prefix + tools.GetMD5Hash(tools.GetUUID()) + postfix)
	}
	return (prefix + tools.GetMD5Hash(base) + postfix)
}
