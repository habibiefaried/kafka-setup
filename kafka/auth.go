package kafkalib

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func (kf *Kafkalib) AuthSSL(tipe int) error {
	c := &kafka.ConfigMap{
		"bootstrap.servers":                   kf.BootstrapServers,
		"security.protocol":                   kf.SecurityProtocol,
		"ssl.ca.location":                     kf.CALocation,                // CA certificate file for SSL
		"ssl.certificate.location":            kf.ClientPublicCertLocation,  // Client certificate
		"ssl.key.location":                    kf.ClientPrivateCertLocation, // Client private key
		"ssl.key.password":                    kf.KeyPassword,               // Private key password
		"enable.ssl.certificate.verification": kf.EnableCertValidation,
	}

	if tipe == AUTH_ADMIN {
		a, err := kafka.NewAdminClient(c)
		if err != nil {
			return err
		}
		kf.ka = a
	} else if tipe == AUTH_PRODUCER {
		p, err := kafka.NewProducer(c)
		if err != nil {
			return err
		}
		kf.kp = p
	} else if tipe == AUTH_CONSUMER {
		c, err := kafka.NewConsumer(c)
		if err != nil {
			return err
		}
		kf.kc = c
	} else {
		return fmt.Errorf("Type unknown, please check your parameter")
	}

	return nil
}

func (kf *Kafkalib) CloseConn() {
	kf.ka.Close()
}
