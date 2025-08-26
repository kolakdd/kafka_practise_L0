package services

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env interface {
	GetDatabaseDSN() string
	GetKafkaBrokerAddreses() []string
	GetKafkaTopic() string
	GetKafkaGroup() string
}

type env struct {
	pgUser     string `env:"APP_DB_USER"`
	pgPassword string `env:"APP_DB_PASSWORD"`
	dbName     string `env:"POSTGRES_DB"`
	pgPort     int    `env:"PG_PORT"`
	pgHost     string `env:"PG_HOST"`

	kafkaHost  string `env:"KAFKA_HOST"`
	kafkaPort  string `env:"KAFKA_PORT"`
	kafkaTopic string `env:"KAFKA_TOPIC"`
	kafkaGroup string `env:"KAFKA_GROUP"`
}

func NewEnv() Env {
	if err := godotenv.Load(); err != nil {
		log.Panic("No .env file found: ", err)
	}

	pgUser := os.Getenv("POSTGRES_USER")
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	pgHost := os.Getenv("PG_HOST")
	dbName := os.Getenv("POSTGRES_DB")
	pgPort := parseInt("PG_PORT")

	kafkaHost := os.Getenv("KAFKA_HOST")
	kafkaPort := os.Getenv("KAFKA_PORT")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	kafkaGroup := os.Getenv("KAFKA_GROUP")

	return &env{pgUser, pgPassword, dbName, pgPort, pgHost, kafkaHost, kafkaPort, kafkaTopic, kafkaGroup}
}

func parseInt(evnKey string) int {
	v := os.Getenv(evnKey)
	vInt, err := strconv.Atoi(v)
	if err != nil {
		log.Panicf("err while parse env key=%s", evnKey)
	}
	return vInt

}

func (e *env) GetKafkaTopic() string {
	return e.kafkaTopic
}

func (e *env) GetKafkaGroup() string {
	return e.kafkaGroup
}

func (e *env) GetDatabaseDSN() string {
	user := e.pgUser
	password := e.pgPassword
	host := e.pgHost
	db := e.dbName
	port := e.pgPort
	return fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=%d sslmode=disable", user, password, host, db, port)
}

func (e *env) GetKafkaBrokerAddreses() []string {
	return []string{fmt.Sprintf("%s:%s", e.kafkaHost, e.kafkaPort)}
}
