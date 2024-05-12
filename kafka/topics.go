package kafkalib

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"
	"time"
)

func (kf *Kafkalib) CreateTopic(topic string, partition int, replicationFactor int) error {
	// Contexts are used to abort or limit the amount of time
	// the Admin call blocks waiting for a result.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create topics on cluster.
	// Set Admin options to wait for the operation to finish
	maxDur, err := time.ParseDuration(kf.MaxTimeout)
	if err != nil {
		return err
	}

	_, err = kf.ka.CreateTopics(
		ctx,
		[]kafka.TopicSpecification{{
			Topic:             topic,
			NumPartitions:     partition,
			ReplicationFactor: replicationFactor}},
		kafka.SetAdminOperationTimeout(maxDur))
	if err != nil {
		return err
	}

	return nil
}

func (kf *Kafkalib) PublishMessage(topic string, message interface{}) error {
	deliveryChan := make(chan kafka.Event)
	client, err := schemaregistry.NewClient(schemaregistry.NewConfig(kf.SchemaRegistryServers))
	if err != nil {
		return fmt.Errorf("Failed to create schema registry client: %s\n", err)
	}

	ser, err := jsonschema.NewSerializer(client, serde.ValueSerde, jsonschema.NewSerializerConfig())
	if err != nil {
		return fmt.Errorf("Failed to create serializer: %s\n", err)
	}

	payload, err := ser.Serialize(topic, message)
	if err != nil {
		return fmt.Errorf("Failed to serialize payload: %s\n", err)
	}

	err = kf.kp.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          payload,
		Headers:        []kafka.Header{},
	}, deliveryChan)
	if err != nil {
		return fmt.Errorf("Produce failed: %v\n", err)
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)
	return nil
}

func (kf *Kafkalib) ConsumeMessage(topic string) (*User, error) {
	client, err := schemaregistry.NewClient(schemaregistry.NewConfig(kf.SchemaRegistryServers))
	if err != nil {
		return nil, fmt.Errorf("Failed to create schema registry client: %s\n", err)
	}
	deser, err := jsonschema.NewDeserializer(client, serde.ValueSerde, jsonschema.NewDeserializerConfig())

	if err != nil {
		return nil, fmt.Errorf("Failed to create deserializer: %s\n", err)
	}

	err = kf.kc.SubscribeTopics([]string{topic}, nil)
	run := true

	for run {
		ev := kf.kc.Poll(100)
		if ev == nil {
			continue
		}

		switch e := ev.(type) {
		case *kafka.Message:
			value := User{}
			err := deser.DeserializeInto(*e.TopicPartition.Topic, e.Value, &value)
			if err != nil {
				return nil, fmt.Errorf("Failed to deserialize payload: %s\n", err)
			} else {
				// fmt.Printf("Message on %s:\n%+v\n", e.TopicPartition, value)
				return &value, nil
			}
		case kafka.Error:
			return nil, fmt.Errorf("Error: %v: %v\n", e.Code(), e)
		default:
			fmt.Printf("Ignored %v\n", e)
		}
	}

	fmt.Printf("Closing consumer\n")
	kf.kc.Close()
	return nil, nil
}
