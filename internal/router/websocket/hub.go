package websocket

import (
	"backend/internal/utils/logger"
	"sync"
	"errors"

	"go.uber.org/zap"
)

type Hub struct {
	// 注册/注销通道
	Register   chan *Client
	Unregister chan *Client

	// 全部活跃的客户端连接
	clients map[*Client]bool

	// 用于内部广播/分发的通道（可选）
	Broadcast chan []byte

	// 添加用户ID到客户端的映射
	userClients map[uint][]*Client // 一个用户可能有多个连接
	mu          sync.RWMutex
	onlineUsers sync.Map // 用于存储在线用户状态
}

func NewHub() *Hub {
	return &Hub{
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Broadcast:   make(chan []byte),
		clients:     make(map[*Client]bool),
		userClients: make(map[uint][]*Client),
		onlineUsers: sync.Map{},
		mu:          sync.RWMutex{},
	}
}

// 添加发送私聊消息的方法
func (h *Hub) SendToUser(targetUserID uint, message []byte) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, exists := h.userClients[targetUserID]; exists {
		for _, client := range clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	} else {
		return errors.New("user not found")
	}
	return nil
}

func (h *Hub) IsUserOnline(userID uint) bool {
	_, online := h.onlineUsers.Load(userID)
	return online
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.clients[client] = true
			// 添加到用户映射
			h.userClients[client.UserID] = append(h.userClients[client.UserID], client)
			h.onlineUsers.Store(client.UserID, true)
			h.mu.Unlock()
			logger.Logger.Info("New client connected",
				zap.Uint("userID", client.UserID),
				zap.Any("client", client))

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				// 从用户映射中移除
				if clients, exists := h.userClients[client.UserID]; exists {
					newClients := make([]*Client, 0)
					for _, c := range clients {
						if c != client {
							newClients = append(newClients, c)
						}
					}
					if len(newClients) == 0 {
						delete(h.userClients, client.UserID)
						h.onlineUsers.Delete(client.UserID)
					} else {
						h.userClients[client.UserID] = newClients
					}
				}
				close(client.send)
			}
			h.mu.Unlock()
			logger.Logger.Info("Client disconnected", zap.Any("client", client))
		case message := <-h.Broadcast:
			// 这里是简单的广播给所有客户端
			// 如果要按用户、房间号分发，需要加上用户ID/房间ID逻辑
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}