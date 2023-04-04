package main

import (
	provider "github.com/lwq/configs/wire"
	entity "github.com/lwq/internal/data/entity"
	"gorm.io/gorm"
)

func MigrateDb() *gorm.DB {

	dbContext, err := provider.GetDbContext()
	if err != nil {
		panic("connect error：" + err.Error())
	}
	return dbContext.GetDb()
}

func main() {
	db := MigrateDb()
	/*
		初始化 User 和 ChatRecord 表
	*/
	db.AutoMigrate(&entity.User{})
	db.AutoMigrate(&entity.ChatRecord{})
}
