from apps.api.embedding import Embedding
import os

#向量缓存
class VectorCache:    
    __embeddingCacheFile = ""
    __embeddingCacheDic = {}
    __embeddingClient = None

    def __init__(self):         
        # todo:路径获取不优雅       
        module_dir = os.path.dirname(os.path.abspath(__name__)).replace("SmartAssistant\\src\\github.com\\lwq\\cmd","")
        self.__embeddingCacheFile = module_dir + "/AppModule/apps/embedding_app/helpcenter_qa/files/vector_cache_file.txt"
        self.__loadEmbeddingCache()
        self.__embeddingClient = Embedding()

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
            result = self.__embeddingClient.create_embedding(text)
            vector = result["data"][0]["embedding"]
            self.__embeddingCacheDic[text] = vector
            self.__writeEmbeddingCacheFile(text, vector)
        return vector