package main

import (
	_ "backend/internal/db"
	"backend/internal/router"
	"backend/internal/utils/logger"
	"fmt"

	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	defer logger.Logger.Sync()

	// 设置 Gin 模式
	gin.SetMode(gin.DebugMode)

	// 启动服务器
	err := router.Router.Run(":8100")
	if err != nil {
		logger.Logger.Error("服务器启动失败", zap.Error(err))
	}

	quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    <-quit
    fmt.Println("Shutting down...")
	router.Cancel()
}
