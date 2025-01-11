package cdc

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avrov2"
	"github.com/mdaneshjoo/db-synchronizer/config"
	"github.com/mdaneshjoo/db-synchronizer/logger"
)

var _log = logger.NewLogger()
var Logger = _log

func newConsumer(cfg *config.Config) *kafka.Consumer {
	kcm := &kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.Brokers,
		"group.id":          cfg.Kafka.GroupID,
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(kcm)

	if err != nil {
		_log.Fatalf("failed to create Kafka consumer: %s", err)
		os.Exit(1)
	}
	_log.Infoln("Consumer Created.")
	return consumer
}

func newDeserializer(cfg *config.Config) *avrov2.Deserializer {
	client, err := schemaregistry.NewClient(schemaregistry.NewConfig(cfg.Kafka.SchemaregistryUrl))

	if err != nil {
		_log.Fatalf("Failed to create schema registry client: %s\n", err)
		os.Exit(1)
	}

	deser, err := avrov2.NewDeserializer(client, serde.ValueSerde, avrov2.NewDeserializerConfig())
	if err != nil {
		_log.Fatalf("Failed to create deserializer: %s\n", err)
		os.Exit(1)
	}
	return deser
}

func getTopics(c *kafka.Consumer, to int) []string {
	metadata, err := c.GetMetadata(nil, true, to)

	if err != nil {
		_log.Fatalf("Failed to get Kafka metadata: %v", err)
		os.Exit(1)
	}
	var topics []string

	for topic := range metadata.Topics {
		if strings.HasPrefix(topic, "postgres") {
			topics = append(topics, topic)
		}

	}
	if len(topics) == 0 {
		_log.Fatalf("No topics found in Kafka broker.")
		os.Exit(1)
	}
	_log.Infof("Litening to Topics:%v \n", topics)
	return topics
}

func Capture(to int, ch chan<- DebeziumPayload) {
	cfg := config.NewConfig()

	consumer := newConsumer(cfg)

	defer consumer.Close()

	deser := newDeserializer(cfg)

	topics := getTopics(consumer, to)

	err := consumer.SubscribeTopics(topics, nil)

	if err != nil {
		_log.Fatalf("Failed to subscribe to topics: %v", err)
		os.Exit(1)
	}

	for {
		ev := consumer.Poll(to)
		if ev == nil {
			continue
		}
		switch e := ev.(type) {
		case *kafka.Message:
			debeziumPayload := DebeziumPayload{}

			if err := debeziumPayload.loadFromAvro(deser, e); err != nil {
				continue
			}

			_log.Debugf("Message Read from: %s\n", e.TopicPartition)

			jsonPayload, err := json.Marshal(debeziumPayload)
			if err != nil {
				_log.Fatalf("Error converting to JSON:", err)
				continue
			}
			_log.Debugf("Message Is: %+v\n", string(jsonPayload))

			if e.Headers != nil {
				_log.Debugf("Headers: %v\n", e.Headers)
			}
			ch <- debeziumPayload
		case kafka.Error:
			_log.Fatalf("%v: %v\n", e.Code(), e)
		default:
			_log.Infof("Ignored %v\n", e)
		}
	}
}
