package embedding

import (
	"fmt"
	"os"
	"strings"

	. "github.lwq.com/Completion"
	. "github.lwq.com/Global"

	"github.com/sashabaranov/go-openai"
)

type EmbeddingClient struct {
	OpenAiClient *openai.Client
}

func readEmbeddingFile() []string {
	content, err := os.ReadFile("./DataFiles/embeddingData.txt")
	if err != nil {
		panic(err)
	}
	contentStr := string(content)
	return strings.Split(contentStr, "\n")
}

func (e *EmbeddingClient) BuildEmbedding() []openai.Embedding {
	input := readEmbeddingFile()
	request := openai.EmbeddingRequest{
		Input: input,
		Model: openai.AdaEmbeddingV2,
	}
	rep, err := e.OpenAiClient.CreateEmbeddings(Ctx, request)
	if err != nil {
		fmt.Println("CreateEmbeddings error：" + err.Error())
		return nil
	}
	//训练文本的结果向量
	embeddings := rep.Data
	return embeddings
}

func (e *EmbeddingClient) distances_from_embeddings(q_embedding []float32, contextEmbeddings []openai.Embedding) {

}

func (e *EmbeddingClient) createContext(question string, contextEmbeddings []openai.Embedding) string {
	request := openai.EmbeddingRequest{
		Input: []string{question},
		Model: openai.AdaEmbeddingV2,
	}
	rep, err := e.OpenAiClient.CreateEmbeddings(Ctx, request)
	if err != nil {
		fmt.Println("Create Question Embeddings error：" + err.Error())
		return question
	}
	//构建问题向量表
	q_embedding := rep.Data[0].Embedding
	e.distances_from_embeddings(q_embedding, contextEmbeddings)

	//returns := make([]string, 0)
	//return strings.Join(returns, "\n\n###\n\n")
	return question
}

func (e *EmbeddingClient) Answer_Question(ctClient *CompletionClient, question string) {
	//构建嵌入索引
	embeddings := e.BuildEmbedding()
	if embeddings != nil {
		//构建问题上下文
		context := e.createContext(question, embeddings)
		//使用text-davinci-003模型
		prompt := "Answer the question based on the context below,and if the question can't be answered based on the context, say \"I don't know\"\n\nContext: " + context + "\n\n---\n\nQuestion: " + question + "\nAnswer:"

		ctClient.CreateTextCompletion(Ctx, prompt, openai.GPT3TextDavinci003)
	}

}
