package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"kafkapractisel0/models"
	"kafkapractisel0/services"
	"log"

	"github.com/segmentio/kafka-go"
)

func StartConsume(env services.Env, c services.OrderService) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        env.GetKafkaBrokerAddreses(),
		Topic:          env.GetKafkaTopic(),
		GroupID:        env.GetKafkaGroup(),
		CommitInterval: 0,
	})
	defer reader.Close()

	// dlq
	writerDLQ := kafka.NewWriter(kafka.WriterConfig{
		Brokers: env.GetKafkaBrokerAddreses(),
		Topic:   "dlq-topic",
	})
	defer writerDLQ.Close()

	ctx := context.Background()
	for {
		bytes, err := reader.ReadMessage(context.Background())
		if err != nil {
			consumerPrinf(fmt.Sprintf("error while read message, err = %v\n", err))
			continue
		}
		err = procceedMessage(bytes, c)
		if err != nil {
			consumerPrinf(fmt.Sprintf("error while procceed message, start send to dlq, len = %v\n", len(bytes.Value)))
			if err := sendDLQ(ctx, err, writerDLQ, bytes); err != nil {
				consumerPrinf(fmt.Sprintf("DLQ error while send message, err = %v\n", err))
			}
			consumerPrinf(fmt.Sprintf("sended to dlq, len = %v\n", len(bytes.Value)))

			continue
		}
		consumerPrinf(fmt.Sprintf("message commited, len = %v\n", len(bytes.Value)))

		err = reader.CommitMessages(context.Background(), bytes)
		if err != nil {
			consumerPrinf(fmt.Sprintf("error while commit messages: err = %v", err))
			if err := sendDLQ(ctx, err, writerDLQ, bytes); err != nil {
				consumerPrinf(fmt.Sprintf("DLQ error while send message, err = %v\n", err))
			}
		}
	}
}

func procceedMessage(bytes kafka.Message, c services.OrderService) error {
	var order models.OrderMessage
	err := json.Unmarshal(bytes.Value, &order)
	if err != nil {
		return fmt.Errorf("error while unmarshal, err = %w", err)
	}
	if err = c.CreateOrder(order); err != nil {
		return fmt.Errorf("error while create order, err = %w", err)
	}

	return nil
}

func sendDLQ(ctx context.Context, err error, writerDLQ *kafka.Writer, bytes kafka.Message) error {
	err = writerDLQ.WriteMessages(ctx, kafka.Message{
		Value: bytes.Value,
		Headers: []kafka.Header{
			{Key: "error-type", Value: []byte(err.Error())},
		},
	})
	return err

}

func consumerPrinf(s string) {
	log.Printf("[Kafka Consumer] %s", s)
}
