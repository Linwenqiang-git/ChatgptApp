package entity

import (
	"gorm.io/gorm"
)

type ChatRecord struct {
	gorm.Model
	UserName string
	Role     string
	Message  string
}

func (ChatRecord) TableName() string {
	return "chat_record"
}
