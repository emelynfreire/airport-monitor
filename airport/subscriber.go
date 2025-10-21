package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"

    "github.com/segmentio/kafka-go"
)

type Flight struct {
    ID      string `json:"id"`
    Airline string `json:"airline"`
    Origin  string `json:"origin"`
    Destiny string `json:"destiny"`
    Time    string `json:"time"`
}

func consume(topic, totemName string) {
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers: []string{"localhost:9092"},
        Topic:   topic,
        GroupID: totemName + "-group",
    })

    for {
        m, err := reader.ReadMessage(context.Background())
        if err != nil {
            log.Fatal(err)
        }

        var flight Flight
        _ = json.Unmarshal(m.Value, &flight)
        fmt.Printf("[%s] New flight: %+v\n", totemName, flight)
    }
}

