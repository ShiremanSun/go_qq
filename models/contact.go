package models

import "gorm.io/gorm"

type Contact struct {
	gorm.Model
	OwnerId  int //谁的关系
	TargetId int //和谁的关系
	Type     int // 关系类型
}

//
func (receiver *Contact) TableName() string {
	return "contact"
}
