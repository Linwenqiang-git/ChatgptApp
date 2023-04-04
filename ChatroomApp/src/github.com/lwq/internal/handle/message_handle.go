package handle

import (
	"context"

	provider "github.com/lwq/configs/wire"
	. "github.com/lwq/internal/data/dto"
	. "github.com/lwq/internal/data/entity"
	_ "github.com/lwq/internal/data/repo"
)

type IMessageRepo interface {
	CreateUser(user User) (int, error)
	GetUser(userName string) (User, error)
	InserUserMessage(records []ChatRecord) (int, error)
	GetUserHistory(userName string) ([]UserHistoryDto, error)
}

var (
	_messageRepo IMessageRepo
)

func init() {
	if _messageRepo == nil {
		var err error
		_messageRepo, err = provider.GetMessageRepo()
		if err != nil {
			panic("Init MessageRepo error" + err.Error())
		}
	}
}

// Interact with chatgpt
func HandleWsMessgae(userName string, message string) (string, error) {
	chatClient, err := provider.GetChatCompletion()
	if err != nil {
		return "", err
	}
	if !chatClient.IsInitChatContext(userName) {
		messages, err := _messageRepo.GetUserHistory(userName)
		if err != nil {
			return "", err
		}
		chatClient.InitChatContext(userName, messages)
	}
	defer func() {
		records := chatClient.GetNewChatContext(userName)
		_messageRepo.InserUserMessage(records)
	}()
	server_msg, err := chatClient.CreateChatCompletionWithContext(context.Background(), message, userName)
	if err != nil {
		return "", err
	}
	return server_msg, nil
}
