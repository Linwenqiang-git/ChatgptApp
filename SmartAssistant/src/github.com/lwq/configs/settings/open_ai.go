package settings

import "gopkg.in/ini.v1"

type OpenaiSetting struct {
	apiKey string
}

func (s *OpenaiSetting) GetOpenaiSetting() string {
	return s.apiKey
}

func NewOpenaiSetting(section *ini.Section) *OpenaiSetting {
	return &OpenaiSetting{
		apiKey: section.Key("api_key").String(),
	}
}
