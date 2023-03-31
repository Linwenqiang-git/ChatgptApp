package file

import (
	"fmt"
	"os"
)

func WriteCsv(content []string) {
	//创建文件
	f, err := os.Create("./DataFiles/embeddingData.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	// 写入UTF-8 BOM
	f.WriteString("\xEF\xBB\xBF")
	//创建一个新的写入文件流
	// w := csv.NewWriter(f)
	// w.WriteAll(data)
	// w.Flush()
}
