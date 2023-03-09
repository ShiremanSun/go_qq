package utils

import jwt "github.com/appleboy/gin-jwt/v2"

func GenerateJwtToken(secret string, issuer string, audience string, expired int64, userId int, userName string) (string, error) {
	sampSecret := []byte(secret)
	token := jwt.New(jwt.)
}
