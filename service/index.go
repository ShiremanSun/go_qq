package service

import (
	gin2 "github.com/gin-gonic/gin"
)

func GetIndex(context *gin2.Context) {
	context.JSON(200, gin2.H{"message": "pong"})
}
