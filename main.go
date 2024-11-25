package main

import (
  "net/http"
  "encoding/json"
  "fmt"
  "os"
)

type configData struct {
  WebUrl string `json:"webhook_url"`
  SseUrl string `json:"sse_endpoint_url"`
}

type pepitoData struct {
  Event string `json:"event"`
  Type string `json:"type"`
  Time int `json:"time"`
  Img string `json:"img"`
}

func main() {
  var config configData
  config_bytes, _ := os.ReadFile("config.json")
  err := json.Unmarshal(config_bytes, &config)
  if err != nil {
    panic(err)
  }

  req, err := http.NewRequest("GET", config.SseUrl, nil)
  if err != nil {
    panic(err)
  }
  req.Header.Set("Accept", "text/event-stream")
  req.Header.Set("Connection", "keep-alive")

  client := &http.Client{}
  res, err := client.Do(req)
  if err != nil {
    panic(err)
  }

  for {
    raw := make([]byte, 1024)

    l, err := res.Body.Read(raw)
    if err != nil {
      fmt.Fprintln(os.Stderr, "error reading data")
      panic(err)
    }
    fmt.Printf("Received message (l): %d: \n%s", l, string(raw))

    var json_string string
    _, err = fmt.Sscanf(string(raw), "data: %s", &json_string)
    if err != nil {
      panic(err)
    }

    var dat pepitoData
    err = json.Unmarshal([]byte(json_string), &dat)
    if err != nil {
      fmt.Fprintln(os.Stderr, "error unmarshalling")
      panic(err)
    }
  }
}

