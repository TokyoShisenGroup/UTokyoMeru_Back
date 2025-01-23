package websocket

import (
    "log"
    "github.com/gorilla/websocket"
    "time"
)

type Client struct {
    hub  *Hub
    conn *websocket.Conn
    send chan []byte // 往这个通道发送的数据会写回客户端
    UserMail string    // 业务上可以加一个UserID等信息
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
            log.Println("readPump error:", err)
            break
        }
        // 这里可以把消息生产到Kafka
        kafkaProducer(message)
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