package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB
var Redis *redis.Client

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	error := viper.ReadInConfig()
	if error != nil {
		fmt.Print(error)
	}
}

func InitMysql() {
	// 添加MySql日志监控
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		})
	db, error := gorm.Open(mysql.Open(viper.GetString("mysql.address")), &gorm.Config{Logger: newLogger})
	DB = db
	if error != nil {
		panic("连接数据库失败")
	}
}

func InitRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleCoon"),
	})

	pong, err := Redis.Ping(context.Background()).Result()

	if err != nil {
		fmt.Print(err)
	} else {
		fmt.Print(pong)
	}
}
