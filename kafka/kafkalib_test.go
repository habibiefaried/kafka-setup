package kafkalib

import (
	"testing"
	//"github.com/stretchr/testify/assert"
)

func TestCreateTopic(t *testing.T) {
	kl := &Kafkalib{
		BootstrapServers:          "broker:9093",
		CALocation:                "../kafka_server_cert.pem",
		ClientPublicCertLocation:  "../kafka_client_cert.pem",
		ClientPrivateCertLocation: "../kafka_client_key.pem",
		KeyPassword:               "datahub",
		EnableCertValidation:      false,
		SecurityProtocol:          "SSL",
		MaxTimeout:                "60s",
	}
	err := kl.AuthSSL()
	if err != nil {
		t.Error(err)
	}
	err = kl.CreateTopic("asw", 1, 1)
	if err != nil {
		t.Error(err)
	}

	defer kl.CloseConn()
}
