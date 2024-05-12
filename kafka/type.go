package kafkalib

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Kafkalib struct {
	ka                        *kafka.AdminClient
	kp                        *kafka.Producer
	kc                        *kafka.Consumer
	BootstrapServers          string
	SchemaRegistryServers     string
	CALocation                string
	ClientPublicCertLocation  string
	ClientPrivateCertLocation string
	KeyPassword               string
	EnableCertValidation      bool
	SecurityProtocol          string
	MaxTimeout                string
}

type User struct {
	Name           string `json:"name"`
	FavoriteNumber int64  `json:"favorite_number"`
	FavoriteColor  string `json:"favorite_color"`
}

const AUTH_ADMIN = 1
const AUTH_PRODUCER = 2
const AUTH_CONSUMER = 3
