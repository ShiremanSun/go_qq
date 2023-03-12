package router

import (
	"QQ/docs"
	"QQ/service"
	jwt "github.com/appleboy/gin-jwt/v2"
	gin2 "github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

func Root() *gin2.Engine {
	docs.SwaggerInfo.BasePath = ""
	r := gin2.Default()
	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET(
		"/index",
		service.GetIndex,
	)

	r.POST("/login", service.AuthMiddleware.LoginHandler)
	r.NoRoute(service.AuthMiddleware.MiddlewareFunc(), func(c *gin2.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin2.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	auth := r.Group("/auth")
	auth.Use(service.AuthMiddleware.MiddlewareFunc())

	r.POST("/user/createUser", service.CreateUser)

	auth.GET("/user/getUsers", service.GetUsers)
	auth.GET("/user/deleteUser", service.DeleteUser)
	auth.POST("/user/updateUser", service.UpdateUser)

	return r
}
