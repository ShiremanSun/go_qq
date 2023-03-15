package main

import (
	"QQ/router"
	"QQ/service"
	"QQ/utils"
)

func main() {

	utils.InitConfig()
	utils.InitMysql()
	service.InitAuth()
	utils.InitRedis()
	r := router.Root()
	r.Run(":8081")

}
