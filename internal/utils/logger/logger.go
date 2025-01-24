package logger

import (
	"backend/config"
	"io"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger

func init() {
	// 确保日志目录存在
	if err := os.MkdirAll("logs", 0766); err != nil {
		panic(err)
	}

	// 业务日志配置
	businessHook := &lumberjack.Logger{
		Filename:   path.Join("logs", config.C.LogConfig.Backend.Filename),
		MaxSize:    config.C.LogConfig.Backend.MaxSize,
		MaxBackups: config.C.LogConfig.Backend.MaxBackups,
		MaxAge:     config.C.LogConfig.Backend.MaxAge,
		Compress:   config.C.LogConfig.Backend.Compress,
	}

	// Gin访问日志配置
	accessHook := &lumberjack.Logger{
		Filename:   path.Join("logs", config.C.LogConfig.Access.Filename),
		MaxSize:    config.C.LogConfig.Access.MaxSize,
		MaxBackups: config.C.LogConfig.Access.MaxBackups,
		MaxAge:     config.C.LogConfig.Access.MaxAge,
		Compress:   config.C.LogConfig.Access.Compress,
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	switch config.C.LogConfig.Level {
	case "debug":
		atomicLevel.SetLevel(zap.DebugLevel)
	case "info":
		atomicLevel.SetLevel(zap.InfoLevel)
	case "warn":
		atomicLevel.SetLevel(zap.WarnLevel)
	case "error":
		atomicLevel.SetLevel(zap.ErrorLevel)
	case "panic":
		atomicLevel.SetLevel(zap.PanicLevel)
	case "fatal":
		atomicLevel.SetLevel(zap.FatalLevel)
	default:
		atomicLevel.SetLevel(zap.InfoLevel)
	}

	// 创建业务日志核心
	businessCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(businessHook)),
		zap.DebugLevel,
	)

	// 创建业务日志记录器
	Logger = zap.New(businessCore, zap.AddCaller())

	// 配置Gin的日志输出
	gin.DisableConsoleColor()
	// 同时写入文件和控制台
	gin.DefaultWriter = io.MultiWriter(accessHook, os.Stdout)
}
