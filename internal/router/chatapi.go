package router

import (
	"backend/internal/router/websocket"

	"github.com/gin-gonic/gin"
)

func ChatHandler(c *gin.Context) {
	// Kafka 生产回调
	kafkaProducerFunc := func(msg []byte) {
		Producer.ProduceMessage(msg)
	}
	websocket.ServeWs(Hub, c, kafkaProducerFunc)
}
