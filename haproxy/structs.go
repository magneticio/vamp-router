package haproxy

import (
	"sync"
)

// placeholder
type Runtime struct {
	Binary string
}

// Main configuration object for load balancers. This contains all variables and is passed to
// the templating engine.
type Config struct {
	Frontends    []*Frontend   `json:"frontends" binding:"required"`
	Backends     []*Backend    `json:"backends" binding:"required"`
	PidFile      string        `json:"-"`
	Mutex        *sync.RWMutex `json:"-"`
	TemplateFile string        `json:"-"`
	ConfigFile   string        `json:"-"`
	JsonFile     string        `json:"-"`
}

// Defines a single haproxy "backend".
type Backend struct {
	Name           string           `json:"name" binding:"required"`
	Mode           string           `json:"mode" binding:"required"`
	BackendServers []*BackendServer `json:"servers" binding:"required"`
	Options        ProxyOptions     `json:"options"`
	ProxyMode      bool             `json:"proxyMode" binding:"required"`
}

// Defines a single haproxy "frontend".
type Frontend struct {
	Name           string       `json:"name" binding:"required"`
	Mode           string       `json:"mode" binding:"required"`
	BindPort       int          `json:"bindPort"`
	BindIp         string       `json:"bindIp"`
	UnixSock       string       `json:"unixSock"`
	SockProtocol   string       `json:"sockProtocol"`
	Options        ProxyOptions `json:"options"`
	DefaultBackend string       `json:"defaultBackend" binding:"required"`
	ACLs           []*ACL       `json:"acls"`
	HttpSpikeLimit SpikeLimit   `json:"httpSpikeLimit,omitempty"`
	TcpSpikeLimit  SpikeLimit   `json:"tcpSpikeLimit,omitempty"`
}

type ProxyOptions struct {
	AbortOnClose    bool `json:"abortOnClose"`
	AllBackups      bool `json:"allBackups"`
	CheckCache      bool `json:"checkCache"`
	ForwardFor      bool `json:"forwardFor"`
	HttpClose       bool `json:"httpClose"`
	HttpCheck       bool `json:"httpCheck"`
	SslHelloCheck   bool `json:"sslHelloCheck"`
	TcpKeepAlive    bool `json:"tcpKeepAlive"`
	TcpLog          bool `json:"tcpLog"`
	TcpSmartAccept  bool `json:"tcpSmartAccept"`
	TcpSmartConnect bool `json:"tcpSmartConnect"`
}

// Defines a server which exists in a backend.
type BackendServer struct {
	Name          string `json:"name" binding:"required"`
	Host          string `json:"host"`
	Port          int    `json:"port"`
	UnixSock      string `json:"unixSock"`
	Weight        int    `json:"weight" binding:"required"`
	MaxConn       int    `json:"maxconn"`
	Check         bool   `json:"check"`
	CheckInterval int    `json:"checkInterval"`
}

// Defines an ACL
type ACL struct {
	Name    string `json:"name" binding:"required"`
	Backend string `json:"backend" binding:"required"`
	Pattern string `json:"pattern" binding:"required"`
}

// Defines a rate limiting setup

type SpikeLimit struct {
	SampleTime string `json:"sampleTime,omitempty" binding:"required"`
	ExpiryTime string `json:"expiryTime,omitempty" binding:"required"`
	Rate       int    `json:"rate,omitempty" binding:"required"`
}

