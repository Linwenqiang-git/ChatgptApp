import os
import sys
import openai
import configparser

from ...API.python.completion import CompletionCient
from embedding import Embedding
from calculatingdata import CalculatingData

#源文件地址
source_file = r"D:\traindatas\systematic_source_data.csv"
#嵌入向量文件地址
embedding_vector_file_name = r"D:\traindatas\systematic_source_data_embedding.csv"


def setOpenaiKey():   
    config_file = configparser.ConfigParser()     
    config_file.read("./././configurations.ini")        
    openai.api_key = config_file['OpenaiSettings']['api_key']


def main():
    setOpenaiKey()
    #自定义文件注意csv文件头部名称需要和代码匹配    
    calculate = CalculatingData(source_file,embedding_vector_file_name)
    embedding = Embedding()
    '''
    1.计算原始文件的token和嵌入文件的向量值    
    '''      
    df = calculate.calculate_token_vector(is_calculate_token=False,is_calculate_vector=False)
    '''    
    2.将步骤二得到的向量字典写入到新的嵌入向量csv文件
    #（写入文件的目标是方便下次直接使用，不需要每次都计算原始文件的向量）
    # 这里将计算好的结果以字典的形式载入内存
    '''
    document_embeddings = embedding.load_embeddings(embedding_vector_file_name)    
    '''
    3.构建用户问题的提示
    '''
    # 用户问答部分
    question = "领猫SCM是什么？" 
    prompt = embedding.build_prompt(question,document_embeddings,df)
    print("\n===============================\n",prompt.encode('gbk', 'ignore').decode('gbk'))   
    '''
    4.使用text-davinci-003 model 回答用户问题
    '''
    cmpletionCient = CompletionCient()
    COMPLETIONS_API_PARAMS = {    
        "temperature": 0.0,
    }
    response_text = cmpletionCient.create_completion(prompt=prompt,**COMPLETIONS_API_PARAMS)        
    print(response_text)   
          

if __name__ == "__main__":
    main()        