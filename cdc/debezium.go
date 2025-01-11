package cdc

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avrov2"
)

type DebeziumPayload struct {
	Op     string      `avro:"op"`
	Before interface{} `avro:"before"`
	After  interface{} `avro:"after"`
	Source struct {
		Table string `avro:"table"`
	} `avro:"source"`
}

func (d *DebeziumPayload) loadFromAvro(deserializer *avrov2.Deserializer, e *kafka.Message) error {
	err := deserializer.DeserializeInto(*e.TopicPartition.Topic, e.Value, d)
	if err != nil {
		_log.Fatalf("Failed to deserialize payload: %s\n", err)
		return err
	}
	return nil
}
