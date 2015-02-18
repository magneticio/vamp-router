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
	GROUP_JSON = "../test/test_group.json"
	SERVER_JSON = "../test/test_server1.json"
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


func TestConfiguration_GetRouteGroups(t *testing.T) {

	groups, err := haConfig.GetRouteGroups("test_route_1")
	if groups[0].Name != "group_a" || err != nil {
		t.Errorf("Failed to get groups")
	}
}

func TestConfiguration_GetRouteGroup(t *testing.T) {

	group, err := haConfig.GetRouteGroup("test_route_1","group_a")
	if group.Name != "group_a" || err != nil {
		t.Errorf("Failed to get group")
	}

	if _, err = haConfig.GetRouteGroup("non_existent_route","group_a"); err == nil {
				t.Errorf("Should return nil on non existent route")
	}
		if _, err = haConfig.GetRouteGroup("test_route_1","non_existent_group"); err == nil {
				t.Errorf("Should return nil on non existent group")
	}

}

func TestConfiguration_PostRouteGroup(t *testing.T) {
	j, err := ioutil.ReadFile(GROUP_JSON)
	var group *Group
	err = json.Unmarshal(j, &group)

	route := "test_route_1"

	if haConfig.AddRouteGroup(route,group) != nil || err != nil {
		t.Errorf("Failed to add route")
	}

}

func TestConfiguration_DeleteRouteGroup(t *testing.T) {

	route := "test_route_1"

	err := haConfig.DeleteRouteGroup(route,"group_c")
	if err != nil {
		t.Errorf("Failed to delete route")
	}

	err = haConfig.DeleteRoute("non_existent_group")
	if err == nil {
		t.Errorf("Should return nil on non existent group")
	}
}

func TestConfiguration_AddGroupServer(t *testing.T) {

	route := "test_route_1"
	group := "group_a"

	j, err := ioutil.ReadFile(GROUP_JSON)
	var server Server

	err = json.Unmarshal(j, &server)

	err = haConfig.AddGroupServer(route,group,&server)
	if err != nil {
		t.Errorf("Failed to create server")
	}

	err = haConfig.AddGroupServer(route,"non_existent_group",&server)
	if err == nil {
		t.Errorf("Should return error on non existent group")
	}
	err = haConfig.AddGroupServer("non_existent_route",group,&server)
	if err == nil {
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

	err := haConfig.DeleteGroupServer(route,group,server)
	if err != nil {
		t.Errorf("Failed to delete server")
	}

	err = haConfig.DeleteGroupServer(route,group,"non_existent_server")
	if err == nil {
		t.Errorf("Should return nil on non existent server")
	}
}

func TestConfiguration_UpdateGroupServer(t *testing.T) {

	j, err := ioutil.ReadFile(SERVER_JSON)

	var server *Server
	err = json.Unmarshal(j, &server)
	serverToUpdate := "server_to_be_updated"
	server.Port = 1234
	routeName := "test_route_2"
	groupName := "group_to_be_updated"

	err = haConfig.UpdateGroupServer(routeName,groupName,serverToUpdate,server)
	if err != nil {
		t.Errorf(err.Error())
	}
	
	server, err = haConfig.GetGroupServer(routeName, groupName, server.Name)
	if err != nil && server.Port != 1234 {
		t.Errorf(err.Error()	)
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
