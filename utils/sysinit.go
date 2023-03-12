package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

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
