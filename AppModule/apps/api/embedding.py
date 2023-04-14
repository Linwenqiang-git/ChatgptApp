import openai
from .api_keys import set_api_key

class Embedding:
    __embedding_model = "text-embedding-ada-002"
    def __ini__(self):
        set_api_key()
    def create_embedding(self,text:str) ->dict:
        response = openai.Embedding.create(model=self.__embedding_model,input=text)
        return response