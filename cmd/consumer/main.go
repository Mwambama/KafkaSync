package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
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

		fileName := string(message.Value)
		fmt.Printf("⬇️  Downloading file: %s...\n", fileName)

		// Simulate download delay
		time.Sleep(1 * time.Second)

		// Simulate saving the file
		outputPath := fmt.Sprintf("./downloads/%s", fileName)
		err = saveDummyFile(outputPath, fileName)
		if err != nil {
			log.Printf("❌ Failed to save %s: %v", fileName, err)
		} else {
			fmt.Printf("✅ Download complete: %s\n", outputPath)
		}
	}
}

// ✅ Move helper functions OUTSIDE main

func saveDummyFile(path string, content string) error {
	err := ensureDir("./downloads")
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte("Downloaded: "+content), 0644)
}

func ensureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// ✅ Add error handling for directory creation
// func ensureDir(dir string) error {
// 	if _, err := os.Stat(dir); os.IsNotExist(err) {
// 		return os.MkdirAll(dir, 0755)
// 	}
// 	return nil
// }
