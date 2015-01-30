package metrics

import (
  "fmt"
)

//  a very simple consumer of a metric stream. It just outputs to the console
type SimpleConsumer struct {
}

func (s *SimpleConsumer) Consume(c chan []byte) {
for {
  value := <- c
  fmt.Printf(string(value[:]))
  }
}