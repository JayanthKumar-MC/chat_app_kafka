package kafka

import (
	"chat-app/config"
	"chat-app/models"
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

func ConsumeMessages() {
	log.Printf("Consuming message")
	kafkaURL := config.Cfg.Kafka.URL
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaURL},
		Topic:   config.Cfg.Kafka.Topic,
		GroupID: "chat-group",
	})
	defer reader.Close()

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Could not read message: %v\n", err)
			continue
		}

		log.Printf("Received message: %s\n", string(m.Value))

		// Parse and save message to database
		var msg models.Message
		// Assuming the message is in JSON format
		err = json.Unmarshal(m.Value, &msg)
		if err != nil {
			log.Printf("Could not parse message: %v\n", err)
			continue
		}

		err = models.UpdateMessageStatus(msg.ID, "delivered")
		if err != nil {
			log.Printf("Could not update message status to delivered: %v\n", err)
		}
	}
}
