package service

import (
	"QQ/models"
	"fmt"
	gin2 "github.com/gin-gonic/gin"
	"strconv"
)

// GetUsers
// @Tags getUsers
// @Success 200 {string} json
// @Router /user/getUsers [get]
func GetUsers(context *gin2.Context) {

	context.JSON(200, gin2.H{"message": models.GetUserList()})
}

// CreateUser
// @Tags CreateUser
// @Param name query string false "用户名"
// @Param pwd query string false "密码"
// @Success 200 {string} json{"code", "message"}
// @Router /user/createUser [get]
func CreateUser(ct *gin2.Context) {
	name := ct.Query("name")
	pwd := ct.Query("pwd")
	user := models.UserBasic{
		Name: name,
		PWD:  pwd,
	}
	user.Salt =
	db := models.CreateUser(&user)
	if db.Error != nil {
		ct.JSON(-1, gin2.H{
			"message": "新增用户失败",
		})
		return
	}
	ct.JSON(200, gin2.H{
		"message": "新增用户成功",
	})

}

// DeleteUser
// @Tags DeleteUser
// @Param id query string false "用户ID"
// @Success 200 {string} json{"code", "message"}
// @Router /user/deleteUser [get]
func DeleteUser(ct *gin2.Context) {

	user := models.UserBasic{}

	id, _ := strconv.Atoi(ct.Query("id"))
	fmt.Print(id)
	user.ID = uint(id)
	db := models.DeleteUser(&user)
	if db.Error != nil {
		ct.JSON(-1, gin2.H{
			"message": "删除用户失败",
		})
		return
	}
	ct.JSON(200, gin2.H{
		"message": "删除用户成功",
	})
}

// UpdateUser
// @Tags UpdateUser
// @Param id formData string false "用户ID"
// @Param name formData string false "用户名"
// @Success 200 {string} json{"code", "message"}
// @Router /user/updateUser [post]
func UpdateUser(ct *gin2.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(ct.PostForm("id"))
	fmt.Print(id)
	user.ID = uint(id)
	db := models.UpdateUser(&user)
	if db.Error != nil {
		ct.JSON(-1, gin2.H{
			"message": "更新用户失败",
		})
		return
	}
	ct.JSON(200, gin2.H{
		"message": "更新用户成功",
	})

}
