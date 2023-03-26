import os
import sys

from chat import continuous_dialogue, initial_dialog_context
root_path = os.path.abspath(__file__)
root_path = '/'.join(root_path.split('/')[:-2])
sys.path.append(root_path)
from API.python.completion import CompletionCient
from embedding import Embedding
from calculatingdata import CalculatingData

#源文件地址
source_file = r"C:\Users\Administrator.DESKTOP-6C71Q8K\Desktop\kb_automation.csv"
#嵌入向量文件地址
embedding_vector_file_name = r"C:\Users\Administrator.DESKTOP-6C71Q8K\Desktop\kb_embedding.csv"
#原始文件头部
source_columns = ['category','small_class','title','content']
#嵌入文件头部
embedding_columns = ['small_class','title']

def main():
    #自定义文件注意csv文件头部名称需要和代码匹配    
    calculate = CalculatingData(source_file,embedding_vector_file_name)
    embedding = Embedding()
    '''
    1.读取准备好的嵌入文件    
    '''    
    #calculate.build_tokens_for_source(source_columns)    
    df = calculate.load_source_file_to_dataframe()
    '''
    2.生成嵌入向量文件    
    '''
    #calculate.generate_embedding_vector_file(embedding,df,embedding_columns)            
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
    question = "领猫SCM是什么？" 
    prompt = embedding.build_prompt(question,document_embeddings,df)
    print("\n===============================\n",prompt.encode('gbk', 'ignore').decode('gbk'))   
    '''
    5.使用text-davinci-003 model 回答用户问题
    '''
    cmpletionCient = CompletionCient()
    COMPLETIONS_API_PARAMS = {    
        "temperature": 0.0,
    }
    response_text = cmpletionCient.create_completion(prompt=prompt,**COMPLETIONS_API_PARAMS)        
    print(response_text)
    #开始连续对话部分
    initial_dialog_context(prompt,response_text)
    user_questions = ['怎样入库？','入库后怎么盘点？','单据怎么上传？']
    for question in user_questions:
        print("\n user_question:",question,"\n")        
        sys_answer = continuous_dialogue(question)
        print(" sys_answer:",sys_answer)        

if __name__ == "__main__":
    main()        