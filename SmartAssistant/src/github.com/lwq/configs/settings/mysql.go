package settings

import (
	"strconv"
	"strings"

	"github.com/spf13/viper"
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

func NewMysqlSetting() *MysqlSetting {
	return &MysqlSetting{
		ip:       viper.GetString("mysql.ip"),
		port:     viper.GetInt("mysql.port"),
		user:     viper.GetString("mysql.user"),
		password: viper.GetString("mysql.password"),
		database: viper.GetString("mysql.database"),
		charset:  viper.GetString("mysql.charset"),
		show_sql: viper.GetBool("mysql.show_sql"),
	}
}
