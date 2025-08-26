package emulator

import (
	"context"
	"encoding/json"
	"fmt"
	"kafkapractisel0/mock"
	"kafkapractisel0/models"
	"kafkapractisel0/services"
	"log"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
)

func StartEmulate(env services.Env, e services.EmulatorService) {
	ctx := context.Background()
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: env.GetKafkaBrokerAddreses(),
		Topic:   env.GetKafkaTopic(),
	})
	defer writer.Close()

	for {
		time.Sleep(time.Second)

		order := newEmulatorMessage(e)

		bytes, _ := json.Marshal(order)
		err := writer.WriteMessages(ctx, kafka.Message{
			Value: bytes,
		})
		if err != nil {
			emulatorPrinf(fmt.Sprintf("error while write message, err = %v", err))
		}
		emulatorPrinf(fmt.Sprintf("message sended , len = %v", len(bytes)))
	}

}

func newEmulatorMessage(e services.EmulatorService) models.OrderMessage {
	deliveryUID, err := e.GetRandomDelivery()
	if err != nil {
		deliveryUID = 0
	}

	items := []int{}
	// Ошибка эмулятора - иногда слайс items - пустой !!!
	for i := 1; i < rand.Intn(9)+1; i++ {
		items = append(items, i)
	}

	message := models.OrderMessage{
		TrackNumber:       "TRACKNUMBER",
		Entry:             "BIL",
		Locale:            models.LocaleRU,
		InternalSignature: "sigma",
		CustomerID:        rand.Intn(10) + 1,
		DeliveryService:   "meest",
		Shardkey:          1,
		SmID:              42,
		OofShard:          52,
		Payment: models.Payment{
			Transaction:  mock.RidGenerator(),
			RequestID:    "",
			Currency:     models.CurrencyRUB,
			Provider:     "XXXPAY",
			Amount:       1337,
			PaymentDt:    time.Now(),
			Bank:         models.BankAlpha,
			DeliveryCost: 1337 - 500,
			GoodsTotal:   500,
			CustomFee:    0,
		},
		DeliveryUID: deliveryUID,
		ItemsID:     items,
	}
	return message

}

func emulatorPrinf(s string) {
	log.Printf("[Kafka Emulator] %s", s)
}
