package main

import (
	"flag"
	"github.com/magneticio/vamp-loadbalancer/api"
	"github.com/magneticio/vamp-loadbalancer/haproxy"
	"github.com/magneticio/vamp-loadbalancer/helpers"
	"github.com/magneticio/vamp-loadbalancer/logging"
	"github.com/magneticio/vamp-loadbalancer/metrics"
	"github.com/magneticio/vamp-loadbalancer/tools"
	"github.com/magneticio/vamp-loadbalancer/zookeeper"
	gologger "github.com/op/go-logging"
	"os"
	"strconv"
)

var (
	// Set all commandline arguments
	port             int
	logPath          string
	configFilePath   string
	templateFilePath string
	jsonFilePath     string
	binaryPath       string
	kafkaSwitch      bool
	kafkaHost        string
	kafkaPort        int
	zooConString     string
	zooConKey        string
	pidFilePath      string
	log              *gologger.Logger
	version          = "0.3.0"
	stream           metrics.Streamer
)

func init() {
	flag.IntVar(&port, "port", 10001, "Port/IP to use for the REST interface. Overrides $PORT0 env variable")
	flag.StringVar(&logPath, "logPath", "/tmp/vamp-loadbalancer.log", "Location of the log file")
	flag.StringVar(&configFilePath, "lbConfigFile", "/tmp/haproxy_new.cfg", "Location of the target HAproxy config file")
	flag.StringVar(&templateFilePath, "lbTemplate", "configuration/templates/haproxy_config.template", "Template file to build HAproxy load balancer config")
	flag.StringVar(&jsonFilePath, "lbJson", "/tmp/vamp_loadbalancer.json", "JSON file to store internal config.")
	flag.StringVar(&binaryPath, "binary", helpers.HaproxyLocation(), "Path to the HAproxy binary")
	flag.BoolVar(&kafkaSwitch, "kafkaSwitch", false, "Switch whether to enable Kafka streaming")
	flag.StringVar(&kafkaHost, "kafkaHost", "", "The hostname or ip address of the Kafka host")
	flag.IntVar(&kafkaPort, "kafkaPort", 9092, "The port of the Kafka host")
	flag.StringVar(&zooConString, "zooConString", "", "A zookeeper ensemble connection string")
	flag.StringVar(&zooConKey, "zooConKey", "magneticio/vamplb", "Zookeeper root key")
	flag.StringVar(&pidFilePath, "pidFile", "/tmp/haproxy-private.pid", "Location of the HAproxy PID file")
}

func main() {

	flag.Parse()

	// resolve flags and environment variables
	tools.SetValueFromEnv(&port, "VAMP_LB_PORT")
	tools.SetValueFromEnv(&logPath, "VAMP_LB_LOG_PATH")
	tools.SetValueFromEnv(&configFilePath, "VAMP_LB_CONFIG_PATH")
	tools.SetValueFromEnv(&templateFilePath, "VAMP_LB_TEMPLATE_PATH")
	tools.SetValueFromEnv(&jsonFilePath, "VAMP_LB_JSON_PATH")
	tools.SetValueFromEnv(&binaryPath, "VAMP_LB_BINARY_PATH")
	tools.SetValueFromEnv(&kafkaSwitch, "VAMP_LB_KAFKA_SWITCH")
	tools.SetValueFromEnv(&kafkaHost, "VAMP_LB_KAFKA_HOST")
	tools.SetValueFromEnv(&kafkaPort, "VAMP_LB_KAFKA_PORT")
	tools.SetValueFromEnv(&zooConString, "VAMP_LB_ZOO_STRING")
	tools.SetValueFromEnv(&zooConKey, "VAMP_LB_ZOO_KEY")
	tools.SetValueFromEnv(&pidFilePath, "VAMP_LB_PID_PATH")

	// setup logging
	log = logging.ConfigureLog(logPath)
	log.Info(logging.PrintLogo(version))

	/*

		HAproxy runtime and configuration setup

	*/

	// setup Haproxy runtime
	haRuntime := haproxy.Runtime{Binary: binaryPath}

	// setup Configuration
	haConfig := haproxy.Config{TemplateFile: templateFilePath, ConfigFile: configFilePath, JsonFile: jsonFilePath, PidFile: pidFilePath}

	log.Notice("Attempting to load config at %s", configFilePath)
	// load config from disk
	err := haConfig.GetConfigFromDisk(haConfig.JsonFile)
	if err != nil {
		log.Notice("Did not find a config, loading example config...")
		err = haConfig.GetConfigFromDisk("examples/example1.json")
		if err != nil {
			log.Warning("Could not load example file from disk...")
		}
	}

	// Render initial config
	err = haConfig.Render()
	if err != nil {
		log.Fatal("Could not render initial config, exiting...")
		os.Exit(1)
	}

	// set the Pid file
	done := haRuntime.SetPid(haConfig.PidFile)
	if done == false {
		log.Notice("Pidfile exists at %s, proceeding...", pidFilePath)
	} else {
		log.Notice("Created new pidfile...")
	}

	// start or reload
	err = haRuntime.Reload(&haConfig)
	if err != nil {
		log.Fatal("Error while reloading haproxy: " + err.Error())
		os.Exit(1)
	}

	/*

		Metric streaming setup

	*/

	log.Notice("Initializing metric streams...")

	// Initialize the stream from a runtime
	stream.Init(&haRuntime, 3000)
	metricsChannel := make(chan metrics.Metric)

	// push the metrics output into the metrics channel
	go stream.Out(metricsChannel)

	// Setup Kafka if required
	if len(kafkaHost) > 0 {

		kafka := metrics.KafkaProducer{}
		kafka.In(metricsChannel)
		kafka.Start(kafkaHost, kafkaPort)

	}

	// Setup SSE Stream
	sseBroker := &metrics.SSEBroker{
		make(map[chan metrics.Metric]bool),
		make(chan (chan metrics.Metric)),
		make(chan (chan metrics.Metric)),
		metricsChannel,
		log,
	}

	sseBroker.In(metricsChannel)
	sseBroker.Start()

	/*

		Zookeeper setup

	*/

	if len(zooConString) > 0 {

		log.Notice("Initializing Zookeeper connection to " + zooConString + zooConKey)
		zkClient := zookeeper.ZkClient{}
		err := zkClient.Init(zooConString, &haConfig, log)

		if err != nil {
			log.Error("Error initializing Zookeeper...")
		}
		zkClient.Watch(zooConKey)
	}

	/*

		Rest API setup

	*/
	log.Notice("Initializing REST API...")
	api.CreateApi(port, &haConfig, &haRuntime, log, sseBroker).Run("0.0.0.0:" + strconv.Itoa(port))

}
