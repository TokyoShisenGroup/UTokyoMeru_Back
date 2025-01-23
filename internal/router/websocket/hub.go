package websocket

import (
    "sync"
    "fmt"
)

type Hub struct {
    // 注册/注销通道
    Register   chan *Client
    Unregister chan *Client

    // 全部活跃的客户端连接
    clients    map[*Client]bool

    // 用于内部广播/分发的通道（可选）
    Broadcast  chan []byte

    mu         sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        Register:   make(chan *Client),
        Unregister: make(chan *Client),
        Broadcast:  make(chan []byte),
        clients:    make(map[*Client]bool),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.Register:
            h.mu.Lock()
            h.clients[client] = true
            h.mu.Unlock()
            fmt.Println("New client connected")

        case client := <-h.Unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            h.mu.Unlock()

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

