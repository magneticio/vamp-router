package haproxy

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

const (
	// TEMPLATE_FILE = "../configuration/templates/haproxy_config.template"
	// CONFIG_FILE   = "/tmp/vamp_lb_test.cfg"
	// EXAMPLE       = "../test/test_config1.json"
	// JSON_FILE     = "/tmp/vamp_lb_test.json"
	// PID_FILE      = "/tmp/vamp_lb_test.pid"
	ROUTE_JSON = "../test/test_route.json"
)

func TestConfiguration_GetRoutes(t *testing.T) {

	routes := haConfig.GetRoutes()
	if routes[0].Name != "test_route_1" {
		t.Errorf("Failed to get all frontends")
	}
}

func TestConfiguration_GetRoute(t *testing.T) {

	route, err := haConfig.GetRoute("test_route_1")
	if route.Name != "test_route_1" && err == nil {
		t.Errorf("Failed to get frontend")
	}

	_, err = haConfig.GetRoute("non_existent_route")
	if err == nil {
		t.Errorf("Should return nil on non existent route")
	}

}

func TestConfiguration_PostRoute(t *testing.T) {
	j, err := ioutil.ReadFile(ROUTE_JSON)
	var route *Route
	err = json.Unmarshal(j, &route)

	if haConfig.AddRoute(route) != nil || err != nil {
		t.Errorf("Failed to add route")
	}
}

func TestConfiguration_UpdateRoute(t *testing.T) {
	j, err := ioutil.ReadFile(ROUTE_JSON)
	var route *Route
	err = json.Unmarshal(j, &route)
	route.Protocol = "tcp"

	err = haConfig.UpdateRoute("test_route_2", route)
	if err != nil {
		t.Errorf("Failed to update route")
	}

	route, err = haConfig.GetRoute("test_route_2")
	if err != nil && route.Protocol != "tcp" {
		t.Errorf("Failed to update route")
	}
}

func TestConfiguration_DeleteRoute(t *testing.T) {

	err := haConfig.DeleteRoute("test_route_2")
	if err != nil {
		t.Errorf("Failed to delete route")
	}

	err = haConfig.DeleteRoute("non_existent_route")
	if err == nil {
		t.Errorf("Should return nil on non existent route")
	}
}
