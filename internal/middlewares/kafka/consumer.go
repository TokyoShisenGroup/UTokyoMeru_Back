package kafka

import (
	"backend/internal/router/websocket"
	"context"
	"log"
	"time"


	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer struct {
	c     *kafka.Consumer
	hub   *websocket.Hub
	topic string
}

func NewConsumer(broker, group, topic string, hub *websocket.Hub) (*Consumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          group,
		"socket.keepalive.enable": true,
		"request.timeout.ms":      30000,
		"socket.timeout.ms":       30000,
		"auto.offset.reset": "earliest",
		"client.id":         "go-kafka-consumer", // 添加客户端ID便于调试
	}

	c, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}
	return &Consumer{c: c, hub: hub, topic: topic}, nil
}

func (cn *Consumer) Start(ctx context.Context) {
	// 订阅Topic
	err := cn.c.SubscribeTopics([]string{cn.topic}, nil)
	if err != nil {
		log.Println("SubscribeTopics error:", err)
		return
	}

	go func() {
		defer cn.c.Close()
		for {
			select {
			case <-ctx.Done():
				log.Println("Consumer stopped")
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
