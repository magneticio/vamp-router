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
	ROUTE_JSON  = "../test/test_route.json"
	GROUP_JSON  = "../test/test_group.json"
	SERVER_JSON = "../test/test_server1.json"
)

func TestConfiguration_GetRoutes(t *testing.T) {

	routes := haConfig.GetRoutes()
	if routes[0].Name != "test_route_1" {
		t.Errorf("Failed to get all frontends")
	}
}

func TestConfiguration_GetRoute(t *testing.T) {

	if route, err := haConfig.GetRoute("test_route_1"); route.Name != "test_route_1" && err == nil {
		t.Errorf("Failed to get frontend")
	}

	if _, err := haConfig.GetRoute("non_existent_route"); err == nil {
		t.Errorf("Should return nil on non existent route")
	}

}

func TestConfiguration_AddRoute(t *testing.T) {
	j, _ := ioutil.ReadFile(ROUTE_JSON)
	var route *Route
	_ = json.Unmarshal(j, &route)

	if haConfig.AddRoute(route) != nil {
		t.Errorf("Failed to add route")
	}

	if haConfig.AddRoute(route) == nil {
		t.Errorf("Adding should fail when a route already exists")
	}
}

func TestConfiguration_UpdateRoute(t *testing.T) {
	j, _ := ioutil.ReadFile(ROUTE_JSON)
	var route *Route
	if err := json.Unmarshal(j, &route); err != nil {
		t.Errorf(err.Error())
	}
	route.Protocol = "tcp"

	if err := haConfig.UpdateRoute("test_route_2", route); err != nil {
		t.Errorf(err.Error())
	}

	if route, err := haConfig.GetRoute("test_route_2"); err != nil && route.Protocol != "tcp" {
		t.Errorf("Failed to update route")
	}
}

func TestConfiguration_GetRouteGroups(t *testing.T) {

	if groups, err := haConfig.GetRouteGroups("test_route_1"); groups[0].Name != "group_a" || err != nil {
		t.Errorf("Failed to get groups")
	}

	if _, err := haConfig.GetRouteGroups("non_existent_group"); err == nil {
		t.Errorf("Should return nil on non existent group")
	}
}

func TestConfiguration_GetRouteGroup(t *testing.T) {

	if group, err := haConfig.GetRouteGroup("test_route_1", "group_a"); group.Name != "group_a" || err != nil {
		t.Errorf("Failed to get group")
	}

	if _, err := haConfig.GetRouteGroup("non_existent_route", "group_a"); err == nil {
		t.Errorf("Should return nil on non existent route")
	}
	if _, err := haConfig.GetRouteGroup("test_route_1", "non_existent_group"); err == nil {
		t.Errorf("Should return nil on non existent group")
	}

}

func TestConfiguration_AddRouteGroups(t *testing.T) {
	j, _ := ioutil.ReadFile(GROUP_JSON)
	var groups []*Group
	_ = json.Unmarshal(j, &groups)

	route := "test_route_1"

	if haConfig.AddRouteGroups(route, groups) != nil {
		t.Errorf("Failed to add route")
	}

	if haConfig.AddRouteGroups(route, groups) == nil {
		t.Errorf("Adding should fail when a group already exists")
	}

	if haConfig.AddRouteGroups("non_existent_group", groups) == nil {
		t.Errorf("Should return nil on non existent route")
	}

}

func TestConfiguration_UpdateRouteGroup(t *testing.T) {
	j, _ := ioutil.ReadFile(GROUP_JSON)
	var groups []*Group
	_ = json.Unmarshal(j, &groups)

	group := groups[0]
	group.Weight = 1

	if err := haConfig.UpdateRouteGroup("test_route_1", groups[0].Name, group); err != nil {
		t.Errorf(err.Error())
	}

}

func TestConfiguration_DeleteRouteGroup(t *testing.T) {

	route := "test_route_1"

	if err := haConfig.DeleteRouteGroup(route, "group_c"); err != nil {
		t.Errorf("Failed to delete route")
	}

	if haConfig.DeleteRouteGroup("non_existent_route", "group_a") == nil {
		t.Errorf("Should return nil on non existent route")
	}

	if haConfig.DeleteRouteGroup(route, "non_existent_group") == nil {
		t.Errorf("Should return nil on non existent group")
	}
}

