package service

import (
	"encoding/json"

	"gorm.io/gorm"

	"backend/internal/model"
	"backend/internal/router/websocket"
	"backend/internal/db"
	"go.uber.org/zap"
)

type SyncService struct {
	db  *gorm.DB
	hub *websocket.Hub
	logger *zap.Logger
}

func NewSyncService(db *gorm.DB, hub *websocket.Hub, logger *zap.Logger) *SyncService {
	return &SyncService{
		db:  db,
		hub: hub,
		logger: logger,
	}
}

func (s *SyncService) SyncOfflineMessages(userID uint) error {
	var messages []model.Message

	var err error
	s.db, err = db.GetDatabaseInstance()
	if err != nil {
		s.logger.Error("获取数据库实例失败", zap.Error(err))
		return err
	}
	// 获取所有未读的离线消息
	if err := s.db.Joins("JOIN message_queues ON messages.id = message_queues.message_id").
		Where("message_queues.user_id = ? AND message_queues.status = Received", userID).
		Find(&messages).Error; err != nil {
		return err
	}

	// 发送离线消息
	for _, msg := range messages {
		messageBytes, _ := json.Marshal(msg)
		if err := s.hub.SendToUser(userID, messageBytes); err == nil {
			s.db.Model(&model.MessageQueue{}).
				Where("message_id = ? AND user_id = ?", msg.ID, userID).
				Update("status", 1)
		}
	}

	return nil
}
