import tiktoken
import pandas as pd

class Tokens:
    __cl100k_base_tokenizer = "cl100k_base"
    def __init__(self):
        
        pass
    
    def generate_tokens(self,source_csv_file: str = 'helpcenter_source.csv'):
        #Load the cl100k_base tokenizer which is designed to work with the ada-002 model
        tokenizer = tiktoken.get_encoding(self.__cl100k_base_tokenizer)
        df = pd.read_csv(source_csv_file)    
        df.columns = ['title', 'heading','content']
        df['tokens'] = df.content.apply(lambda x: len(tokenizer.encode(x)))
        #df.tokens.hist()
        #最新的嵌入模型可以处理多达 8191 个输入标记的输入,如果tokens数量过大，需要拆分文本
        df.to_csv(source_csv_file,index=False)
        pass

