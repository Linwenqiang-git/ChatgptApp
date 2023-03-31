package handle

import (
	"context"

	provider "github.com/lwq/configs/wire"
)

// Interact with chatgpt
func HandleWsMessgae(user string, message string) (string, error) {
	chatClient, err := provider.GetChatCompletion()
	if err != nil {
		return "", err
	}
	server_msg, err := chatClient.CreateChatCompletionWithContext(context.Background(), message, user)
	if err != nil {
		return "", err
	}
	return server_msg, nil
}
