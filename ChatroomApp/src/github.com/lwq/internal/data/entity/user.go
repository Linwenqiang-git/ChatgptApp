package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name string `gorm:"size:100;not null"`
}

func (User) TableName() string {
	return "user"
}
