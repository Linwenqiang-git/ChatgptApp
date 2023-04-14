package configs

import (
	"fmt"
	"sync"
	. "github.com/lwq/configs/settings"
	"github.com/spf13/viper"
)

var configure Configure
var once sync.Once

type Configure struct {
	OpenaiSetting    *OpenaiSetting
	MysqlSetting     *MysqlSetting
	IpcServerSetting *IpcServerSetting
}

func ProvideConfigure() (Configure, error) {
	once.Do(func() {
		viper.SetConfigName("appsettings")    // 配置文件名
		viper.SetConfigType("ini")            // 配置文件类型
		viper.AddConfigPath("D:\\traindatas") // 配置文件路径(需配置本地机密文件地址)
		// 加载配置文件
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("fatal error config file: %s", err))
		}
		//init setting
		openaiSetting := NewOpenaiSetting()
		mysqlSetting := NewMysqlSetting()
		ipcServerSetting := NewIpcServerSetting()
		configure = Configure{
			OpenaiSetting:    openaiSetting,
			MysqlSetting:     mysqlSetting,
			IpcServerSetting: ipcServerSetting,
		}
	})
	return configure, nil
}
