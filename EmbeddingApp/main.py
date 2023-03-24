import openai
from embedding import Embedding
from calculatingdata import CalculatingData


COMPLETIONS_MODEL = "text-davinci-003" #完成模型
source_file = "./EmbeddingApp/helpcenter_source.csv"
embedding_vector_file_name = "./EmbeddingApp/helpcenter_source_embeddings.csv"

def main():
    #自定义文件注意csv文件头部名称需要和代码匹配    
    calculate = CalculatingData(source_file,embedding_vector_file_name)
    embedding = Embedding()
    '''
    1.读取准备好的嵌入文件    
    '''
    #calculate.build_tokens_for_source()
    df = calculate.load_source_file_to_dataframe()
    '''
    2.生成嵌入向量文件    
    '''
    #calculate.generate_embedding_vector_file(embedding,df)    
    '''
    3.将步骤二得到的向量字典写入到新的嵌入向量csv文件
    #（写入文件的目标是方便下次直接使用，不需要每次都计算原始文件的向量）
    # 这里将计算好的结果以字典的形式载入内存
    '''
    document_embeddings = embedding.load_embeddings(embedding_vector_file_name)
    
    '''
    4.构建用户问题的提示
    '''
    # 用户问答部分
    question = "今天吃什么" 
    prompt = embedding.build_prompt(question,document_embeddings,df)
    print("\n===============================\n",prompt.encode('gbk', 'ignore').decode('gbk'))   
    '''
    5.使用text-davinci-003 model 回答用户问题
    '''
    COMPLETIONS_API_PARAMS = {
        # 这里的温度（0-1）参数给到0，如果大于0，ai会自己丰富感情色彩，但是因为这是基于可预测的上下文环境，所以这里使用0即可
        "temperature": 0.0,
        "max_tokens": 300,
        "model": COMPLETIONS_MODEL,
    }
    response = openai.Completion.create(
                prompt=prompt,
                **COMPLETIONS_API_PARAMS
            )
    response_text = response["choices"][0]["text"].strip(" \n")
    print(response_text)

if __name__ == "__main__":
    main()        