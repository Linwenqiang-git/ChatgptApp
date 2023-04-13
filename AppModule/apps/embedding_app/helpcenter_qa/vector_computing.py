import numpy as np
import pandas as pd
import tiktoken

import time
from .vectorcache import VectorCache


class VectorComputing:
    __max_section_length = 500    
    __embedding_model = "text-embedding-ada-002"
    __separator_len = 0
    __separator = "\n* "    
    __vectorCache = VectorCache()
    def __init__(self):
        ENCODING = "gpt2"  # encoding for text-davinci-003
        encoding = tiktoken.get_encoding(ENCODING)        
        self.__separator_len = len(encoding.encode(self.__separator))        

    #内部方法  
        
    # 返回两个向量之间的相似度
    def __vector_similarity(self,x: list[float], y: list[float]) -> float:    
        return np.dot(np.array(x), np.array(y))
    
    #排序返回,搜索内容与文档嵌入最相似的部分
    def __order_document_sections_by_query_similarity(self,query: str, context_embeddings: dict[(str, str), np.array]) -> list[(float, (str, str))]:    
        query_embedding = self.create_embedding(query,self.__embedding_model)        
        document_similarities = sorted([
            (self.__vector_similarity(query_embedding, doc_embedding), doc_index) for doc_index, doc_embedding in context_embeddings.items()
        ], reverse=True)    
        return document_similarities

    #计算原始文件的嵌入向量
    def compute_source_doc_embedding(self,df: pd.DataFrame) -> dict[tuple[str, str], list[float]]:  
        result = {}
        call_count_minute = 0
        T1 = time.perf_counter()
        for idx, r in df.iterrows():
            #接口调用上限为每分钟60次
            if call_count_minute == 59:
                call_count_minute = 0
                T2 = time.perf_counter()
                time_diff = T2 - T1
                if time_diff < 60:
                    time.sleep(60.00-time_diff+5)
                T1 = time.perf_counter()
            if r.tokens < 8100:
                result[idx] = self.create_embedding(r.content,self.__embedding_model)
                call_count_minute+=1
        return result
        
    #从指定的csv文件，读取嵌入向量信息
    def load_embeddings(self,csv_file_name: str) -> dict[tuple[str, str], list[float]]: 
        df = pd.read_csv(csv_file_name)
        max_dim = max([int(c) for c in df.columns if c != "title" and c != "small_class"])
        return {
            (r.small_class,r.title): [r[str(i)] for i in range(max_dim + 1)] for _, r in df.iterrows()
        }        
    
    #为指定的文本生成嵌入向量
    def create_embedding(self,text: str, model: str) -> list[float]:
        vector = self.__vectorCache.addOrGetEmbeddingCache(text)
        return vector
    
    #构建问题提示
    def build_prompt(self,question: str, context_embeddings: dict, df: pd.DataFrame) -> str:
        # 根据问题获取相似嵌入
        most_relevant_document_sections = self.__order_document_sections_by_query_similarity(question, context_embeddings)
        
        chosen_sections = []
        chosen_sections_len = 0
        chosen_sections_indexes = []
        
        for _, section_index in most_relevant_document_sections:
            document_section = df[(df['small_class']==section_index[0]) & (df['title']==section_index[1])]
            if document_section.empty:
                continue
            
            chosen_sections_len += int(document_section['tokens']) + self.__separator_len
            if chosen_sections_len > self.__max_section_length:
                break
                
            if len(document_section['content'].values) > 0:
                chosen_sections.append(self.__separator + document_section['content'].values[0].replace("\n", " "))
            chosen_sections_indexes.append(str(section_index))
                
        print(f"Selected {len(chosen_sections)} document sections:")
        print("\n".join(chosen_sections_indexes))
        
        #基于嵌入上下文来回答问题，如果上下文不理解，则回答 “我不知道”
        header = """Answer the question as truthfully as possible using the provided context,and if the answer is not contained within the text below,say "我不知道."\n\nContext:\n"""
        return header + ''.join(chosen_sections) + "\n\n Q: " + question + "\n A:"   
