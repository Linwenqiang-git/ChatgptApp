package data

import (
	. "github.com/lwq/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DbContext struct {
	db *gorm.DB
}

func (c DbContext) GetDb() *gorm.DB {
	return c.db
}

func ProvideDbContext(configure Configure) (DbContext, error) {
	dsn := configure.MysqlSetting.GetConnectInfo()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return DbContext{db: db}, err
}
