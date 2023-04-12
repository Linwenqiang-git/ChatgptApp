//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	conf "github.com/lwq/configs"
	dbContext "github.com/lwq/internal/data"
	repo "github.com/lwq/internal/data/repo"
	chatgpt "github.com/lwq/third_party/chatgpt"
	compleion "github.com/lwq/third_party/chatgpt/completion"
)

var configureSet = wire.NewSet(conf.ProvideConfigure)

var openAiClientSet = wire.NewSet(configureSet, chatgpt.ProvideOpenAiClient)

var completionSet = wire.NewSet(openAiClientSet, compleion.ProvideCompletionClient, compleion.ProvideChatCompletionClient)

var dbContextSet = wire.NewSet(configureSet, dbContext.ProvideDbContext)

var repoSet = wire.NewSet(dbContextSet, repo.ProvideMessageRepo)

func GetConfigure() (conf.Configure, error) {
	wire.Build(configureSet)
	return conf.Configure{}, nil
}

func GetCompletionClient() (compleion.CompletionClient, error) {
	wire.Build(completionSet)
	return compleion.CompletionClient{}, nil
}

func GetChatCompletion() (compleion.ChatCompletionClient, error) {
	wire.Build(completionSet)
	return compleion.ChatCompletionClient{}, nil
}

func GetDbContext() (dbContext.DbContext, error) {
	wire.Build(dbContextSet)
	return dbContext.DbContext{}, nil
}

func GetMessageRepo() (repo.MessageRepo, error) {
	wire.Build(repoSet)
	return repo.MessageRepo{}, nil
}
