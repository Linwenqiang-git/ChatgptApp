import tiktoken
from transformers import GPT2TokenizerFast
import pandas as pd

class Tokens:
    __cl100k_base_tokenizer = "cl100k_base"
    def __init__(self):
        
        pass
    
    def generate_tokens(self,source_csv_file: str,columns:list[str]):
        #Load the cl100k_base tokenizer which is designed to work with the ada-002 model
        #tokenizer = tiktoken.get_encoding(self.__cl100k_base_tokenizer)
        tokenizer = GPT2TokenizerFast.from_pretrained("gpt2")        
        df = pd.read_csv(source_csv_file, encoding ='gb18030')    
        df = df.loc[:, ~df.columns.str.contains('^Unnamed')]
        df.columns = columns
        df['tokens'] = df.content.apply(lambda x: len(tokenizer(x)['input_ids']))
        #df.tokens.hist()
        #最新的嵌入模型可以处理多达 8191 个输入标记的输入,如果tokens数量过大，需要拆分文本
        df.to_csv(source_csv_file,index=False)
        pass

    #计算文本tokens
    def calculate_tokens(self,text:str)->int:
        tokenizer = GPT2TokenizerFast.from_pretrained("gpt2")        
        result = tokenizer(text)['input_ids']
        return len(result)        