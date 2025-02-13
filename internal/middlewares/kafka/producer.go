// internal/kafka/producer.go
package kafka

import (
	"backend/internal/utils/logger"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type Producer struct {
	P     *kafka.Producer
	Topic string
}

func NewProducer(broker, topic string) (*Producer, error) {
	logger.Logger.Info("Attempting to connect to Kafka", zap.String("broker", broker))
	config := &kafka.ConfigMap{
		"bootstrap.servers":       broker,
		"socket.keepalive.enable": true,
		"request.timeout.ms":      30000,
		"socket.timeout.ms":       30000,
		"client.id":               "go-kafka-producer",
		"acks":                    "all",
		"retries":                 3,
		"retry.backoff.ms":        100,
	}

	p, err := kafka.NewProducer(config)
	if err != nil {
		logger.Logger.Error("Failed to create producer", zap.Error(err))
		return nil, err
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					logger.Logger.Error("Failed to deliver message",
						zap.Error(ev.TopicPartition.Error),
						zap.String("topic", *ev.TopicPartition.Topic))
				}
			case kafka.Error:
				logger.Logger.Error("Kafka error",
					zap.Error(ev),
					zap.String("code", ev.Code().String()))
			}
		}
	}()

	logger.Logger.Info("Successfully connected to Kafka")
	return &Producer{P: p, Topic: topic}, nil
}

// ProduceMessage 将消息发送到Kafka
func (pr *Producer) ProduceMessage(msg []byte) error {
    deliveryChan := make(chan kafka.Event)
    err := pr.P.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{
            Topic:     &pr.Topic,
            Partition: kafka.PartitionAny,
        },
        Value: msg,
    }, deliveryChan)

    if err != nil {
        logger.Logger.Error("Failed to produce message", zap.Error(err))
        return err
    }

    e := <-deliveryChan
    m := e.(*kafka.Message)

    if m.TopicPartition.Error != nil {
        logger.Logger.Error("Message delivery failed", zap.Error(m.TopicPartition.Error))
        return m.TopicPartition.Error
    }

    close(deliveryChan)
    return nil
}

// 关闭 producer
func (pr *Producer) Close() {
	pr.P.Close()
}
