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

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	for _, be := range c.Backends {
		if be.Name == backend {
			for _, srv := range be.BackendServers {
				if srv.Name == server {
					srv.Weight = weight
				}
			}
		}
	}

	err := c.RenderAndPersist()
	if err != nil {
		return err
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

// gets a frontend
func (c *Config) AddFrontend(frontend *Frontend) error {

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	c.Frontends = append(c.Frontends, frontend)

	return nil

}

// deletes a frontend
func (c *Config) DeleteFrontend(name string) bool {

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
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

// get the acls from a frontend
func (c *Config) GetAcls(frontend string) []*ACL {

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	var acls []*ACL

	for _, fe := range c.Frontends {
		if fe.Name == frontend {
			acls = fe.ACLs

		}
	}
	return acls
}

// get the acls from a frontend
func (c *Config) AddAcl(frontend string, acl *ACL) {

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	for _, fe := range c.Frontends {
		if fe.Name == frontend {
			fe.ACLs = append(fe.ACLs, acl)
		}
	}
}

// delete an ACL from a frontend
func (c *Config) DeleteAcl(frontendName string, aclName string) bool {

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	result := false

	for _, fe := range c.Frontends {
		if fe.Name == frontendName {
			for i, acl := range fe.ACLs {
				if acl.Name == aclName {
					fe.ACLs = append(fe.ACLs[:i], fe.ACLs[i+1:]...)
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

func (c *Config) GetServer(backendName string, serverName string) *BackendServer {

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	var result *BackendServer

	for _, be := range c.Backends {
		if be.Name == backendName {
			for _, srv := range be.BackendServers {
				if srv.Name == serverName {
					result = srv
					break
				}
			}
		}
	}
	return result
}

func (c *Config) DeleteServer(backendName string, serverName string) bool {

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	result := false

	for _, be := range c.Backends {
		if be.Name == backendName {
			for i, srv := range be.BackendServers {
				if srv.Name == serverName {
					be.BackendServers = append(be.BackendServers[:i], be.BackendServers[i+1:]...)
					result = true
					break
				}
			}
		}
	}
	return result
}

// gets all servers of a specific backend
func (c *Config) GetServers(backendName string) []*BackendServer {

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	var result []*BackendServer

	for _, be := range c.Backends {
		if be.Name == backendName {
			result = be.BackendServers
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

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

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
