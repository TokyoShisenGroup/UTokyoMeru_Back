package service

import (
	"backend/internal/db"
	"backend/internal/middlewares/kafka"
	"backend/internal/model"
	"backend/internal/router/websocket"
	"context"
	"encoding/json"
	"time"
	"go.uber.org/zap"
)

type MessageService struct {
	producer   *kafka.Producer
	consumer   *kafka.Consumer
	hub        *websocket.Hub
	logger     *zap.Logger
}

func NewMessageService(producer *kafka.Producer, consumer *kafka.Consumer, hub *websocket.Hub, logger *zap.Logger) *MessageService {
	return &MessageService{
		producer: producer,
		consumer: consumer,
		hub:      hub,
		logger: logger,
	}
}

// 处理接收到的消息
func (s *MessageService) HandleIncomingMessage(msg *websocket.Message) error {
	crudmsg := db.MessageCRUD{}
	crudmq := db.MessageQueueCRUD{}
	// 1. 保存消息到数据库
	message := &model.Message{
		From:    msg.FromID,
		To:      msg.ToID,
		Content: msg.Content,
		Status:  "Received",
	}
	if err := crudmsg.CreateByObject(message); err != nil {
		return err
	}
	// 2. 创建消息队列记录
	queueItem := &model.MessageQueue{
		UserID: msg.ToID,
		MessageID:    message.ID,
		Status:       "Received",
	}
	if err := crudmq.CreateByObject(queueItem); err != nil {
		return err
	}
	// 3. 发送到Kafka的发送队列
	messageBytes, _ := json.Marshal(message)
	s.producer.ProduceMessage(messageBytes)

	return nil
}

// 处理发送队列的消息
func (s *MessageService) ProcessSendQueue(ctx context.Context) {
	// 使用外部传入的 context，便于控制生命周期
	for msg := range s.consumer.Consume(ctx) {
		var message model.Message
		if err := json.Unmarshal(msg.Value, &message); err != nil {
			s.logger.Error("Failed to unmarshal message", zap.Error(err))
			continue
		}

		// 检查用户是否在线
		if s.hub.IsUserOnline(message.To) {
			if err := s.hub.SendToUser(message.To, msg.Value); err == nil {
				s.markMessageAsSent(message.ID)
				continue
			}
		}
		
		s.retryLater(message.ID)
	}
}

// 处理重试队列
func (s *MessageService) ProcessRetryQueue() {
	crudmq := db.MessageQueueCRUD{}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()  // 防止资源泄露
	
	for range ticker.C {
		var queueItems []model.MessageQueue
		queueItems, err := crudmq.FindAllByField("status", "Received", "created_at", "asc")
		if err != nil {
			continue
		}
		
		for _, item := range queueItems {
			message := item.Message
			
			if s.hub.IsUserOnline(item.UserID) {
				messageBytes, _ := json.Marshal(message)
				if err := s.hub.SendToUser(item.UserID, messageBytes); err == nil {
					s.markMessageAsSent(item.MessageID)
				}
			}
		}
	}
}

func (s *MessageService) markMessageAsSent(messageID uint) {
	crudmq := db.MessageQueueCRUD{}
	crudmsg := db.MessageCRUD{}

	message, err := crudmq.FindByID(messageID)
	if err != nil {
		return
	}

	message.Status = "Sent"
	message.Message.Status = "Sent"
	crudmq.UpdateByObject(message)
	crudmsg.UpdateByObject(&message.Message)
}

func (s *MessageService) retryLater(messageID uint) {
	// 更新重试次数和下次重试时间
}