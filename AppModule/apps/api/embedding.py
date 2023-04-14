import openai

class Embedding:
    __embedding_model = "text-embedding-ada-002"
    def __ini__(self):
        pass
    def create_embedding(self,text:str) ->dict:
        response = openai.Embedding.create(model=self.__embedding_model,input=text)
        return response