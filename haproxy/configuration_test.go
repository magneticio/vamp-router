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

func TestConfiguration_GetConfigFromDisk(t *testing.T) {
	err := haConfig.GetConfigFromDisk(EXAMPLE)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}

func TestConfiguration_SetWeight(t *testing.T) {
	err := haConfig.SetWeight("test_be_1", "test_be_1_a", 20)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}

func TestConfiguration_Render(t *testing.T) {
	err := haConfig.Render()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
}

func TestConfiguration_Persist(t *testing.T) {
	err := haConfig.Persist()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	os.Remove(CONFIG_FILE)
	os.Remove(JSON_FILE)
}

func TestConfiguration_RenderAndPersist(t *testing.T) {
	err := haConfig.RenderAndPersist()
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	os.Remove(CONFIG_FILE)
	os.Remove(JSON_FILE)
}
