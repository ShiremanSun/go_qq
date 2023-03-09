package utils

import (
	"QQ/models"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

var AuthMiddleware *jwt.GinJWTMiddleware

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

func InitAuth() {
	authMiddleware, error := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "UTC",
		Key:         []byte(viper.GetString("auth.secret")),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.UserBasic); ok {
				return jwt.MapClaims{
					identityKey: v.Name,
				}
			}
			return jwt.MapClaims{}
		},
	})

	AuthMiddleware = authMiddleware
}
