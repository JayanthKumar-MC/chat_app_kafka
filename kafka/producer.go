package kafka

import (
	"chat-app/config"
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

var writer *kafka.Writer

/*func InitProducer() {
	kafkaURL := config.Cfg.Kafka.URL
	log.Printf("Connecting to Kafka on %s\n", kafkaURL)
	writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaURL},
		Topic:   config.Cfg.Kafka.Topic,
	})
}*/

func InitProducer() error {
	kafkaURL := config.Cfg.Kafka.URL
	log.Printf("Connecting to Kafka on %s\n", kafkaURL)
	writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaURL},
		Topic:   config.Cfg.Kafka.Topic,
	})
	// Example check: Attempt to connect or validate writer is not nil
	// This is a placeholder, adapt based on actual kafka-go API capabilities
	if writer == nil {
		log.Printf("Failed to initialize Kafka writer")
		return fmt.Errorf("failed to initialize Kafka writer")
	}
	return nil
}

func ProduceMessage(msg []byte) error {
	if writer == nil {
		log.Printf("Kafka writer is not initialized")
		return fmt.Errorf("Kafka writer is not initialized")
	}
	if err := writer.WriteMessages(context.TODO(), kafka.Message{Value: msg}); err != nil {
		log.Printf("Could not write message to Kafka: %v\n", err)
		return err
	}
	return nil
}

func CloseProducer() {
	writer.Close()
}
