package kafkalib

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Kafkalib struct {
	ka                        *kafka.AdminClient
	BootstrapServers          string
	CALocation                string
	ClientPublicCertLocation  string
	ClientPrivateCertLocation string
	KeyPassword               string
	EnableCertValidation      bool
	SecurityProtocol          string
	MaxTimeout                string
}
