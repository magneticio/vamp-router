package haproxy

import (
	"encoding/json"
	"io/ioutil"
	"sync"
	"os"
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


// get the acls from a frontend
func (c *Config) GetAcls(frontend string) [] *ACL {

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()

	var acls [] *ACL

	for _, fe := range c.Frontends {
		if fe.Name == frontend {
			 acls = fe.ACLs

			}
		}

return acls

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

	err =	c.Persist()	
	if err != nil {
		return err
	}

	return nil 	
}
