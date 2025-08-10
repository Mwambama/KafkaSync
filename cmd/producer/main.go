package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type DownloadNotification struct {
	Hash     string `json:"info_hash"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

func main() {
	brokerAddress := "localhost:9094"
	topic := "kafkasync-files"

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{brokerAddress},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	fmt.Println("✅ Kafka writer ready. Enter file details:")

	for {
		var hash, name, location string

		fmt.Print("🔢 Enter hash: ")
		fmt.Scanln(&hash)

		fmt.Print("📄 Enter file name: ")
		fmt.Scanln(&name)

		fmt.Print("🌐 Enter remote location path: ")
		fmt.Scanln(&location)

		notification := DownloadNotification{
			Hash:     hash,
			Name:     name,
			Location: location,
		}

		payload, err := json.Marshal(notification)
		if err != nil {
			log.Printf("❌ JSON encode failed: %v", err)
			continue
		}

		err = writer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(name),
			Value: payload,
		})

		if err != nil {
			log.Printf("❌ Failed to send message: %v", err)
		} else {
			fmt.Printf("📨 Sent to Kafka: %+v\n", notification)
		}
	}
}
