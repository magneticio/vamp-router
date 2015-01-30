package logging

import (
  lumberjack "github.com/natefinch/lumberjack"
  gologging "github.com/op/go-logging"
  "io"
  "os"
)

type Logger struct {
  Log *gologging.Logger
}

func PrintLogo(version string) string {

var logo =`

██╗   ██╗ █████╗ ███╗   ███╗██████╗
██║   ██║██╔══██╗████╗ ████║██╔══██╗
██║   ██║███████║██╔████╔██║██████╔╝
╚██╗ ██╔╝██╔══██║██║╚██╔╝██║██╔═══╝
 ╚████╔╝ ██║  ██║██║ ╚═╝ ██║██║
  ╚═══╝  ╚═╝  ╚═╝╚═╝     ╚═╝╚═╝
                       loadbalancer
                       version ` + version + `
                       by magnetic.io
                                      `

return logo                       

}

func ConfigureLog(logPath string) *gologging.Logger {

  var log = gologging.MustGetLogger("vamp-loadbalancer")
  var backend *gologging.LogBackend
  var format = gologging.MustStringFormatter(
    "%{color}%{time:15:04:05.000} %{shortfunc} %{level:.4s} ==> %{color:reset} %{message}",
  )

  // mix in the Lumberjack logger so we can have rotation on log files
  if len(logPath) > 0 {
    backend = gologging.NewLogBackend(io.MultiWriter(&lumberjack.Logger{
      Filename:   logPath,
      MaxSize:    50, // megabytes
      MaxBackups: 2,  //days
      MaxAge:     14,
    }, os.Stdout), "", 0)
  }

  backendFormatter := gologging.NewBackendFormatter(backend, format)
  gologging.SetBackend(backendFormatter)

  return log 

}