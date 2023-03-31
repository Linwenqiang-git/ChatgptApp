package completion

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

func (client *CompletionClient) CreateTextCompletion(ctx context.Context, prompt string, fineTuneModel string) {
	reqest := openai.CompletionRequest{
		Model:  fineTuneModel,
		Prompt: prompt,
	}
	rep, err := client.OpenAiClient.Client.CreateCompletion(ctx, reqest)
	if err != nil {
		fmt.Println("CreateCompletion Error:" + err.Error())
		return
	}
	fmt.Println("CreateCompletion Success:" + rep.ID)
	fmt.Println("CreateCompletion.Object:" + rep.Object)
	fmt.Println("==================")
	if rep.Choices != nil && len(rep.Choices) > 0 {
		for _, choice := range rep.Choices {
			fmt.Println("Text：" + choice.Text)
			fmt.Println("FinishReason:" + choice.FinishReason)
		}
	}
}
