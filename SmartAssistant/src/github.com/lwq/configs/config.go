package configs

import (
	"path"
	"runtime"

	. "github.com/lwq/configs/settings"
	"gopkg.in/ini.v1"
)

type Configure struct {
	OpenaiSetting    *OpenaiSetting
	MysqlSetting     *MysqlSetting
	IpcServerSetting *IpcServerSetting
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
	cfg, err := ini.Load(absPath + "/appsettings.ini")
	if err != nil {
		panic(err)
	}
	//init setting
	openaiSetting := NewOpenaiSetting(cfg.Section("OpenaiSettings"))
	mysqlSetting := NewMysqlSetting(cfg.Section("Mysql"))
	ipcServerSetting := NewIpcServerSetting(cfg.Section("IpcServer"))
	return Configure{
		OpenaiSetting:    openaiSetting,
		MysqlSetting:     mysqlSetting,
		IpcServerSetting: ipcServerSetting,
	}, err
}
