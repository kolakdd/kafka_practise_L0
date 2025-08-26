package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"kafkapractisel0/models"
	"kafkapractisel0/repo/cache"
	"kafkapractisel0/services"
	"log"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	env   services.Env
	o     services.OrderService
	cache cache.CacheRepo
}

func NewConsumer(env services.Env, o services.OrderService, cache cache.CacheRepo) Consumer {
	return Consumer{env: env, o: o, cache: cache}
}

func (c *Consumer) StartConsume() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        c.env.GetKafkaBrokerAddreses(),
		Topic:          c.env.GetKafkaTopic(),
		GroupID:        c.env.GetKafkaGroup(),
		CommitInterval: 0,
	})
	defer reader.Close()

	// dlq
	writerDLQ := kafka.NewWriter(kafka.WriterConfig{
		Brokers: c.env.GetKafkaBrokerAddreses(),
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
		err = c.procceedMessage(bytes)
		if err != nil {
			consumerPrinf(fmt.Sprintf("error while procceed message, start send to dlq, len = %v\n", len(bytes.Value)))
			if err := c.sendDLQ(ctx, err, writerDLQ, bytes); err != nil {
				consumerPrinf(fmt.Sprintf("DLQ error while send message, err = %v\n", err))
			}
			consumerPrinf(fmt.Sprintf("sended to dlq, len = %v\n", len(bytes.Value)))
			continue
		}
		consumerPrinf(fmt.Sprintf("message commited, len = %v\n", len(bytes.Value)))

		err = reader.CommitMessages(context.Background(), bytes)
		if err != nil {
			consumerPrinf(fmt.Sprintf("error while commit messages: err = %v", err))
			if err := c.sendDLQ(ctx, err, writerDLQ, bytes); err != nil {
				consumerPrinf(fmt.Sprintf("DLQ error while send message, err = %v\n", err))
			}
		}
	}
}

func (c *Consumer) procceedMessage(bytes kafka.Message) error {
	var order models.OrderMessage
	err := json.Unmarshal(bytes.Value, &order)
	if err != nil {
		return fmt.Errorf("error while unmarshal, err = %w", err)
	}
	orderDB, err := c.o.CreateOrder(order)
	if err != nil {
		return fmt.Errorf("error while create order, err = %w", err)
	}
	orderDBBytes, err := json.Marshal(orderDB)
	if err != nil {
		return fmt.Errorf("error while marshal, err = %w", err)
	}
	c.cache.Set(orderDB.UID, orderDBBytes)
	return nil
}

func (c *Consumer) sendDLQ(ctx context.Context, err error, writerDLQ *kafka.Writer, bytes kafka.Message) error {
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