func TestConfiguration_GetGroupServers(t *testing.T) {

	if servers, err := haConfig.GetGroupServers("test_route_1", "group_a"); err != nil {
		t.Errorf("Failed to get servers")
	} else {
		if servers[0].Name != "paas.55f73f0d-6087-4964-a70e-b1ca1d5b24cd" {
			t.Errorf("Failed to get servers")
		}
	}

	if _, err := haConfig.GetGroupServers("non_existent_route", "group_a"); err == nil {
		t.Errorf("Should return nil on non existent route")
	}

	if _, err := haConfig.GetGroupServers("test_route_1", "non_existent_group"); err == nil {
		t.Errorf("Should return nil on non existent group")
	}

}

func TestConfiguration_GetGroupServer(t *testing.T) {

	if _, err := haConfig.GetGroupServer("test_route_1", "group_a", "paas.55f73f0d-6087-4964-a70e-b1ca1d5b24cd"); err != nil {
		t.Errorf("Failed to get server")
	}

	if _, err := haConfig.GetGroupServer("test_route_1", "group_a", "non_existent_server"); err == nil {
		t.Errorf("Should return nil on non existent server")
	}
}

func TestConfiguration_AddGroupServer(t *testing.T) {

	route := "test_route_1"
	group := "group_a"

	j, _ := ioutil.ReadFile(GROUP_JSON)
	var server Server

	_ = json.Unmarshal(j, &server)

	if err := haConfig.AddGroupServer(route, group, &server); err != nil {
		t.Errorf(err.Error())
	}

	if err := haConfig.AddGroupServer(route, group, &server); err == nil {
		t.Errorf("Adding should fail when a server already exists")
	}

	if err := haConfig.AddGroupServer(route, "non_existent_group", &server); err == nil {
		t.Errorf("Should return error on non existent group")
	}

	if err := haConfig.AddGroupServer("non_existent_route", group, &server); err == nil {
		t.Errorf("Should return error on non existent route")
	}

	// test should be activated when "exists" checking is in the haproxy package
	// server.Name = "paas.55f73f0d-6087-4964-a70e-b1ca1d5b24cd"
	// err = haConfig.AddGroupServer(route,group,&server)
	// if err == nil {
	// 	t.Errorf("Should return error on trying to create an already existing server")
	// }

}

func TestConfiguration_DeleteGroupServer(t *testing.T) {

	route := "test_route_1"
	group := "group_a"
	server := "paas.55f73f0d-6087-4964-a70e-b1ca1d5b24cd"

	if err := haConfig.DeleteGroupServer(route, group, server); err != nil {
		t.Errorf("Failed to delete server")
	}

	if err := haConfig.DeleteGroupServer(route, group, "non_existent_server"); err == nil {
		t.Errorf("Should return nil on non existent server")
	}
}

func TestConfiguration_UpdateGroupServer(t *testing.T) {

	j, _ := ioutil.ReadFile(SERVER_JSON)

	var server *Server
	_ = json.Unmarshal(j, &server)
	serverToUpdate := "server_to_be_updated"
	server.Port = 1234
	routeName := "test_route_2"
	groupName := "group_to_be_updated"

	if err := haConfig.UpdateGroupServer(routeName, groupName, serverToUpdate, server); err != nil {
		t.Errorf(err.Error())
	}

	if server, err := haConfig.GetGroupServer(routeName, groupName, server.Name); err != nil && server.Port != 1234 {
		t.Errorf(err.Error())
	}
}

func TestConfiguration_DeleteRoute(t *testing.T) {

	if err := haConfig.DeleteRoute("test_route_2"); err != nil {
		t.Errorf("Failed to delete route")
	}

	if err := haConfig.DeleteRoute("non_existent_route"); err == nil {
		t.Errorf("Should return nil on non existent route")
	}
}
