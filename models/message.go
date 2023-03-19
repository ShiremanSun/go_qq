package models

import (
	"github.com/fatih/set"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	FromId      string
	TargetId    string
	MessageType string
	MediaType   int
	Content     string
	Pic         string
	Url         string
	Desc        string
	Amount      string
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Con       *websocket.Conn
	DataQueue chan []byte
	GroupSet  set.Interface // 加入的群组
}
