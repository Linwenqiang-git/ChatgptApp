package completion

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

func (client *CompletionClient) CreateCodeCompletion(ctx context.Context) {
	reqest := openai.CompletionRequest{
		Model:  openai.GPT3TextDavinci003,
		Prompt: "写一个冒泡排序算法",
	}
	rep, err := client.OpenAiClient.Client.CreateCompletion(ctx, reqest)
	if err != nil {
		fmt.Println("Create Code Completion Error:" + err.Error())
		return
	}
	fmt.Println("Create Code Completion Success:" + rep.Choices[0].Text)
}
