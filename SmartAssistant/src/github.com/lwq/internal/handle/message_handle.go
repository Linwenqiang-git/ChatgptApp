package handle

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	provider "github.com/lwq/utils/wire"
	entity "github.com/lwq/internal/data/entity"
	_ "github.com/lwq/internal/data/repo"
	. "github.com/lwq/internal/shared/dto"
	"gorm.io/gorm"
)

type IMessageRepo interface {
	CreateUser(user entity.User) (int, error)
	GetUser(userName string) (*entity.User, error)
	InserUserMessage(records []entity.ChatRecord) (int, error)
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
func HandleWsMessgae(userName string, sendChan chan []byte, message string) error {
	_, err := getOrCreateUser(userName)
	if err != nil {
		return err
	}
	chatClient, err := provider.GetChatCompletion()
	if err != nil {
		return err
	}
	if !chatClient.IsInitChatContext(userName) {
		messages, err := _messageRepo.GetUserHistory(userName)
		if err != nil {
			return err
		}
		chatClient.InitChatContext(userName, messages)
	}

	responseStream, err := chatClient.CreateChatCompletionWithContext(context.Background(), message, userName)
	if err != nil {
		return err
	}
	defer responseStream.Close()
	var strBuilder strings.Builder
	for {
		response, err := responseStream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("Stream finished")
			sendChan <- []byte("end")
			break
		}
		if err != nil {
			return err
		}
		content := response.Choices[0].Delta.Content
		strBuilder.WriteString(content)
		sendChan <- []byte(content)
	}
	content := strBuilder.String()
	chatClient.AddContext(content, userName)
	go func() {
		records := chatClient.GetNewChatContext(userName)
		_messageRepo.InserUserMessage(records)
	}()
	return nil
}

func getOrCreateUser(userName string) (*entity.User, error) {
	user, err := _messageRepo.GetUser(userName)
	if err != nil {
		return nil, err
	}
	if user == nil {
		user = &entity.User{
			Name: userName,
			Model: gorm.Model{
				CreatedAt: time.Now(),
			},
		}
		_messageRepo.CreateUser(*user)
		log.Printf("register %s success", user.Name)
	}
	return user, nil
}
