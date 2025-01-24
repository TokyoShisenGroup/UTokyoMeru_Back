// internal/kafka/producer.go
package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"backend/internal/utils/logger"
	"go.uber.org/zap"
)

type Producer struct {
	p     *kafka.Producer
	topic string
}

func NewProducer(broker, topic string) (*Producer, error) {
	logger.Logger.Info("Attempting to connect to Kafka at: %s", zap.String("broker", broker))
	config := &kafka.ConfigMap{
		"bootstrap.servers": broker,
		"socket.keepalive.enable": true,
		"request.timeout.ms":      30000,
		"socket.timeout.ms":       30000,
		"auto.offset.reset": "earliest",
		"client.id":         "go-kafka-producer", // 添加客户端ID便于调试
	}

	p, err := kafka.NewProducer(config)
	if err != nil {
		logger.Logger.Error("Failed to create producer: %v", zap.Error(err))
		return nil, err
	}
	logger.Logger.Info("Successfully connected to Kafka")
	return &Producer{p: p, topic: topic}, nil
}

// ProduceMessage 将消息发送到Kafka
func (pr *Producer) ProduceMessage(msg []byte) {
	err := pr.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &pr.topic,
			Partition: kafka.PartitionAny,
		},
		Value: msg,
	}, nil)
	if err != nil {
		logger.Logger.Error("Produce error:", zap.Error(err))
	}
}

// 关闭 producer
func (pr *Producer) Close() {
	pr.p.Close()
}
