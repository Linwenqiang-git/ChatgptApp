package completion

import (
	"context"
	"fmt"

	. "github.com/lwq/third_party/chatgpt"
	"github.com/sashabaranov/go-openai"
)

type ChatCompletionError struct {
	message string
	err     error
}

func (e *ChatCompletionError) Error() string {
	return fmt.Sprintf("%s: %v", e.message, e.err)
}

type CompletionClient struct {
	OpenAiClient
}
type ChatCompletionClient struct {
	OpenAiClient
	userContext map[string][]openai.ChatCompletionMessage
}

func ProvideCompletionClient(openAiClient OpenAiClient) (CompletionClient, error) {
	return CompletionClient{OpenAiClient: openAiClient}, nil
}
func ProvideChatCompletionClient(openAiClient OpenAiClient) (ChatCompletionClient, error) {
	return ChatCompletionClient{
		OpenAiClient: openAiClient,
		userContext:  make(map[string][]openai.ChatCompletionMessage),
	}, nil
}

// 聊天模式下支持的模型类型
const (
	GPT3Dot5Turbo0301 = openai.GPT3Dot5Turbo0301
	GPT3Dot5Turbo     = openai.GPT3Dot5Turbo
	MaxMessageLength  = 1024 // 假设最大长度为1024
)

// 创建聊天
func (c *CompletionClient) CreateChatCompletion(ctx context.Context, userInput string) (string, error) {
	if len(userInput) == 0 || len(userInput) > MaxMessageLength {
		return "", &ChatCompletionError{"invalid user input", nil}
	}
	request := openai.ChatCompletionRequest{
		Model: GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userInput,
			},
		},
	}
	response, err := c.OpenAiClient.Client.CreateChatCompletion(ctx, request)
	if err != nil {
		return "", &ChatCompletionError{"failed to create chat completion", err}
	}
	if len(response.Choices) == 0 {
		return "", &ChatCompletionError{"no chat message returned", nil}
	}

	return response.Choices[0].Message.Content, nil
}

// 创建上下文聊天
func (c *ChatCompletionClient) CreateChatCompletionWithContext(ctx context.Context, userInput string, userName string) (string, error) {
	if len(userInput) == 0 || len(userInput) > MaxMessageLength {
		return "", &ChatCompletionError{"invalid user input", nil}
	}
	if c.userContext == nil {
		c.userContext = make(map[string][]openai.ChatCompletionMessage)
	}
	if userName == "" {
		userName = "Mr.zhang"
	}
	historyMessages, ok := c.userContext[userName]
	if !ok {
		c.userContext[userName] = []openai.ChatCompletionMessage{}
	}
	userMessage := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userInput,
	}
	historyMessages = append(historyMessages, userMessage)
	request := openai.ChatCompletionRequest{
		Model:    GPT3Dot5Turbo,
		Messages: historyMessages,
	}
	response, err := c.OpenAiClient.Client.CreateChatCompletion(ctx, request)
	if err != nil {
		return "", &ChatCompletionError{"failed to create chat completion", err}
	}
	if len(response.Choices) == 0 {
		return "", &ChatCompletionError{"no chat message returned", nil}
	}
	serverMessage := openai.ChatCompletionMessage{
		Role:    response.Choices[0].Message.Role,
		Content: response.Choices[0].Message.Content,
	}
	historyMessages = append(historyMessages, serverMessage)
	c.userContext[userName] = historyMessages
	return serverMessage.Content, nil
}
