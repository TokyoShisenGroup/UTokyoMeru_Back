package config

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
)

type Config struct {
	DbConfig    DbConfig
	RedisConfig RedisConfig
	LogConfig   LogConfig
}

type DbConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
	Timezone string `json:"timezone"`
}

type RedisConfig struct {
	Host      string `json:"host"`
	Port      int64  `json:"port"`
	Password  string `json:"password"`
	Db        int    `json:"db"`
	Pool_size int    `json:"pool_size"`
}

type LogConfig struct {
	Backend struct {
		Filename   string `json:"filename"`
		MaxSize    int    `json:"maxsize"`
		MaxBackups int    `json:"maxbackups"`
		MaxAge     int    `json:"maxage"`
		Compress   bool   `json:"compress"`
	} `json:"backend"`
	Access struct {
		Filename   string `json:"filename"`
		MaxSize    int    `json:"maxsize"`
		MaxBackups int    `json:"maxbackups"`
		MaxAge     int    `json:"maxage"`
		Compress   bool   `json:"compress"`
	} `json:"access"`
	Level string `json:"level"`
}

var C Config

func init() {
	var f = flag.String("f", "config.yaml", "config file")
	flag.Parse()
	conf.MustLoad(*f, &C)
}

func (c *Config) ToDSN() string {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		c.DbConfig.Host, c.DbConfig.User, c.DbConfig.Password, c.DbConfig.Dbname, c.DbConfig.Port, c.DbConfig.SSLMode, c.DbConfig.Timezone)
	//fmt.Println("database config loaded:", dsn)
	return dsn
}
