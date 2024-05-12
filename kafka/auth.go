package kafkalib

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func (kf *Kafkalib) AuthSSL() error {
	a, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers":                   kf.BootstrapServers,
		"security.protocol":                   kf.SecurityProtocol,
		"ssl.ca.location":                     kf.CALocation,                // CA certificate file for SSL
		"ssl.certificate.location":            kf.ClientPublicCertLocation,  // Client certificate
		"ssl.key.location":                    kf.ClientPrivateCertLocation, // Client private key
		"ssl.key.password":                    kf.KeyPassword,               // Private key password
		"enable.ssl.certificate.verification": kf.EnableCertValidation,
	})
	if err != nil {
		return err
	}
	kf.ka = a

	return nil
}

func (kf *Kafkalib) CloseConn() {
	kf.ka.Close()
}
