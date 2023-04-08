package chatgpt

import (
	. "github.com/lwq/configs"
	openai "github.com/sashabaranov/go-openai"
)

type OpenAiClient struct {
	Client *openai.Client
}

func ProvideOpenAiClient(configure Configure) (OpenAiClient, error) {
	apiKey := configure.OpenaiSetting.GetOpenaiSetting()
	openai_client := openai.NewClient(apiKey)
	client := OpenAiClient{
		Client: openai_client,
	}
	return client, nil
}
