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
	err := kl.AuthSSL(AUTH_ADMIN)
	if err != nil {
		t.Error(err)
	}
	err = kl.CreateTopic("test", 1, 1)
	if err != nil {
		t.Error(err)
	}

	defer kl.CloseConn()
}

func TestPublishConsumeMessage(t *testing.T) {
	kl := &Kafkalib{
		BootstrapServers:          "broker:9093",
		SchemaRegistryServers:     "http://broker:8081",
		CALocation:                "../kafka_server_cert.pem",
		ClientPublicCertLocation:  "../kafka_client_cert.pem",
		ClientPrivateCertLocation: "../kafka_client_key.pem",
		KeyPassword:               "datahub",
		EnableCertValidation:      false,
		SecurityProtocol:          "SSL",
		MaxTimeout:                "60s",
		ConsumerGroupName:         "group1",
	}

	err := kl.AuthSSL(AUTH_PRODUCER)
	if err != nil {
		t.Error(err)
	}

	err = kl.PublishMessage("test", &User{
		Name:           "First user",
		FavoriteNumber: 42,
		FavoriteColor:  "blue",
	})
	if err != nil {
		t.Error(nil)
	}

	err = kl.AuthSSL(AUTH_CONSUMER)
	if err != nil {
		t.Error(err)
	}

	m, err := kl.ConsumeMessage("test")
	if err != nil {
		t.Error(nil)
	}

	t.Logf("%+v", m)
}
