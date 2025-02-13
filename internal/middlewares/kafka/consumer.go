package kafka

import (
	"backend/internal/router/websocket"
	"context"
	"fmt"
	"time"

	"backend/internal/utils/logger"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type Consumer struct {
	c     *kafka.Consumer
	hub   *websocket.Hub
	topic string
}

func NewConsumer(broker, group, topic string, hub *websocket.Hub) (*Consumer, error) {
	if broker == "" || group == "" || topic == "" {
		return nil, fmt.Errorf("broker, group and topic cannot be empty")
	}

	config := &kafka.ConfigMap{
		"bootstrap.servers":       broker,
		"group.id":                group,
		"socket.keepalive.enable": true,
		"socket.timeout.ms":       30000,
		"auto.offset.reset":       "earliest",
		"client.id":               "go-kafka-consumer",
		"enable.auto.commit":      true,
		"auto.commit.interval.ms": 5000,
		"session.timeout.ms":      10000,
		"max.poll.interval.ms":    300000,
	}

	logger.Logger.Info("Creating new Kafka consumer",
		zap.String("broker", broker),
		zap.String("group", group),
		zap.String("topic", topic))

	c, err := kafka.NewConsumer(config)
	if err != nil {
		logger.Logger.Error("Failed to create consumer",
			zap.Error(err),
			zap.String("broker", broker))
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	// 验证连接
	metadata, err := c.GetMetadata(&topic, false, 5000)
	if err != nil {
		logger.Logger.Error("Failed to get metadata",
			zap.Error(err),
			zap.String("topic", topic))
		c.Close()
		return nil, fmt.Errorf("failed to get metadata: %w", err)
	}

	// 检查 topic 是否存在
	if _, exists := metadata.Topics[topic]; !exists {
		logger.Logger.Error("Topic does not exist",
			zap.String("topic", topic))
		c.Close()
		return nil, fmt.Errorf("topic %s does not exist", topic)
	}

	consumer := &Consumer{
		c:     c,
		hub:   hub,
		topic: topic,
	}

	logger.Logger.Info("Successfully created Kafka consumer",
		zap.String("topic", topic))

	return consumer, nil
}

func (cn *Consumer) Start(ctx context.Context) {
	// 订阅Topic
	err := cn.c.SubscribeTopics([]string{cn.topic}, nil)
	if err != nil {
		logger.Logger.Error("SubscribeTopics error:", zap.Error(err))
		return
	}

	go func() {
		defer cn.c.Close()
		for {
			select {
			case <-ctx.Done():
				logger.Logger.Info("Consumer stopped")
				return
			default:
				// 拉取消息
				msg, err := cn.c.ReadMessage(500 * time.Millisecond)
				if err == nil && msg != nil {
					// 这里拿到Kafka消息，然后把消息推给 Hub
					// 在Hub里可以做广播、或者根据用户/房间ID做定向推送
					cn.hub.Broadcast <- msg.Value
				}
			}
		}
	}()
}

func (cn *Consumer) Consume(ctx context.Context) <-chan kafka.Message {
	// 创建消息通道
	messageChan := make(chan kafka.Message, 100) // 设置缓冲区大小为100
	if cn == nil {
		fmt.Println("Consumer is nil")
		close(messageChan)
		return messageChan
	}
	// 订阅Topic
	err := cn.c.SubscribeTopics([]string{cn.topic}, nil)
	if err != nil {
		logger.Logger.Error("SubscribeTopics error:", zap.Error(err))
		close(messageChan)
		return messageChan
	}

	// 启动消费协程
	go func() {
		defer close(messageChan)
		defer cn.c.Close()

		for {
			select {
			case <-ctx.Done():
				logger.Logger.Info("Consumer stopped by context")
				return
			default:
				// 拉取消息，设置超时时间为500ms
				msg, err := cn.c.ReadMessage(500 * time.Millisecond)
				if err != nil {
					if err.(kafka.Error).Code() != kafka.ErrTimedOut {
						logger.Logger.Error("Failed to read message:", zap.Error(err))
					}
					continue
				}

				if msg == nil {
					continue
				}

				// 发送消息到通道
				select {
				case messageChan <- *msg:
					// 消息发送成功
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return messageChan
}
