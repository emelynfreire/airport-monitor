package main

import (
    "context"
    "encoding/json"
    "log"
    "time"

    "github.com/segmentio/kafka-go"
)

type Flight struct {
    ID        string `json:"id"`
    Airline   string `json:"airline"`
    Origin    string `json:"origin"`
    Destiny   string `json:"destiny"`
    Time      string `json:"time"` // ISO8601 string
}

func publish(flight Flight, topic string) {
    writer := &kafka.Writer{
        Addr:     kafka.TCP("localhost:9092"),
        Topic:    topic,
        Balancer: &kafka.LeastBytes{},
    }
    defer writer.Close()

    msg, _ := json.Marshal(flight)
    err := writer.WriteMessages(context.Background(), kafka.Message{
        Key:   []byte(flight.ID),
        Value: msg,
    })
    if err != nil {
        log.Println("Failed to publish:", err)
    } else {
        log.Println("Published to", topic)
    }
}

func main() {
    flight1 := Flight{
        ID:      "AZ1234",
        Airline: "AZUL",
        Origin:  "Fortaleza",
        Destiny: "Natal",
        Time:    time.Now().Format(time.RFC3339),
    }

    flight2 := Flight{
        ID:      "G3456",
        Airline: "GOL",
        Origin:  "Natal",
        Destiny: "Recife",
        Time:    time.Now().Add(1 * time.Hour).Format(time.RFC3339),
    }

    publish(flight1, "arrival-flights")
    publish(flight2, "departure-flights")
}

