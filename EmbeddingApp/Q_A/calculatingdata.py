from tokens import Tokens
import pandas as pd
from embedding import Embedding

class CalculatingData:
    __source_file_name = ""
    __embedding_vector_file_name = ""
    def __init__(self,source_file_name,embedding_vector_file_name):
        self.__source_file_name = source_file_name
        self.__embedding_vector_file_name = embedding_vector_file_name
        pass

    # 为原始文件计算token
    def build_tokens_for_source(self,columns:list[str]):
        tokensObj = Tokens()
        tokensObj.generate_tokens(self.__source_file_name,columns)
        pass

    #将初始文件载入DataFrame
    def load_source_file_to_dataframe(self) -> pd.DataFrame:
        return pd.read_csv(self.__source_file_name)

    #生成嵌入向量文件
    def generate_embedding_vector_file(self,embedding:Embedding,df:pd.DataFrame,embedding_columns:list[str]):   
        df = df.set_index(embedding_columns)
        source_vector = embedding.compute_source_doc_embedding(df)
        #将计算结果写进新的文件
        self.__generate_embedding_vector_file(source_vector)   

     #计算提tokens数量
    def calculate_prompt_tokens(serlf,prompt:str):
        
        pass 

    #生成嵌入向量文件
    def __generate_embedding_vector_file(self,data:dict[tuple[str, str], list[float]]):    
        #固定csv文件头部为title、heading
        dic = {'small_class':[],'title':[]}    
        for k,vector in data.items():                
            dic['small_class'].append(k[0])
            dic['title'].append(k[1] if k[1] != None else '')
            index = 0
            for v in vector:
                if str(index) in dic:
                    dic[str(index)].append(v)            
                else:
                    dic[str(index)] = [v]                
                index+=1
        # print(dic)
        df = pd.DataFrame(dic)
        df.to_csv(self.__embedding_vector_file_name,index=False) 