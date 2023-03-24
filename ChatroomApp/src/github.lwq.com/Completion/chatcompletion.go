package completion

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type CompletionClient struct {
	OpenAiClient openai.Client
}

func init() {
	fmt.Println("This is Completion Client")
}

//聊天模式下支持的模型类型
const (
	GPT3Dot5Turbo0301 = openai.GPT3Dot5Turbo0301
	GPT3Dot5Turbo     = openai.GPT3Dot5Turbo
)

func (client *CompletionClient) CreateChatCompletion(ctx context.Context, message string) string {
	request := openai.ChatCompletionRequest{
		Model: GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: message,
			},
		},
	}
	rep, err := client.OpenAiClient.CreateChatCompletion(ctx, request)
	if err != nil {
		return err.Error()
	}
	return rep.Choices[0].Message.Content
}
