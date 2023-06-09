// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/lwq/configs"
	"github.com/lwq/internal/data"
	"github.com/lwq/internal/data/repo"
	"github.com/lwq/third_party/chatgpt"
	"github.com/lwq/third_party/chatgpt/completion"
)

// Injectors from wire.go:

func GetConfigure() (configs.Configure, error) {
	configure, err := configs.ProvideConfigure()
	if err != nil {
		return configs.Configure{}, err
	}
	return configure, nil
}

func GetCompletionClient() (completion.CompletionClient, error) {
	configure, err := configs.ProvideConfigure()
	if err != nil {
		return completion.CompletionClient{}, err
	}
	openAiClient, err := chatgpt.ProvideOpenAiClient(configure)
	if err != nil {
		return completion.CompletionClient{}, err
	}
	completionClient, err := completion.ProvideCompletionClient(openAiClient)
	if err != nil {
		return completion.CompletionClient{}, err
	}
	return completionClient, nil
}

func GetChatCompletion() (completion.ChatCompletionClient, error) {
	configure, err := configs.ProvideConfigure()
	if err != nil {
		return completion.ChatCompletionClient{}, err
	}
	openAiClient, err := chatgpt.ProvideOpenAiClient(configure)
	if err != nil {
		return completion.ChatCompletionClient{}, err
	}
	chatCompletionClient, err := completion.ProvideChatCompletionClient(openAiClient)
	if err != nil {
		return completion.ChatCompletionClient{}, err
	}
	return chatCompletionClient, nil
}

func GetDbContext() (data.DbContext, error) {
	configure, err := configs.ProvideConfigure()
	if err != nil {
		return data.DbContext{}, err
	}
	dbContext, err := data.ProvideDbContext(configure)
	if err != nil {
		return data.DbContext{}, err
	}
	return dbContext, nil
}

func GetMessageRepo() (repo.MessageRepo, error) {
	configure, err := configs.ProvideConfigure()
	if err != nil {
		return repo.MessageRepo{}, err
	}
	dbContext, err := data.ProvideDbContext(configure)
	if err != nil {
		return repo.MessageRepo{}, err
	}
	messageRepo := repo.ProvideMessageRepo(dbContext)
	return messageRepo, nil
}

// wire.go:

var configureSet = wire.NewSet(configs.ProvideConfigure)

var openAiClientSet = wire.NewSet(configureSet, chatgpt.ProvideOpenAiClient)

var completionSet = wire.NewSet(openAiClientSet, completion.ProvideCompletionClient, completion.ProvideChatCompletionClient)

var dbContextSet = wire.NewSet(configureSet, data.ProvideDbContext)

var repoSet = wire.NewSet(dbContextSet, repo.ProvideMessageRepo)
