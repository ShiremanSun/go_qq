package main

import (
	"QQ/models"
	"QQ/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMysql()
	utils.DB.AutoMigrate(&models.UserBasic{})
}
