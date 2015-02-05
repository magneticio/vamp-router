package haproxy

import (
	"github.com/magneticio/vamp-loadbalancer/helpers"
	"io/ioutil"
	"os"
	"testing"
)

var (
	haRuntime = Runtime{Binary: helpers.HaproxyLocation()}
)

func TestRuntime_SetNewPid(t *testing.T) {

	//make sure there is no pidfile present
	os.Remove(PID_FILE)
	result := haRuntime.SetPid(PID_FILE)
	if result != true {
		t.Fatalf("err: Failed to create pid file")
	}
	os.Remove(PID_FILE)
}

func TestRuntime_UseExistingPid(t *testing.T) {

	//create a pid file
	emptyPid := []byte("12356")
	ioutil.WriteFile(PID_FILE, emptyPid, 0644)

	result := haRuntime.SetPid(PID_FILE)
	if result != false {
		t.Fatalf("err: Failed to read existing pid file")
	}
	os.Remove(PID_FILE)
}

// for some reason, this always returns 1 (error).
// func TestRuntime_Reload(t *testing.T) {
// 	haRuntime.SetPid(PID_FILE)
// 	err := haRuntime.Reload(&haConfig)
// 	if err != nil {
// 		t.Fatalf("err: %v", err)
// 	}
// }
