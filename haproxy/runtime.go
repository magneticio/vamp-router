package haproxy

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/magneticio/vamp-loadbalancer/tools"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// returns an error if the file was already there
func (r *Runtime) SetPid(pidfile string) error {

	//Create and empty pid file on the specified location, if not already there
	if _, err := os.Stat(pidfile); err != nil {
		emptyPid := []byte("")
		ioutil.WriteFile(pidfile, emptyPid, 0644)
		return nil
	}
	return errors.New("file already there")
}

// Reload runtime with configuration
func (r *Runtime) Reload(c *Config) error {

	pid, err := ioutil.ReadFile(c.PidFile)
	if err != nil {
		return err
	}

	/*  Setup all the command line parameters so we get an executable similar to
	    /usr/local/bin/haproxy -f resources/haproxy_new.cfg -p resources/haproxy-private.pid -sf 1234

	*/
	arg0 := "-f"
	arg1 := c.ConfigFile
	arg2 := "-p"
	arg3 := c.PidFile
	arg4 := "-D"
	arg5 := "-sf"
	arg6 := strings.Trim(string(pid), "\n")
	var cmd *exec.Cmd

	// fmt.Println(r.Binary + " " + arg0 + " " + arg1 + " " + arg2 + " " + arg3 + " " + arg4 + " " + arg5 + " " + arg6)
	// If this is the first run, the PID value will be empty, otherwise it will be > 0
	if len(arg6) > 0 {
		cmd = exec.Command(r.Binary, arg0, arg1, arg2, arg3, arg4, arg5, arg6)
	} else {
		cmd = exec.Command(r.Binary, arg0, arg1, arg2, arg3, arg4)
	}

	var out bytes.Buffer
	cmd.Stdout = &out

	cmdErr := cmd.Run()
	if cmdErr != nil {
		return cmdErr
	}

	return nil
}

// Sets the weight of a backend
func (r *Runtime) SetWeight(backend string, server string, weight int) (string, error) {

	result, err := r.cmd("set weight " + backend + "/" + server + " " + strconv.Itoa(weight) + "\n")

	if err != nil {
		return "", err
	} else {
		return result, nil
	}

}

// Adds an ACL.
// We need to match a frontend name to an id. This is somewhat awkard.
// func (r *Runtime) SetAcl(frontend string, acl string, pattern string) (string, error) {

// 	result, err := r.cmd("add acl " + acl + pattern)

// 	if err != nil {
// 		return "", err
// 	} else {
// 		return result, nil
// 	}
// }

// Gets basic info on haproxy process
func (r *Runtime) GetInfo() (Info, error) {
	var Info Info
	result, err := r.cmd("show info \n")
	if err != nil {
		return Info, err
	} else {
		result, err := tools.MultiLineToJson(result)
		if err != nil {
			return Info, err
		} else {
			err := json.Unmarshal([]byte(result), &Info)
			if err != nil {
				return Info, err
			} else {
				return Info, nil
			}
		}
	}

}

/* get the basic stats in CSV format

@parameter statsType takes the form of:
- all
- frontend
- backend
*/
func (r *Runtime) GetStats(statsType string) ([]StatsGroup, error) {

	var Stats []StatsGroup
	var cmdString string

	switch statsType {
	case "all":
		cmdString = "show stat -1\n"
	case "backend":
		cmdString = "show stat -1 2 -1\n"
	case "frontend":
		cmdString = "show stat -1 1 -1\n"
	case "server":
		cmdString = "show stat -1 4 -1\n"
	}

	result, err := r.cmd(cmdString)
	if err != nil {
		return Stats, err
	} else {
		result, err := tools.CsvToJson(strings.Trim(result, "# "))
		if err != nil {
			return Stats, err
		} else {
			err := json.Unmarshal([]byte(result), &Stats)
			if err != nil {
				return Stats, err
			} else {
				return Stats, nil
			}
		}

	}
}

// Executes a arbitrary HAproxy command on the unix socket
func (r *Runtime) cmd(cmd string) (string, error) {

	// connect to haproxy
	conn, err_conn := net.Dial("unix", "/tmp/haproxy.stats.sock")
	defer conn.Close()

	if err_conn != nil {
		return "", errors.New("Unable to connect to Haproxy socket")
	} else {

		fmt.Fprint(conn, cmd)

		response := ""

		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			response += (scanner.Text() + "\n")
		}
		if err := scanner.Err(); err != nil {
			return "", err
		} else {
			return response, nil
		}

	}
}
