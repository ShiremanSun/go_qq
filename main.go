package main

import (
	"QQ/router"
	"QQ/utils"
)

func main() {

	utils.InitConfig()
	utils.InitMysql()
	r := router.Root()
	r.Run(":8081")

}
