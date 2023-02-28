package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// 小写，将字符串转为md5
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tempStr := h.Sum(nil)
	return hex.EncodeToString(tempStr)
}

// 大写，将字符串转为md5
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

func MakePWS(plainPWD string, salt string) string {
	return Md5Encode(plainPWD + salt)
}

func ValidatePwd(plainPwd string, salt string, pwd string) bool {
	return MakePWS(plainPwd, salt) == pwd
}
