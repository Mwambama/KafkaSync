package main

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {
	//brokerAddress := "localhost:9092"
	//brokerAddress := "host.docker.internal:9092"
	brokerAddress := "localhost:9094"

	topic := "kafkasync-files"

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{brokerAddress},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	// ✅ Test connection with a dummy message
	// err := writer.WriteMessages(nil, kafka.Message{
	// 	Key:   []byte("test"),
	// 	Value: []byte("Kafka test connection message"),
	// })
	ctx := context.TODO()
	err := writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte("test"),
		Value: []byte("Kafka test connection message"),
	})

	if err != nil {
		log.Fatalf("❌ Kafka connection test failed: %v\n", err)
	} else {
		fmt.Println("✅ Kafka test message sent successfully!")
	}

	fmt.Println("✅ Kafka writer ready. Type a file name and press ENTER to send.")

	// 🔁 Loop for interactive input
	for {
		fmt.Print("📤 Enter file name: ")
		var fileName string
		_, err := fmt.Scanln(&fileName)
		if err != nil {
			log.Printf("❌ Error reading input: %v\n", err)
			continue
		}

		// Send to Kafka
		// err = writer.WriteMessages(
		// 	nil,
		// 	kafka.Message{
		// 		Key:   []byte(fileName),
		// 		Value: []byte(fileName),
		// 	},
		// )

		err = writer.WriteMessages(ctx, kafka.Message{
			Key:   []byte(fileName),
			Value: []byte(fileName),
		})

		if err != nil {
			log.Printf("❌ Failed to send message: %v\n", err)
		} else {
			fmt.Printf("📨 Sent to Kafka: %s\n", fileName)
		}
	}
}
