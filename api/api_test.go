package api

import (
	"github.com/gin-gonic/gin"
	"github.com/magneticio/vamp-loadbalancer/haproxy"
	"github.com/magneticio/vamp-loadbalancer/helpers"
	"testing"
)

const (
	TEMPLATE_FILE = "../configuration/templates/haproxy_config.template"
	CONFIG_FILE   = "/tmp/vamp_lb_test.cfg"
	EXAMPLE       = "../test/test_config1.json"
	JSON_FILE     = "/tmp/vamp_lb_test.json"
	PID_FILE      = "/tmp/vamp_lb_test.pid"
	LOG_PATH      = "/tmp/vamp_lb_test.log"
)

var (
	haConfig  = haproxy.Config{TemplateFile: TEMPLATE_FILE, ConfigFile: CONFIG_FILE, JsonFile: JSON_FILE, PidFile: PID_FILE}
	haRuntime = haproxy.Runtime{Binary: helpers.HaproxyLocation()}
	c         *gin.Context
)

func TestApi_GetConfig(t *testing.T) {
	c.Set("haConfig", &haConfig)
	// config := c.MustGet("haConfig").(*haproxy.Config)
	// GetConfig(c)
	// if err != nil {
	// 	t.Fatalf("err: %v", err)
	// }
}
