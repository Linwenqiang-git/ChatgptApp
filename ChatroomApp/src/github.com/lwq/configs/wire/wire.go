//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	conf "github.com/lwq/internal/conf"
	chatgpt "github.com/lwq/third_party/chatgpt"
	compleion "github.com/lwq/third_party/chatgpt/completion"
)

var OpenAiClientSet = wire.NewSet(conf.ProvideConfigure, chatgpt.ProvideOpenAiClient)

var CompletionSet = wire.NewSet(OpenAiClientSet, compleion.ProvideCompletionClient, compleion.ProvideChatCompletionClient)

func GetCompletionClient() (compleion.CompletionClient, error) {
	wire.Build(CompletionSet)
	return compleion.CompletionClient{}, nil
}

func GetChatCompletion() (compleion.ChatCompletionClient, error) {
	wire.Build(CompletionSet)
	return compleion.ChatCompletionClient{}, nil
}
