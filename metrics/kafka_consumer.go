package metrics

import (

  "github.com/Shopify/sarama"
  "log"
  "strconv"
  "time"
)


type KafkaMetricsConsumer struct {
  Name string
  Metrics chan []byte
}

func(k *KafkaMetricsConsumer) Consume(m chan []byte) {
  k.Metrics = m
}

func (k *KafkaMetricsConsumer) Init(host string, port int) {

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

  go k.pushMetrics(producer)

}

func (k *KafkaMetricsConsumer) pushMetrics(producer *sarama.Producer) {
  for {
        metric := <- k.Metrics
        log.Println(string(metric   ))
        err := producer.SendMessage("loadbalancer.all", sarama.StringEncoder("lbmetrics"), sarama.StringEncoder(metric))
        if err != nil {
          log.Println("error sending to Kafka ")
        }
      } 
}