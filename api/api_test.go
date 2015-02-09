package api

import (
	"github.com/gin-gonic/gin"
	"github.com/magneticio/vamp-loadbalancer/haproxy"
	"github.com/magneticio/vamp-loadbalancer/helpers"
	"net/http"
	"net/http/httptest"
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
	haConfig  = *haproxy.Config{TemplateFile: TEMPLATE_FILE, ConfigFile: CONFIG_FILE, JsonFile: JSON_FILE, PidFile: PID_FILE}
	haRuntime = *haproxy.Runtime{Binary: helpers.HaproxyLocation()}
	r         = gin.New()
)

func init() {
	r.Use(HaproxyMiddleware(haConfig, haRuntime))
	r.Use(LoggerMiddleware(log.Logger))
	r.Use(gin.Recovery())
	r.Static("/www", "./www")
	v1 := r.Group("/v1")
}

func TestApi_GetConfig(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/config", nil)
	w := httptest.NewRecorder()

	v1.GET("/v1/config", GetFrontends)
	v1.ServeHTTP(w, req)

	if w.Body.String() != "{\"frontends\":\"[]\"}\n" {
		t.Errorf("Response should be {\"foo\":\"bar\"}, was: %s", w.Body.String())
	}

}
