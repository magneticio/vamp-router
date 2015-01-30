package metrics

type SSEConsumer struct {
  Name string
  Metrics chan []byte
}

func(k *SSEConsumer) Consume(m chan []byte) {
  k.Metrics = m
}