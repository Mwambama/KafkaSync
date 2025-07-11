package main

import (
	"context" // Import context for managing request lifecycle
	"fmt"
	"log"
	"time" // Import time for sleep functionality

	"github.com/segmentio/kafka-go"
)

func main() {
	// Define Kafka config
	brokerAddress := "localhost:9094"
	topic := "kafkasync-files"
	groupID := "file-consumer-group"

	// Step 1: Create a Kafka Reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{brokerAddress},
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       1,           // Smallest message size to fetch
		MaxBytes:       10e6,        // 10MB max size
		CommitInterval: time.Second, // Commit offsets every second
	})
	defer reader.Close()

	fmt.Println("✅ Kafka consumer is now listening for messages...")

	// Step 2: Consume messages in a loop
	for {
		ctx := context.TODO()
		message, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("❌ Error reading message: %v", err)
			continue
		}

		fmt.Printf("📥 Received file name: %s\n", string(message.Value))
	}
}
