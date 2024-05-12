package kafkalib

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Kafkalib struct {
	ka                        *kafka.AdminClient
	kp                        *kafka.Producer
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

const AUTH_ADMIN = 1
const AUTH_PRODUCER = 2
const AUTH_CONSUMER = 3
