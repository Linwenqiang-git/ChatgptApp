//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	conf "github.com/lwq/internal/conf"
	dbContext "github.com/lwq/internal/data"
	repo "github.com/lwq/internal/data/repo"
	chatgpt "github.com/lwq/third_party/chatgpt"
	compleion "github.com/lwq/third_party/chatgpt/completion"
)

var openAiClientSet = wire.NewSet(conf.ProvideConfigure, chatgpt.ProvideOpenAiClient)

var completionSet = wire.NewSet(openAiClientSet, compleion.ProvideCompletionClient, compleion.ProvideChatCompletionClient)

var dbContextSet = wire.NewSet(conf.ProvideConfigure, dbContext.ProvideDbContext)

var repoSet = wire.NewSet(dbContextSet, repo.ProvideMessageRepo)

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
