package api

import (
	"github.com/magneticio/vamp-router/haproxy"
	"github.com/magneticio/vamp-router/helpers"
	"github.com/magneticio/vamp-router/logging"
	gologger "github.com/op/go-logging"
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
	haConfig  = haproxy.Config{TemplateFile: TEMPLATE_FILE, ConfigFile: CONFIG_FILE, JsonFile: JSON_FILE, PidFile: PID_FILE}
	haRuntime = haproxy.Runtime{Binary: helpers.HaproxyLocation()}
	log       = logging.ConfigureLog(LOG_PATH)
)

func TestApi_GetConfig(t *testing.T) {

	api.CreateApi(port, &haConfig, &haRuntime, log, sseBroker).Run("0.0.0.0:" + strconv.Itoa(port))

	req, _ := http.NewRequest("GET", "/v1/config", nil)
	w := httptest.NewRecorder()

	v1.ServeHTTP(w, req)

	if w.Body.String() != "{\"frontends\":\"[]\"}\n" {
		t.Errorf("Response should be {\"foo\":\"bar\"}, was: %s", w.Body.String())
	}

}
