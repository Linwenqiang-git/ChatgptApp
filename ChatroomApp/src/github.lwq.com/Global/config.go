package global

import (
	"path"
	"runtime"
	"strings"

	"gopkg.in/ini.v1"
)

type GlobalConfig struct {
	OpenaiSetting *OpenaiSetting
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
func LoadGlobalConfig() *GlobalConfig {
	absPath := getCurrentAbPathByCaller()
	if absPath == "" {
		panic("get project path err")
	}
	absPath = strings.Replace(absPath,"ChatroomApp/src/github.lwq.com/Global", "", 1)
	cfg, err := ini.Load(absPath + "configurations.ini")
	if err != nil {
		panic(err)
	}
	print(cfg.Section("OpenaiSettings").Key("api_key").String())
	//init openai setting
	openaiSetting := &OpenaiSetting{
		apiKey: cfg.Section("OpenaiSettings").Key("api_key").String(),
	}
	return &GlobalConfig{OpenaiSetting: openaiSetting}
}

type OpenaiSetting struct {
	apiKey string
}

func (s *OpenaiSetting) GetSetting() string {
	return s.apiKey
}
