package main

import (
	"flag"
	// "os"

	"fmt"
	"github.com/magneticio/vamp-loadbalancer/logging"
	"github.com/magneticio/vamp-loadbalancer/haproxy"
	"github.com/magneticio/vamp-loadbalancer/api"
	"github.com/magneticio/vamp-loadbalancer/metrics"
	gologger "github.com/op/go-logging"
	"strconv"
	"os"

)

var (

	// Set all commandline arguments
	port             int
	logPath          string
	configFilePath   string
	templateFilePath string
	binaryPath       string
	kafkaSwitch      bool
	kafkaHost        string
	kafkaPort        int
	zooConString     string
	zooConKey        string
	pidFilePath      string
	log	*gologger.Logger
	version = "0.1"	
	stream metrics.Streamer						
)

func init() {
	flag.IntVar(&port, "port", 10001, "Port/IP to use for the REST interface. Overrides $PORT0 env variable")
	flag.StringVar(&logPath, "logPath", "/tmp/vamp-loadbalancer.log", "Location of the log file")
	flag.StringVar(&configFilePath, "lbConfigFile", "/tmp/haproxy_new.cfg", "Location of the target HAproxy config file")
	flag.StringVar(&templateFilePath, "lbTemplate", "configuration/templates/haproxy_config.template", "Template file to build HAproxy load balancer config")
	flag.StringVar(&binaryPath, "binary", "/usr/local/sbin/haproxy", "Path to the HAproxy binary")
	flag.BoolVar(&kafkaSwitch, "kafkaSwitch", false, "Switch whether to enable Kafka streaming")
	flag.StringVar(&kafkaHost, "kafkaHost", "localhost", "The hostname or ip address of the Kafka host")
	flag.IntVar(&kafkaPort, "kafkaPort", 9092, "The port of the Kafka host")
	flag.StringVar(&zooConString, "zooConString", "localhost", "A zookeeper ensemble connection string")
	flag.StringVar(&zooConKey, "zooConKey", "magnetic", "Zookeeper root key")
	flag.StringVar(&pidFilePath, "pidFile", "/tmp/haproxy-private.pid", "Location of the HAproxy PID file")
}

func main() {

	flag.Parse()

	// resolve flags and environment variables
	setValueFromEnv(&port, "VAMP_LB_PORT")
	setValueFromEnv(&logPath, "VAMP_LB_LOG_PATH")
	setValueFromEnv(&configFilePath, "VAMP_LB_CONFIG_PATH")
	setValueFromEnv(&templateFilePath, "VAMP_LB_TEMPLATE_PATH")
	setValueFromEnv(&binaryPath, "VAMP_LB_BINARY_PATH")
	setValueFromEnv(&kafkaSwitch, "VAMP_LB_KAFKA_SWITCH")
	setValueFromEnv(&kafkaHost, "VAMP_LB_KAFKA_HOST")
	setValueFromEnv(&kafkaPort, "VAMP_LB_KAFKA_PORT")
	setValueFromEnv(&zooConString, "VAMP_LB_ZOO_STRING")
	setValueFromEnv(&zooConKey, "VAMP_LB_ZOO_KEY")
	setValueFromEnv(&pidFilePath, "VAMP_LB_PID_PATH")
	
	// setup logging
	log = logging.ConfigureLog(logPath)
	log.Info(logging.PrintLogo(version))

	// setup Haproxy runtime
	haRuntime := haproxy.Runtime{ Binary: binaryPath }

	// setup Configuration
	haConfig := haproxy.Config{ 
		TemplateFile: templateFilePath,
		ConfigFile: configFilePath,
		PidFile: pidFilePath,
	}

	log.Notice("Attempting to load config from disk..")
	// load config from disk
	err := haConfig.GetConfigFromDisk(haConfig.ConfigFile)
	if err != nil {
		log.Warning("Failed to read config from disk, loading example config...")
		err = haConfig.GetConfigFromDisk("examples/example1.json")
		if err != nil {
			log.Warning("Could not load example file from disk...")
		}
	}

	// set the Pid file
	err = haRuntime.SetPid(haConfig.PidFile)
	if err != nil {
		log.Notice("Pidfile exists, proceeding...")
	}
	err = haConfig.Render()

	log.Notice("Initializing metric streams...")

	// Initialize the stream from a runtime
	stream.Init(&haRuntime)

	metricsChannel := make(chan []byte)

	// push the JSON feed into the metrics channel
	go stream.ToJson(metricsChannel)

	// Consume the metrics channel
	// consumer := metrics.SimpleConsumer{}
	// consumer.Consume(metricsChannel)

	kafkaConsumer := metrics.KafkaMetricsConsumer{Name: "kafka", Metrics: metricsChannel}
	kafkaConsumer.Init(kafkaHost, kafkaPort)


	log.Notice("Initializing REST Api...")

	api.CreateApi(port, &haConfig, &haRuntime)



}

func setValueFromEnv(field interface{}, envVar string) {

	env := os.Getenv(envVar)
	if len(env) > 0 {

		switch v := field.(type) { 
    	default:
        fmt.Printf("unexpected type %T while parsing flags",v)
    	case *int:
        		field, _ = strconv.Atoi(env)
    	case *string:
    	  		field = env
    	case *bool:
    	  		field, _ = strconv.ParseBool(env)
    } 
	}
}	





