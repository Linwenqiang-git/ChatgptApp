package fine_tunes

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

type FineTuneClientInterface interface {
	CreateFineTune(context.Context) (string, error)
	UploadFineTuneFile(context.Context, string, string) (string, error)
}
type FineTuneClient struct {
	OpenAiClient openai.Client
}

func init() {
	fmt.Println("this is fine_tunes init func")
}

// 创建 fine_tune 的基础 Model
const (
	ada     = "ada"
	babbage = "babbage"
	curie   = "curie"
	davinci = "davinci"
)

// 上传微调文件
func (c *FineTuneClient) uploadFineTuneFile(ctx context.Context, fileName string, filePath string) (fileId string, err error) {
	request := openai.FileRequest{
		FileName: fileName,
		FilePath: filePath,
		Purpose:  "fine-tune",
	}
	response, err := c.OpenAiClient.CreateFile(ctx, request)
	if err == nil {
		fileId = response.ID
	}
	return
}

func (c *FineTuneClient) CreateFineTune(ctx context.Context) (fineTuenModel string, err error) {
	fileId, err := c.uploadFineTuneFile(ctx, "@fine_tunes_data.jsonl", "D:\\GoStudy\\src\\github.linwenqiang.com\\OpenAi\\fine_tunes_data.jsonl")
	if err != nil {
		fmt.Printf("UploadFineTuneFile error:%s", err.Error())
		return
	}
	request := openai.FineTuneRequest{
		TrainingFile: fileId,
		Model:        davinci,
	}
	rep, err := c.OpenAiClient.CreateFineTune(ctx, request)
	if err != nil {
		fmt.Println("CreateFineTune error：" + err.Error())
	}
	fmt.Println("CreateFineTune rep：")
	fmt.Println("FineTunedModel：" + rep.FineTunedModel)
	fmt.Println("Model：" + rep.Model)
	fmt.Println("Status：" + rep.Status)
	fmt.Println("ID：" + rep.ID)
	fmt.Println("Object：" + rep.Object)

	fineTuenModel = rep.FineTunedModel

	return
}

func (c *FineTuneClient) FineTuneList(ctx context.Context) {
	rep, err := c.OpenAiClient.ListFineTunes(ctx)
	if err != nil {
		fmt.Println("FineTuneList error：", err.Error())
	}
	fmt.Println("ListFineTunes Rep:")
	data := rep.Data
	for _, tune := range data {
		fmt.Print("Tune ID:" + tune.ID + " ")
		fmt.Print("Tune Object:" + tune.Object + " ")
		fmt.Print("Tune Model:" + tune.Model + " ")
		fmt.Print("Tune FineTunedModel:" + tune.FineTunedModel + " ")
		fmt.Print("Tune Status:" + tune.Status + " \n")
	}
}
