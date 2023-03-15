package main

import (
	"QQ/utils"
)

func main() {
	utils.InitConfig()
	//utils.InitMysql()
	//utils.DB.AutoMigrate(&models.UserBasic{})
	utils.InitRedis()
}