// Struct to hold the output from the /stats endpoint
type StatsGroup struct {
	Pxname         string `json:"pxname"`
	Svname         string `json:"svname"`
	Qcur           string `json:"qcur"`
	Qmax           string `json:"qmax"`
	Scur           string `json:"scur"`
	Smax           string `json:"smax"`
	Slim           string `json:"slim"`
	Stot           string `json:"stot"`
	Bin            string `json:"bin"`
	Bout           string `json:"bout"`
	Dreq           string `json:"dreq"`
	Dresp          string `json:"dresp"`
	Ereq           string `json:"ereq"`
	Econ           string `json:"econ"`
	Eresp          string `json:"eresp"`
	Wretr          string `json:"wretr"`
	Wredis         string `json:"wredis"`
	Status         string `json:"status"`
	Weight         string `json:"weight"`
	Act            string `json:"act"`
	Bck            string `json:"bck"`
	Chkfail        string `json:"chkfail"`
	Chkdown        string `json:"chkdown"`
	Lastchg        string `json:"lastchg"`
	Downtime       string `json:"downtime"`
	Qlimit         string `json:"qlimit"`
	Pid            string `json:"pid"`
	Iid            string `json:"iid"`
	Sid            string `json:"sid"`
	Throttle       string `json:"throttle"`
	Lbtot          string `json:"lbtot"`
	Tracked        string `json:"tracked"`
	_Type          string `json:"type"`
	Rate           string `json:"rate"`
	Rate_lim       string `json:"rate_lim"`
	Rate_max       string `json:"rate_max"`
	Check_status   string `json:"check_status"`
	Check_code     string `json:"check_code"`
	Check_duration string `json:"check_duration"`
	Hrsp_1xx       string `json:"hrsp_1xx"`
	Hrsp_2xx       string `json:"hrsp_2xx"`
	Hrsp_3xx       string `json:"hrsp_3xx"`
	Hrsp_4xx       string `json:"hrsp_4xx"`
	Hrsp_5xx       string `json:"hrsp_5xx"`
	Hrsp_other     string `json:"hrsp_other"`
	Hanafail       string `json:"hanafail"`
	Req_rate       string `json:"req_rate"`
	Req_rate_max   string `json:"req_rate_max"`
	Req_tot        string `json:"req_tot"`
	Cli_abrt       string `json:"cli_abrt"`
	Srv_abrt       string `json:"srv_abrt"`
	Comp_in        string `json:"comp_in"`
	Comp_out       string `json:"comp_out"`
	Comp_byp       string `json:"comp_byp"`
	Comp_rsp       string `json:"comp_rsp"`
	Lastsess       string `json:"lastsess"`
	Last_chk       string `json:"last_chk"`
	Last_agt       string `json:"last_agt"`
	Qtime          string `json:"qtime"`
	Ctime          string `json:"ctime"`
	Rtime          string `json:"rtime"`
	Ttime          string `json:"ttime"`
}

// struct to hold the output from the /info endpoint
type Info struct {
	Name                        string `json:"Name"`
	Version                     string `json:"Version"`
	Release_date                string `json:"Release_date"`
	Nbproc                      string `json:"Nbproc"`
	Process_num                 string `json:"Process_num"`
	Pid                         string `json:"Pid"`
	Uptime                      string `json:"Uptime"`
	Uptime_sec                  string `json:"Uptime_sec"`
	Memmax_MB                   string `json:"Memmax_MB"`
	Ulimitn                     string `json:"Ulimit-n"`
	Maxsock                     string `json:"Maxsock"`
	Maxconn                     string `json:"Maxconn"`
	Hard_maxconn                string `json:"Hard_maxconn"`
	CurrConns                   string `json:"CurrConns"`
	CumConns                    string `json:"CumConns"`
	CumReq                      string `json:"CumReq"`
	MaxSslConns                 string `json:"MaxSslConns"`
	CurrSslConns                string `json:"CurrSslConns"`
	CumSslConns                 string `json:"CumSslConns"`
	Maxpipes                    string `json:"Maxpipes"`
	PipesUsed                   string `json:"PipesUsed"`
	PipesFree                   string `json:"PipesFree"`
	ConnRate                    string `json:"ConnRate"`
	ConnRateLimit               string `json:"ConnRateLimit"`
	MaxConnRate                 string `json:"MaxConnRate"`
	SessRate                    string `json:"SessRate"`
	SessRateLimit               string `json:"SessRateLimit"`
	MaxSessRate                 string `json:"MaxSessRate"`
	SslRate                     string `json:"SslRate"`
	SslRateLimit                string `json:"SslRateLimit"`
	MaxSslRate                  string `json:"MaxSslRate"`
	SslFrontendKeyRate          string `json:"SslFrontendKeyRate"`
	SslFrontendMaxKeyRate       string `json:"SslFrontendMaxKeyRate"`
	SslFrontendSessionReuse_pct string `json:"SslFrontendSessionReuse_pct"`
	SslBackendKeyRate           string `json:"SslBackendKeyRate"`
	SslBackendMaxKeyRate        string `json:"SslBackendMaxKeyRate"`
	SslCacheLookups             string `json:"SslCacheLookups"`
	SslCacheMisses              string `json:"SslCacheMisses"`
	CompressBpsIn               string `json:"CompressBpsIn"`
	CompressBpsOut              string `json:"CompressBpsOut"`
	CompressBpsRateLim          string `json:"CompressBpsRateLim"`
	ZlibMemUsage                string `json:"ZlibMemUsage"`
	MaxZlibMemUsage             string `json:"MaxZlibMemUsage"`
	Tasks                       string `json:"Tasks"`
	Run_queue                   string `json:"Run_queue"`
	Idle_pct                    string `json:"Idle_pct"`
	node                        string `json:"node"`
	description                 string `json:"description"`
}
