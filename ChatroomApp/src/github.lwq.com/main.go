package main

import (
	"context"

	. "github.lwq.com/Completion"

	openai "github.com/sashabaranov/go-openai"

	//. "github.lwq.com/Embedding"
	_ "github.lwq.com/FineTunes"
	. "github.lwq.com/Global"

	//. "github.lwq.com/Model"
	. "github.lwq.com/Ws/Server"
)

func createOpenAiClient() *openai.Client {
	client := openai.NewClient(Chatgpt_token)
	return client
}

func main() {
	client := createOpenAiClient()
	Ctx = context.Background()
	CtClient = &CompletionClient{OpenAiClient: *client}
	// embedClient := &EmbeddingClient{OpenAiClient: client}
	// embedClient.Answer_Question(CtClient, "控件")

	//CtClient.CreateTextCompletion(Ctx, "text-embedding-ada-002")
	//completionClient.CreateTextCompletion(ctx, GPT3Dot5Turbo)
	//respStr := CtClient.CreateChatCompletion(Ctx, "excel读取的总行数和实际行数不一致")
	//fmt.Println(respStr)

	// modelClient := &ModelClient{OpenAiClient: *client}
	// modelClient.ModelList(Ctx)

	//fineTuneClient := &FineTuneClient{OpenAiClient: *client}
	//fineTuneClient.FineTuneList(ctx)

	server := CreateServer()
	server.Start()
	select {}
}
