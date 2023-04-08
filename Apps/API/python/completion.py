from EmbeddingApp.Q_A.tokens import Tokens
import openai

class CompletionCient:
    __text_completio_model = {
        'text-davinci-003':'text-davinci-003',
        'text-davinci-002':'text-davinci-002'
    }
    __chat_completion_mode ={
        'gpt-3.5-turbo':'gpt-3.5-turbo'
    }
    def __init__(self):
        pass
    
    def create_completion(self,prompt:str,**kwargs) -> str:
        token_client = Tokens()
        max_token = 3000
        prompt_tokens = token_client.calculate_tokens(prompt)
        can_use_max_tokens = max_token - prompt_tokens
        if can_use_max_tokens <= 0:
            raise ValueError('prompt tokens too large:' + str(prompt_tokens))
        response = openai.Completion.create(
                prompt=prompt,
                max_tokens=can_use_max_tokens,
                model=self.__text_completio_model['text-davinci-003'],
                **kwargs)
        #print(response)
        response_text = response["choices"][0]["text"].strip(" \n")
        return response_text
        
    def create_chat_completion(self,messages:list[dict],**kwargs)-> dict:
        rep = openai.ChatCompletion.create(
            model=self.__chat_completion_mode['gpt-3.5-turbo'],
            messages=messages,
            **kwargs)
        message = rep["choices"][0]["message"]
        return message        