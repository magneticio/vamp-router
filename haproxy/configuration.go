package haproxy

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"text/template"
	"errors"
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


// the transactions methods are kept separate so we can chain an arbitrary set of operations
// on the Config object within one transaction. Alas, this burdons the developer with extra housekeeping
// but gives you more control over the flow of mutations and reads without risking deadlocks or duplicating
// locks and unlocks inside of methods.
func (c *Config) BeginWriteTrans() {
	c.Mutex.Lock()
}

func (c *Config) EndWriteTrans() {
	c.Mutex.Unlock()
}

func (c *Config) BeginReadTrans() {
	c.Mutex.RLock()
}

func (c *Config) EndReadTrans() {
	c.Mutex.RUnlock()
}


// gets all frontends
func (c *Config) GetFrontends() []*Frontend {

	return c.Frontends
}

// gets a frontend
func (c *Config) GetFrontend(name string) *Frontend {

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

	c.Frontends = append(c.Frontends, frontend)
	return nil

}

// deletes a frontend
func (c *Config) DeleteFrontend(name string) error {

	for i, fe := range c.Frontends {
		if fe.Name == name {
			c.Frontends = append(c.Frontends[:i], c.Frontends[i+1:]...)
			return nil
		}
	}
	return errors.New("no such frontend found")
}

// get the filters from a frontend
func (c *Config) GetFilters(frontend string) []*Filter {

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

	for _, fe := range c.Frontends {
		if fe.Name == frontend {
			fe.Filters = append(fe.Filters, filter)
		}
	}
	return nil
}

// delete a Filter from a frontend
func (c *Config) DeleteFilter(frontendName string, filterName string) bool {

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

	return c.Backends
}

// adds a frontend
func (c *Config) AddBackend(backend *Backend) error {

	c.Backends = append(c.Backends, backend)
	return nil

}

/* Deleting a backend is tricky. Frontends have a default backend. Removing that backend and then reloading
the configuration will crash Haproxy. This means some extra protection is put into this method to check
if this backend is still used. If not, it can be deleted.
*/
func (c *Config) DeleteBackend(name string) error {

	if err := c.BackendUsed(name); err != nil {
		return err 
	} else {
		for i, be := range c.Backends {
			if be.Name == name {
				c.Backends = append(c.Backends[:i], c.Backends[i+1:]...)
				return nil
			}
		}
		return errors.New("no such backend found")
	} 
}


func (c *Config) GetServer(backendName string, serverName string) *ServerDetail {

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
func (c *Config) AddServer(backendName string, server *ServerDetail) error {

	for _, be := range c.Backends {
		if be.Name == backendName {
			be.Servers = append(be.Servers, server)
			return nil
		}
	}
	return errors.New("No such backend found")
}

func (c *Config) DeleteServer(backendName string, serverName string) error {
	for _, be := range c.Backends {
		if be.Name == backendName {
			for i, srv := range be.Servers {
				if srv.Name == serverName {
					be.Servers = append(be.Servers[:i], be.Servers[i+1:]...)
					return nil
				}
			}
		}
	}
	return errors.New("no such server found")
}

// gets all servers of a specific backend
func (c *Config) GetServers(backendName string) []*ServerDetail {

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

// helper function to check if a Frontend exists
func (c *Config) FrontendExists(name string) bool {
	
	for _,frontend := range c.Frontends {
		if frontend.Name == name {
			return true
		}
	}
	return false
}

// helper function to check if a Backend exists
func (c *Config) BackendExists(name string) bool {
	
	for _,backend := range c.Backends {
		if backend.Name == name {
			return true
		}
	}
	return false
}


// helper function to check if a Backend is used by a Frontend as a default backend or a filter destination
func (c *Config) BackendUsed(name string) error {

	if c.BackendExists(name) {
		for _,frontend := range c.Frontends {
			if frontend.DefaultBackend == name {
				return errors.New("Backend still in use by: " + frontend.Name)
			}
			for _,filter := range frontend.Filters {
				if filter.Destination == name {
					return errors.New("Backend still in use by: " + frontend.Name + ".Filters." + filter.Name) 
				}
			}
		}

	}
	return nil
}


// helper function to check if a Route exists
func (c *Config) RouteExists(name string) bool {
	
	for _,route := range c.Routes {
		if route.Name == name {
			return true
		}
	}
	return false
}

// helper function to check if a Group exists
func (c *Config) GroupExists(routeName string, groupName string) bool {
	
	for _, rt := range c.Routes {
		if rt.Name == routeName {
			for _, grp := range rt.Groups {
				if grp.Name == groupName {
					return true
				}
			}
		}
	}
	return false
}

// helper function to check if a Server exists in a specific Group
func (c *Config) ServerExists(routeName string, groupName string, serverName string) bool {
	
	for _, rt := range c.Routes {
		if rt.Name == routeName {
			for _, grp := range rt.Groups {
				if grp.Name == groupName {
					for _,server := range grp.Servers {
						if server.Name == serverName{
							return true
						}
					}
				}
			}
		}
	}
	return false
}

// helper function to create a Backend or Frontend name based on a Route and Group
func GroupName(routeName string, groupName string) string {
	return routeName + "." + groupName
}
func RouteName(routeName string, groupName string) string {
	return routeName + "." + groupName
}

func BackendName(routeName string, groupName string) string {
	return routeName + "." + groupName
}

func FrontendName(routeName string, groupName string) string {
	return routeName + "." + groupName
}

func ServerName(routeName string, groupName string) string {
	return routeName + "." + groupName
}
