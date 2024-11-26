package main

import (
  "github.com/r3labs/sse/v2"
  "encoding/json"
  "io"
  "fmt"
  "os"
  "net/http"
  "bytes"
)

type configData struct {
  WebUrl string `json:"webhook_url"`
  SseUrl string `json:"sse_endpoint_url"`
  SseStream string `json:"sse_stream"`
  AvatarUrl string `json:"webhook_avatar_url"`
  Username string `json:"webhook_username"`
}

type streamData struct {
  Event string `json:"event"`
  Type string `json:"type"`
  Time int `json:"time"`
  Img string `json:"img"`
}

type webhookPayload struct {
  Content   string `json:"content"`
  Username  string `json:"username"`
  AvatarURL string `json:"avatar_url"`
}


func main() {
  if len(os.Args) < 2 {
    panic("no config")
  }

  var config configData
  config_bytes, _ := os.ReadFile(os.Args[1])
  err := json.Unmarshal(config_bytes, &config)
  check(err, "error unmarshaling config file")

  events := make(chan *sse.Event)

  sse_client := sse.NewClient(config.SseUrl)
  sse_client.SubscribeChan(config.SseStream, events)

  client := &http.Client{}

  for event := range events {
    var data streamData
    err = json.Unmarshal(event.Data, &data)
    check(err, "unable to unmarshal stream data")
    fmt.Printf("Received data: event: %s, type: %s, time: %d, img: %s\n", data.Event, data.Type, data.Time, data.Img)

    if data.Event == "heartbeat" || data.Type == "" || data.Img == "" {
      continue
    }
    message := fmt.Sprintf("Pepito is %s: %s", data.Type, data.Img)
    payload := webhookPayload{ Content: message, AvatarURL: config.AvatarUrl, Username: "Pepito-Bot"}
    bts, err := json.Marshal(payload)
    check(err, "can't marshal payload")

    req, err := http.NewRequest("POST", config.WebUrl, bytes.NewReader(bts))
    check(err, "")
    req.Header.Set("Content-Type", "application/json")

    resp, err := client.Do(req)
    check(err, "")

    bd, err := io.ReadAll(resp.Body)
    check(err, "")
    fmt.Printf("%s\n", string(bd))
  }
}
func check(err error, msg string) {
  if err != nil {
    fmt.Fprintf(os.Stderr, msg)
    panic(err)
  }
}
