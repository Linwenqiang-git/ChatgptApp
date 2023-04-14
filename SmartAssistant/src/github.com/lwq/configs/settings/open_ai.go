package settings

import "github.com/spf13/viper"

type OpenaiSetting struct {
	apiKey string
}

func (s *OpenaiSetting) GetOpenaiSetting() string {
	return s.apiKey
}

func NewOpenaiSetting() *OpenaiSetting {
	return &OpenaiSetting{
		apiKey: viper.GetString("openai_settings.api_key"),
	}
}
