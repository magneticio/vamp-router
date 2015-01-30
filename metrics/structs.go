package metrics

type Metric struct {

  Name string `json:"name"`
  Value int `json:"value"`
  Timestamp int64 `json:"timestamp"`

}