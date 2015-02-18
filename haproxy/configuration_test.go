package haproxy

import (
	"os"
	"testing"
)

const (
	TEMPLATE_FILE = "../configuration/templates/haproxy_config.template"
	CONFIG_FILE   = "/tmp/vamp_lb_test.cfg"
	EXAMPLE       = "../test/test_config1.json"
	JSON_FILE     = "/tmp/vamp_lb_test.json"
	PID_FILE      = "/tmp/vamp_lb_test.pid"
)

var (
	haConfig = Config{TemplateFile: TEMPLATE_FILE, ConfigFile: CONFIG_FILE, JsonFile: JSON_FILE, PidFile: PID_FILE}
)


// Runtime

func TestConfiguration_GetConfigFromDisk(t *testing.T) {
	err := haConfig.GetConfigFromDisk(EXAMPLE)
	if err != nil {
		t.Errorf("err: %v", err)
	}

	err = haConfig.GetConfigFromDisk("/this_is_really_something_wrong")
	if err == nil {
		t.Errorf("Expected an error")
	}
}

func TestConfiguration_SetWeight(t *testing.T) {
	err := haConfig.SetWeight("test_be_1", "test_be_1_a", 20)
	if err != nil {
		t.Errorf("err: %v", err)
	}
}


// Frontends

func TestConfiguration_GetFrontends(t *testing.T) {
	result := haConfig.GetFrontends()
	if result[0].Name != "test_fe_1" {
		t.Errorf("Failed to get frontends array")
	}

}

func TestConfiguration_GetFrontend(t *testing.T) {
	result := haConfig.GetFrontend("test_fe_1")
	if result.Name != "test_fe_1" {
		t.Errorf("Failed to get frontend")
	}

}

func TestConfiguration_AddFrontend(t *testing.T) {

	fe := Frontend{Name: "my_test_frontend", Mode: "http", DefaultBackend: "test_be_1"}
	err := haConfig.AddFrontend(&fe)
	if err != nil {
		t.Errorf("Failed to add frontend")
	}
	if haConfig.Frontends[3].Name != "my_test_frontend" {
		t.Errorf("Failed to add frontend")
	}
}

func TestConfiguration_DeleteFrontend(t *testing.T) {

	if err := haConfig.DeleteFrontend("test_fe_2"); err != nil {
		t.Errorf("Failed to remove frontend")
	}

	if err := haConfig.DeleteFrontend("non_existing_backend"); err == nil {
		t.Errorf("Backend should not be removed")
	}
}

func TestConfiguration_GetFilters(t *testing.T) {

	filters := haConfig.GetFilters("test_fe_1")
	if filters[0].Name != "uses_internetexplorer" {
		t.Errorf("Could not retrieve Filter")
	}
}

func TestConfiguration_AddFilter(t *testing.T) {

	filter := Filter{Name: "uses_firefox",Condition: "hdr_sub(user-agent) Mozilla", Destination: "test_be_1_b"}
	err := haConfig.AddFilter("test_fe_1", &filter)
	if err != nil {
		t.Errorf("Could not add Filter")
	}
	if haConfig.Frontends[0].Filters[1].Name != "uses_firefox" {
		t.Errorf("Could not add Filter")
	}
}

// Backends

func TestConfiguration_DeleteBackend(t *testing.T) {


	if err := haConfig.DeleteBackend("test_be_1"); err == nil {
		t.Errorf("Backend should not be removed because it is still in use")
	}

	if err := haConfig.DeleteBackend("deletable_backend"); err != nil {
		t.Errorf("Could not delete backend that should be deletable")
	}

	if err := haConfig.DeleteFrontend("non_existing_backend"); err == nil {
		t.Errorf("Backend should not be removed")
	}
}



// Rendering & Persisting

func TestConfiguration_Render(t *testing.T) {
	err := haConfig.Render()
	if err != nil {
		t.Errorf("err: %v", err)
	}
}

func TestConfiguration_Persist(t *testing.T) {
	err := haConfig.Persist()
	if err != nil {
		t.Errorf("err: %v", err)
	}
	os.Remove(CONFIG_FILE)
	os.Remove(JSON_FILE)
}

func TestConfiguration_RenderAndPersist(t *testing.T) {
	err := haConfig.RenderAndPersist()
	if err != nil {
		t.Errorf("err: %v", err)
	}
	os.Remove(CONFIG_FILE)
	os.Remove(JSON_FILE)
}
