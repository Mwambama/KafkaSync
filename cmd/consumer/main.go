package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {

	// 🔧 Set up logging to both terminal and file
	// 🔧 Create or open the log file
	logFile, err := os.OpenFile("consumer.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("❌ Failed to open log file: %v\n", err)
		return
	}
	defer logFile.Close()

	// 🔧 Set log output to both file and terminal
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.SetFlags(log.LstdFlags | log.Lmsgprefix)

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
	//printing to console
	// fmt.Println("✅ Kafka consumer is now listening for messages...")
	// logging to file
	log.Println("✅ Kafka consumer is now listening for messages...")

	// Step 2: Consume messages in a loop
	for {
		ctx := context.TODO()
		message, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("❌ Error reading message: %v", err)
			continue
		}

		fileName := string(message.Value)
		log.Printf("⬇️  Downloading file: %s...\n", fileName)

		// Simulate download delay
		time.Sleep(1 * time.Second)

		// Simulate saving the file
		outputPath := fmt.Sprintf("./downloads/%s", fileName)
		//err = saveDummyFile(outputPath, fileName)
		// ✅ Save the file to the downloads directory
		// retry logic
		maxRetries := 3
		for attempt := 1; attempt <= maxRetries; attempt++ {
			err = saveDummyFile(outputPath, fileName)
			if err == nil {
				log.Printf("✅ Download complete: %s\n", outputPath)
				break
			}
			log.Printf("⚠️  Attempt %d: Failed to save %s: %v", attempt, fileName, err)
			time.Sleep(1 * time.Second)
		}

		if err != nil {
			log.Printf("❌ All retries failed for %s\n", fileName)
		}

		if err != nil {
			log.Printf("❌ Failed to save %s: %v", fileName, err)
		} else {
			log.Printf("✅ Download complete: %s\n", outputPath)
		}
	}
}

// saveDummyFile simulates saving a file by writing a dummy content to the specified path.
// In a real application, this would be replaced with actual file download logic.

func saveDummyFile(path string, content string) error {
	err := ensureDir("./downloads")
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte("Downloaded: "+content), 0644)
	//return errors.New("fake fail")
}

func ensureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}
