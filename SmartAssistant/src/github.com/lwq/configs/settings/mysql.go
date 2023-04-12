package settings

import (
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
)

type MysqlSetting struct {
	ip       string
	port     int
	user     string
	password string
	database string
	charset  string
	show_sql bool
}

func (s *MysqlSetting) GetConnectInfo() string {
	//"user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	var dsnBuilder strings.Builder
	dsnBuilder.WriteString(s.user)
	dsnBuilder.WriteString(":")
	dsnBuilder.WriteString(s.password)
	dsnBuilder.WriteString("@tcp")
	dsnBuilder.WriteString("(")
	dsnBuilder.WriteString(s.ip)
	dsnBuilder.WriteString(":")
	dsnBuilder.WriteString(strconv.Itoa(s.port))
	dsnBuilder.WriteString(")/")
	dsnBuilder.WriteString(s.database)
	dsnBuilder.WriteString("?charset=")
	dsnBuilder.WriteString(s.charset)
	dsnBuilder.WriteString("&parseTime=True&loc=Local")
	return dsnBuilder.String()
}

func NewMysqlSetting(section *ini.Section) *MysqlSetting {
	return &MysqlSetting{
		ip:       section.Key("ip").String(),
		port:     section.Key("port").MustInt(),
		user:     section.Key("user").String(),
		password: section.Key("password").String(),
		database: section.Key("database").String(),
		charset:  section.Key("charset").String(),
		show_sql: section.Key("show_sql").MustBool(),
	}
}
