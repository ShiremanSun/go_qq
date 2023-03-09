package router

import (
	"QQ/docs"
	"QQ/service"
	gin2 "github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Root() *gin2.Engine {
	docs.SwaggerInfo.BasePath = ""
	r := gin2.Default()
	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET(
		"/index",
		service.GetIndex,
	)
	r.GET("/user/getUsers", service.GetUsers)
	r.POST("/user/createUser", service.CreateUser)
	r.GET("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)
	r.POST("/user/login", service.Login)
	r.POST("/login")

	return r
}
