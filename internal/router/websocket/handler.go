package websocket

import (
	"net/http"
	"github.com/gorilla/websocket"
	"log"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        // 生产环境需要校验Origin
        return true
    },
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, kafkaProducer func([]byte)) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Upgrade error:", err)
        return
    }
    client := &Client{
        hub:  hub,
        conn: conn,
        send: make(chan []byte, 256),
    }
    hub.Register <- client

    go client.writePump()
    go client.readPump(kafkaProducer)
}