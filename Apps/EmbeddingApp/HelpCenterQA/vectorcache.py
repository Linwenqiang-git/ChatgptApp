import openai

#向量缓存
class VectorCache:
    __embedding_model = "text-embedding-ada-002"
    __embeddingCacheFile = "./EmbeddingApp/Q_A/Files/vector_cache_file.txt"
    __embeddingCacheDic = {}

    def __init__(self):
        self.__loadEmbeddingCache()

    def __loadEmbeddingCache(self):        
        with open(self.__embeddingCacheFile, 'r') as f:
            for line in f:
                parts = line.strip().split('\t')
                if len(parts) == 2:
                    self.__embeddingCacheDic[parts[0]] = [float(x) for x in parts[1].split(',')]
        

    def __writeEmbeddingCacheFile(self,text:str,vector:list[float]):
        with open(self.__embeddingCacheFile, 'a') as f:
            f.write(f"{text}\t{','.join(str(x) for x in vector)}\n")
        

    def addOrGetEmbeddingCache(self,text:str) ->list[float]:
        vector = self.__embeddingCacheDic.get(text,[])
        if len(vector) == 0:
            #计算文本的嵌入向量
            result = openai.Embedding.create(model=self.__embedding_model,input=text)
            vector = result["data"][0]["embedding"]
            self.__embeddingCacheDic[text] = vector
            self.__writeEmbeddingCacheFile(text, vector)
        return vector