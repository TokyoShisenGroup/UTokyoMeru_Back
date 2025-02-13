package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte // 往这个通道发送的数据会写回客户端
	UserID uint
}

type Message struct {
	Type      string `json:"type"`      // "chat" 表示私聊
	FromID    uint   `json:"from_id"`   // 发送者ID
	ToID      uint   `json:"to_id"`     // 接收者ID
	Content   string `json:"content"`   // 消息内容
	Timestamp int64  `json:"timestamp"` // 时间戳
}

// 读协程：从 WebSocket 读消息 -> 处理 或 生产到 Kafka
func (c *Client) readPump(kafkaProducer func([]byte)) {
	defer func() {
		c.hub.Unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	// 设置读超时，保持心跳
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// 解析消息
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("error parsing message: %v", err)
			continue
		}

		// 设置发送者ID和时间戳
		msg.FromID = c.UserID
		msg.Timestamp = time.Now().Unix()

		// 重新序列化消息
		messageBytes, _ := json.Marshal(msg)

		// 如果是私聊消息，发送给特定用户
		if msg.Type == "chat" {
			c.hub.SendToUser(msg.ToID, messageBytes)
		}

		// 同时发送到Kafka用于消息持久化
		kafkaProducer(messageBytes)
	}
}

// 写协程：监听 c.send 通道，将消息写回客户端
func (c *Client) writePump() {
	ticker := time.NewTicker(50 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// 通道关闭，关闭连接
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// 将消息写回客户端
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("writePump error:", err)
				return
			}

		case <-ticker.C:
			// 定时发送 ping，保持连接
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("writePump ping error:", err)
				return
			}
		}
	}
}
