package model

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type ModelClient struct {
	OpenAiClient openai.Client
}

func (c *ModelClient) ModelList(ctx context.Context) {
	rep, err := c.OpenAiClient.ListModels(ctx)
	if err != nil {
		fmt.Println("ModelList error：" + err.Error())
		return
	}
	fmt.Println("ModelList rep：")
	for _, model := range rep.Models {
		fmt.Println(" ID ：" + model.ID)
		fmt.Println(" Object ：" + model.Object)
		fmt.Println(" Parent ：" + model.Parent)
		fmt.Print(" CreatedAt ：")
		fmt.Println(model.CreatedAt)
		fmt.Println("======================")
	}
}
