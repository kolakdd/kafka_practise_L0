package models

import (
	"time"
)

type PaymentCurrency string
type BankType string
type LocaleType string

const (
	CurrencyRUB PaymentCurrency = "RUB"
	CurrencyUSD PaymentCurrency = "USD"
	CurrencyEUR PaymentCurrency = "EUR"

	BankAlpha BankType = "alpha"
	BankSber  BankType = "sber"
	BankTbank BankType = "tbank"

	LocaleENG LocaleType = "eng"
	LocaleRU  LocaleType = "ru"
	LocaleKZ  LocaleType = "kz"
)

type Delivery struct {
	UID       int       `db:"uid" json:"uid"`
	Name      string    `db:"name" json:"name"`
	Phone     string    `db:"phone" json:"phone"`
	Zip       string    `db:"zip" json:"zip"`
	City      string    `db:"city" json:"city"`
	Address   string    `db:"address" json:"address"`
	Region    string    `db:"region" json:"region"`
	Email     string    `db:"email" json:"email"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Payment struct {
	Transaction  string          `db:"transaction" json:"transaction"`
	RequestID    string          `db:"request_id" json:"request_id"`
	Currency     PaymentCurrency `db:"currency" json:"currency"`
	Provider     string          `db:"provider" json:"provider"`
	Amount       int64           `db:"amount" json:"amount"`
	PaymentDt    time.Time       `db:"payment_dt" json:"payment_dt"`
	Bank         BankType        `db:"bank" json:"bank"`
	DeliveryCost int64           `db:"delivery_cost" json:"delivery_cost"`
	GoodsTotal   int64           `db:"goods_total" json:"goods_total"`
	CustomFee    int64           `db:"custom_fee" json:"custom_fee"`
	CreatedAt    time.Time       `db:"created_at" json:"created_at"`
}

type Item struct {
	ChrtID      int             `db:"chrt_id" json:"chrt_id"`
	TrackNumber string          `db:"track_number" json:"track_number"`
	Price       int64           `db:"price" json:"price"`
	Rid         string          `db:"rid" json:"rid"`
	Name        string          `db:"name" json:"name"`
	Sale        int             `db:"sale" json:"sale"`
	Size        int16           `db:"size" json:"size"`
	Currency    PaymentCurrency `db:"currency" json:"currency"`
	TotalPrice  int64           `db:"total_price" json:"total_price"`
	NmID        int64           `db:"nm_id" json:"nm_id"`
	Brand       string          `db:"brand" json:"brand"`
	Status      int16           `db:"status" json:"status"`
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`
}

type Customer struct {
	UID       int       `db:"uid" json:"uid"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Order struct {
	UID               int        `db:"uid" json:"uid"`
	TrackNumber       string     `db:"track_number" json:"track_number"`
	Entry             string     `db:"entry" json:"entry"`
	DeliveryID        int        `db:"delivery" json:"-"`
	PaymentID         string     `db:"payment" json:"-"`
	Locale            LocaleType `db:"locale" json:"locale"`
	InternalSignature string     `db:"internal_signature" json:"internal_signature"`
	CustomerID        int        `db:"customer_id" json:"customer_id"`
	DeliveryService   string     `db:"delivery_service" json:"delivery_service"`
	Shardkey          int16      `db:"shardkey" json:"shardkey"`
	SmID              int16      `db:"sm_id" json:"sm_id"`
	OofShard          int16      `db:"oof_shard" json:"oof_shard"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`

	Customer Customer `json:"customer"`
	Delivery Delivery `json:"delivery"`
	Payment  Payment  `json:"payment"`

	Items []Item `json:"items"`
}

type OrderMessage struct {
	TrackNumber       string     `db:"track_number" json:"track_number"`
	Entry             string     `db:"entry" json:"entry"`
	Locale            LocaleType `db:"locale" json:"locale"`
	InternalSignature string     `db:"internal_signature" json:"internal_signature"`
	DeliveryService   string     `db:"delivery_service" json:"delivery_service"`
	Shardkey          int16      `db:"shardkey" json:"shardkey"`
	SmID              int16      `db:"sm_id" json:"sm_id"`
	OofShard          int16      `db:"oof_shard" json:"oof_shard"`

	CustomerID  int `db:"customer_id" json:"customer_id"`
	DeliveryUID int `db:"delivery" json:"delivery"`

	Payment Payment `json:"payment"`

	ItemsID []int
}
