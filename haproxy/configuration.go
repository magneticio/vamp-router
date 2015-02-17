package haproxy

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"text/template"
)

// Load a config from disk
func (c *Config) GetConfigFromDisk(file string) error {
	s, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(s, &c)
	if err != nil {
		return err
	}

	c.Mutex = new(sync.RWMutex)
	return err
}

// updates the weight of a server of a specific backend with a new weight
func (c *Config) SetWeight(backend string, server string, weight int) error {

	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	for _, be := range c.Backends {
		if be.Name == backend {
			for _, srv := range be.Servers {
				if srv.Name == server {
					srv.Weight = weight
				}
			}
		}
	}

	return nil
}

// gets all frontends
func (c *Config) GetFrontends() []*Frontend {

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	return c.Frontends
}

// gets a frontend
func (c *Config) GetFrontend(name string) *Frontend {

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	var result *Frontend

	for _, fe := range c.Frontends {
		if fe.Name == name {
			result = fe
			break

		}
	}
	return result
}

// adds a frontend
func (c *Config) AddFrontend(frontend *Frontend) error {

	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	c.Frontends = append(c.Frontends, frontend)

	return nil

}

// deletes a frontend
func (c *Config) DeleteFrontend(name string) bool {

	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	result := false

	for i, fe := range c.Frontends {
		if fe.Name == name {
			c.Frontends = append(c.Frontends[:i], c.Frontends[i+1:]...)
			result = true
			break
		}
	}
	return result
}

// get the filters from a frontend
func (c *Config) GetFilters(frontend string) []*Filter {

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	var filters []*Filter

	for _, fe := range c.Frontends {
		if fe.Name == frontend {
			filters = fe.Filters

		}
	}
	return filters
}

// set the filter on a frontend
func (c *Config) AddFilter(frontend string, filter *Filter) error {

	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	for _, fe := range c.Frontends {
		if fe.Name == frontend {
			fe.Filters = append(fe.Filters, filter)
		}
	}
	return nil
}

// delete a Filter from a frontend
func (c *Config) DeleteFilter(frontendName string, filterName string) bool {

	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	result := false

	for _, fe := range c.Frontends {
		if fe.Name == frontendName {
			for i, filter := range fe.Filters {
				if filter.Name == filterName {
					fe.Filters = append(fe.Filters[:i], fe.Filters[i+1:]...)
					result = true
					break
				}
			}
		}
	}
	return result
}

// gets a backend
func (c *Config) GetBackend(backend string) *Backend {

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	var result *Backend

	for _, be := range c.Backends {
		if be.Name == backend {
			result = be
			break

		}
	}
	return result
}

// gets all backends
func (c *Config) GetBackends() []*Backend {

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	return c.Backends
}

// deletes a backend
func (c *Config) DeleteBackend(name string) bool {

	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	result := false

	for i, be := range c.Backends {
		if be.Name == name {
			c.Backends = append(c.Backends[:i], c.Backends[i+1:]...)
			result = true
			break
		}
	}
	return result
}

func (c *Config) GetServer(backendName string, serverName string) *ServerDetail {

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	var result *ServerDetail

	for _, be := range c.Backends {
		if be.Name == backendName {
			for _, srv := range be.Servers {
				if srv.Name == serverName {
					result = srv
					break
				}
			}
		}
	}
	return result
}

// adds a Server
func (c *Config) AddServer(backendName string, server *ServerDetail) bool {

	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	result := false

	for _, be := range c.Backends {
		if be.Name == backendName {
			be.Servers = append(be.Servers, server)
			result = true
			break
		}
	}
	return result
}

func (c *Config) DeleteServer(backendName string, serverName string) bool {

	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	result := false

	for _, be := range c.Backends {
		if be.Name == backendName {
			for i, srv := range be.Servers {
				if srv.Name == serverName {
					be.Servers = append(be.Servers[:i], be.Servers[i+1:]...)
					result = true
					break
				}
			}
		}
	}
	return result
}

// gets all servers of a specific backend
func (c *Config) GetServers(backendName string) []*ServerDetail {

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	var result []*ServerDetail

	for _, be := range c.Backends {
		if be.Name == backendName {
			result = be.Servers
			break
		}
	}
	return result
}

// Render a config object to a HAproxy config file
func (c *Config) Render() error {

	// read the template
	f, err := ioutil.ReadFile(c.TemplateFile)
	if err != nil {
		return err
	}

	// create a file for the config
	fp, err := os.OpenFile(c.ConfigFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer fp.Close()

	// render the template
	t := template.Must(template.New(c.TemplateFile).Parse(string(f)))
	err = t.Execute(fp, &c)
	if err != nil {
		return err
	}

	return nil
}

// save the JSON config to disk
func (c *Config) Persist() error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(c.JsonFile, b, 0666)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) RenderAndPersist() error {

	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	err := c.Render()
	if err != nil {
		return err
	}

	err = c.Persist()
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) RouteExists(name string) bool {
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	
	for _,route := range c.Routes {
		if route.Name == name {
			return true
		}
	}
	return false
}
