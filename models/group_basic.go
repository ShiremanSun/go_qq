package models

import "gorm.io/gorm"

type GroupBasic struct {
	gorm.Model
	Name    string
	OwnerId string
	Icon    string
	Type    int
	Desc    string
}

func (receiver *GroupBasic) TableName() string {
	return "group_basic"
}
