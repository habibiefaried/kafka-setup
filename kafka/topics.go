package kafkalib

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"time"
)

func (kf *Kafkalib) CreateTopic(topic string, partition int, replicationFactor int) error {
	// Contexts are used to abort or limit the amount of time
	// the Admin call blocks waiting for a result.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create topics on cluster.
	// Set Admin options to wait for the operation to finish (or at most 60s)
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
