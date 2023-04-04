package conf

import (
	"path"
	"runtime"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
)

type Configure struct {
	OpenaiSetting *OpenaiSetting
	MysqlSetting  *Mysql
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
func ProvideConfigure() (Configure, error) {
	absPath := getCurrentAbPathByCaller()
	if absPath == "" {
		panic("get project path err")
	}
	cfg, err := ini.Load(absPath + "/configurations.ini")
	if err != nil {
		panic(err)
	}
	//init openai setting
	openaiSetting := &OpenaiSetting{
		apiKey: cfg.Section("OpenaiSettings").Key("api_key").String(),
	}
	//init mysql setting
	mysqlSection := cfg.Section("Mysql")
	mysqlSetting := &Mysql{
		ip:       mysqlSection.Key("ip").String(),
		port:     mysqlSection.Key("port").MustInt(),
		user:     mysqlSection.Key("user").String(),
		password: mysqlSection.Key("password").String(),
		database: mysqlSection.Key("database").String(),
		charset:  mysqlSection.Key("charset").String(),
		show_sql: mysqlSection.Key("show_sql").MustBool(),
	}
	return Configure{
		OpenaiSetting: openaiSetting,
		MysqlSetting:  mysqlSetting,
	}, err
}

type OpenaiSetting struct {
	apiKey string
}

func (s *OpenaiSetting) GetOpenaiSetting() string {
	return s.apiKey
}

type Mysql struct {
	ip       string
	port     int
	user     string
	password string
	database string
	charset  string
	show_sql bool
}

func (s *Mysql) GetConnectInfo() string {
	//"user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	var dsn_builder strings.Builder
	dsn_builder.WriteString(s.user)
	dsn_builder.WriteString(":")
	dsn_builder.WriteString(s.password)
	dsn_builder.WriteString("@tcp")
	dsn_builder.WriteString("(")
	dsn_builder.WriteString(s.ip)
	dsn_builder.WriteString(":")
	dsn_builder.WriteString(strconv.Itoa(s.port))
	dsn_builder.WriteString(")/")
	dsn_builder.WriteString(s.database)
	dsn_builder.WriteString("?charset=")
	dsn_builder.WriteString(s.charset)
	dsn_builder.WriteString("&parseTime=True&loc=Local")
	return dsn_builder.String()
}
