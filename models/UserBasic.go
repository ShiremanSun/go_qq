package models

import (
	"QQ/utils"
	"gorm.io/gorm"
)

/*
实体类
*/

type UserBasic struct {
	gorm.Model
	Name          string
	PWD           string
	Phone         string
	Email         string
	Identity      string
	ClientIP      string
	ClientPort    string
	LoginTime     uint64
	HeartbeatTime uint64 "json:heartbeat_time"
	LogOutTime    uint64 "json:log_out_time"
	IsLogOut      bool
	DeviceInfo    string
	Salt          string
}

/*
添加创建表方法，返回表名
函数前的括号，相当于kotlin中的扩展方法
*/

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() (data []*UserBasic) {
	data = make([]*UserBasic, 10)
	utils.DB.Find(&data)
	return data
}

func CreateUser(user *UserBasic) *gorm.DB {
	return utils.DB.Create(user)
}

func DeleteUser(user *UserBasic) *gorm.DB {
	return utils.DB.Delete(user)
}

func UpdateUser(user *UserBasic) *gorm.DB {
	return utils.DB.Model(user).Updates(UserBasic{Name: user.Name})
}

func FindUserByName(name string) *UserBasic {
	user := &UserBasic{}
	utils.DB.Where("name = ?", name).Find(user)
	return user
}

func FindUserByNameAndPWD(name string, pwd string) *UserBasic {
	user := &UserBasic{}
	utils.DB.Where("name = ? and pwd = ?", name, pwd).Find(user)
	return user
}
