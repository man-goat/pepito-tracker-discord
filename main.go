package main

import (
  "github.com/r3labs/sse"
  "encoding/json"
  "fmt"
  "os"
)

type configData struct {
  WebUrl string `json:"webhook_url"`
  SseUrl string `json:"sse_endpoint_url"`
  SseStream string `json:"sse_stream"`
}

type streamData struct {
  Event string `json:"event"`
  Type string `json:"type"`
  Time int `json:"time"`
  Img string `json:"img"`
}

func main() {
  if len(os.Args) < 2 {
    panic("no config")
  }

  var config configData
  config_bytes, _ := os.ReadFile(os.Args[1])
  err := json.Unmarshal(config_bytes, &config)
  check(err, "error unmarshaling config file")

  // done := make(chan bool)

  events := make(chan *sse.Event)

  sse_client := sse.NewClient(config.SseUrl)
  sse_client.SubscribeChan(config.SseStream, events)

  for event := range events {
    var data streamData
    err = json.Unmarshal(event.Data, &data)
    check(err, "unable to unmarshal stream data")
    fmt.Printf("Received data: event: %s, type: %s, time: %d, img: %s\n", data.Event, data.Type, data.Time, data.Img)
  }
}

func check(err error, msg string) {
  if err != nil {
    fmt.Fprintf(os.Stderr, msg)
    panic(err)
  }
}
