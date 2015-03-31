package main

import (
	"flag"
	"github.com/magneticio/vamp-router/api"
	"github.com/magneticio/vamp-router/haproxy"
	"github.com/magneticio/vamp-router/helpers"
	"github.com/magneticio/vamp-router/logging"
	"github.com/magneticio/vamp-router/metrics"
	"github.com/magneticio/vamp-router/tools"
	"github.com/magneticio/vamp-router/zookeeper"
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
	kafkaHost        string
	kafkaPort        int
	zooConString     string
	zooConKey        string
	pidFilePath      string
	headless         bool
	log              *gologger.Logger
	version          = "0.7.0-RC2"
	stream           metrics.Streamer
	workDir          helpers.WorkDir
	customWorkDir    string
)

func init() {
	flag.IntVar(&port, "port", 10001, "Port/IP to use for the REST interface. Overrides $PORT0 env variable")
	flag.StringVar(&logPath, "logPath", "/logs/vamp-router.log", "Location of the log file")
	flag.StringVar(&configFilePath, "configFile", "/haproxy_new.cfg", "Target location of the generated HAproxy config file")
	flag.StringVar(&templateFilePath, "template", "configuration/templates/haproxy_config.template", "Template file to build HAproxy config")
	flag.StringVar(&jsonFilePath, "json", "/vamp_router.json", "JSON file to store internal config.")
	flag.StringVar(&binaryPath, "binary", helpers.HaproxyLocation(), "Path to the HAproxy binary")
	flag.StringVar(&kafkaHost, "kafkaHost", "", "The hostname or ip address of the Kafka host")
	flag.IntVar(&kafkaPort, "kafkaPort", 9092, "The port of the Kafka host")
	flag.StringVar(&zooConString, "zooConString", "", "A zookeeper ensemble connection string")
	flag.StringVar(&zooConKey, "zooConKey", "magneticio/vamplb", "Zookeeper root key")
	flag.StringVar(&pidFilePath, "pidFile", "/haproxy-private.pid", "Location of the HAproxy PID file")
	flag.StringVar(&customWorkDir, "customWorkDir", "", "Custom working directory for logs, configs and sockets")
	flag.BoolVar(&headless, "headless", false, "Run without any logging output to the console")
}

func main() {

	flag.Parse()

	// resolve flags and environment variables
	tools.SetValueFromEnv(&port, "VAMP_RT_PORT")
	tools.SetValueFromEnv(&logPath, "VAMP_RT_LOG_PATH")
	tools.SetValueFromEnv(&configFilePath, "VAMP_RT_CONFIG_PATH")
	tools.SetValueFromEnv(&templateFilePath, "VAMP_RT_TEMPLATE_PATH")
	tools.SetValueFromEnv(&jsonFilePath, "VAMP_RT_JSON_PATH")
	tools.SetValueFromEnv(&binaryPath, "VAMP_RT_BINARY_PATH")
	tools.SetValueFromEnv(&kafkaHost, "VAMP_RT_KAFKA_HOST")
	tools.SetValueFromEnv(&kafkaPort, "VAMP_RT_KAFKA_PORT")
	tools.SetValueFromEnv(&zooConString, "VAMP_RT_ZOO_STRING")
	tools.SetValueFromEnv(&zooConKey, "VAMP_RT_ZOO_KEY")
	tools.SetValueFromEnv(&pidFilePath, "VAMP_RT_PID_PATH")
	tools.SetValueFromEnv(&customWorkDir, "VAMP_RT_CUSTOM_WORKDIR")
	tools.SetValueFromEnv(&headless, "VAMP_RT_HEADLESS")

	//create working dir if not already there
	if err := workDir.Create(customWorkDir); err != nil {
		panic(err)
	}

	// setup logging
	log = logging.ConfigureLog(workDir.Dir()+logPath, headless)
	log.Info(logging.PrintLogo(version))

	/*

		HAproxy runtime and configuration setup

	*/

	// setup Haproxy runtime
	haRuntime := haproxy.Runtime{Binary: binaryPath}

	// setup Configuration
	haConfig := haproxy.Config{TemplateFile: templateFilePath,
		ConfigFile: workDir.Dir() + configFilePath,
		JsonFile:   workDir.Dir() + jsonFilePath,
		PidFile:    workDir.Dir() + pidFilePath,
		WorkingDir: workDir.Dir(),
	}

	log.Notice("Attempting to load config at %s", workDir.Dir()+configFilePath)
	// load config from disk
	err := haConfig.GetConfigFromDisk(haConfig.JsonFile)

	if err != nil {
		log.Notice("Did not find a config...initialzing empty config")
		haConfig.InitializeConfig()
	}

	// Render initial config
	err = haConfig.Render()
	if err != nil {
		log.Fatal("Could not render initial config, exiting...")
		os.Exit(1)
	}

	// set the Pid file
	if err := haRuntime.SetPid(haConfig.PidFile); err != nil {
		log.Notice("Pidfile exists at %s, proceeding...", workDir.Dir()+pidFilePath)
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
	stream.Init(&haRuntime, 3000, log)

	// Setup Kafka if required
	if len(kafkaHost) > 0 {

		kafkaChannel := make(chan metrics.Metric)
		stream.AddClient(kafkaChannel)

		kafka := metrics.KafkaProducer{Log: log}
		kafka.In(kafkaChannel)
		kafka.Start(kafkaHost, kafkaPort)

	}

	sseChannel := make(chan metrics.Metric)
	stream.AddClient(sseChannel)

	// Always setup SSE Stream
	sseBroker := &metrics.SSEBroker{
		make(map[chan metrics.Metric]bool),
		make(chan (chan metrics.Metric)),
		make(chan (chan metrics.Metric)),
		sseChannel,
		log,
	}

	sseBroker.In(sseChannel)
	sseBroker.Start()

	go stream.Start()

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
	restApi := api.CreateApi(log, &haConfig, &haRuntime, sseBroker)
	restApi.Run("0.0.0.0:" + strconv.Itoa(port))

}
