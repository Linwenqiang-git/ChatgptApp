import os
import sys
root_path = os.path.abspath(__file__)
root_path = '/'.join(root_path.split('/')[:-2])
sys.path.append(root_path)
from API.python.completion import CompletionCient


messages = []
cmpletionCient = CompletionCient()

#初始化对话
def initial_dialog_context(user_message,system_answer):
    dialogs = [{"role": "user", "content": user_message},
               {"role": "system", "content": system_answer}]
    messages.extend(dialogs)
    
#连续对话
def continuous_dialogue(user_message) -> str:
    dialog = {"role": "user", "content": user_message}
    messages.append(dialog)
    COMPLETIONS_API_PARAMS = {     
        "temperature": 0.2,
    }
    sys_answer = cmpletionCient.create_chat_completion(messages,**COMPLETIONS_API_PARAMS)
    messages.append(sys_answer)
    return sys_answer['content']    