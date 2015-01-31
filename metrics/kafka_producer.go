package metrics

import (

  "github.com/Shopify/sarama"
  "encoding/json"
  "log"
  "strconv"
  "time"
)


type KafkaProducer struct {
  metricsChannel chan Metric
}

func(k *KafkaProducer) In(c chan Metric) {
  k.metricsChannel = c
}

func (k *KafkaProducer) Start(host string, port int) {

  connection := host + ":" + strconv.Itoa(port)

  log.Println("Connecting to Kafka on " + connection + "...")

  clientConfig := sarama.NewClientConfig()
  clientConfig.WaitForElection = (10 * time.Second)


  client, err := sarama.NewClient("client_id", []string{connection}, clientConfig)
  if err != nil {
    panic(err)
  } else {
    log.Println("Connection to Kafka successful")
  }

  /**
  *  Create a producer with some specific setting
  */
  producerConfig := sarama.NewProducerConfig()

  // if delivering messages async,  buffer them for at most MaxBufferTime
  producerConfig.MaxBufferTime = (2 * time.Second)

  // max bytes in buffer
  producerConfig.MaxBufferedBytes = 51200

  // Use zip compression
  producerConfig.Compression = 0

  // We are just streaming metrics, so don't not wait for any Kafka Acks.
  producerConfig.RequiredAcks = -1

  producer, err := sarama.NewProducer(client, producerConfig)
  if err != nil {
    panic(err)
  }

  go k.produce(producer)

}

func (k *KafkaProducer) produce(producer *sarama.Producer) {
  for {
        metric := <- k.metricsChannel
        json, err := json.MarshalIndent(metric, "", " ")
        if err != nil {
          return
        }
        err = producer.SendMessage("loadbalancer.all", sarama.StringEncoder("lbmetrics"), sarama.StringEncoder(json))
        if err != nil {
          log.Println("error sending to Kafka ")
        }
      } 
}