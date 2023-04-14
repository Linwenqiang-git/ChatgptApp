import openai

_openAiKey = ""
def set_api_key():       
    openai.api_key = _openAiKey

def set_openai_key(api_key):
    _openAiKey = api_key