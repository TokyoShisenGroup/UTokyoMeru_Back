package db

import (
	"backend/config"
	"backend/internal/model"
	"context"
	"fmt"
	"log"
	"time"
	"os"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

)

var (
	DB *gorm.DB
	dsn        string
)

var rdb *redis.Client

var c = context.Background()

// 初始化redis
func RedisInit() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%d",
			config.C.RedisConfig.Host,
			config.C.RedisConfig.Port,
		),
		Password: config.C.RedisConfig.Password,  // Redis 密码
		DB:       config.C.RedisConfig.Db,        // Redis 数据库编号
		PoolSize: config.C.RedisConfig.Pool_size, // 连接池大小
	})

	_, err = rdb.Ping(c).Result() // 通过 Ping 来测试连接
	fmt.Printf("connect redis success!\n")
	return
}

func GetDatabaseInstance() (*gorm.DB, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return DB, nil
}

func init() {
	gormLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
        logger.Config{
            // 慢 SQL 阈值，超过这个时间的查询会被记录
            SlowThreshold: time.Second,
            // 日志级别
            LogLevel: logger.Warn,    // 记录所有 SQL
            // LogLevel: logger.Warn, // 只记录慢查询和错误
            // 是否忽略 ErrRecordNotFound 错误
            IgnoreRecordNotFoundError: true,
            // 是否启用彩色打印
            Colorful: true,
        },
    )
	dsn = config.C.ToDSN()
	
	// 3. 创建数据库连接
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        // 使用配置好的日志
        Logger: gormLogger,
    })
	if err != nil {
		fmt.Errorf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(&model.User{}, &model.Good{}, &model.Comment{}, &model.Message{}, &model.MessageQueue{})
	if err != nil {
		fmt.Errorf("failed to migrate database: %v", err)
	}

	// 4. 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Errorf("failed to get database instance: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 5. 保存全局数据库实例
	DB = db

	if err := RedisInit(); err != nil {
		fmt.Errorf("failed to initialize redis: %v", err)
	}

	fmt.Println("Database initialized")
}
