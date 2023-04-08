import os
import sys
root_path = os.path.abspath(__file__)
root_path = '/'.join(root_path.split('/')[:-2])
sys.path.append(root_path)
from API.python.completion import CompletionCient


messages = []
cmpletionCient = CompletionCient()

#初始化对话
def __initial_dialog_context(user_message,system_answer):
    if user_message == "":
        return        
    dialogs = [{"role": "user", "content": user_message},
               {"role": "system", "content": system_answer}]
    messages.extend(dialogs)
    
#连续对话
def __continuous_dialogue(user_message) -> str:
    dialog = {"role": "user", "content": user_message}
    messages.append(dialog)
    COMPLETIONS_API_PARAMS = {     
        "temperature": 0.2,
    }
    sys_answer = cmpletionCient.create_chat_completion(messages,**COMPLETIONS_API_PARAMS)
    messages.append(sys_answer)
    return sys_answer['content']    

def call_continuous_dialogue(prompt,response_text):
    is_continue_chat = False
    if is_continue_chat:
        __initial_dialog_context(prompt,response_text)
        user_questions = ['怎样入库？','入库后怎么盘点？','单据怎么上传？']
        for question in user_questions:
            print("\n user_question:",question,"\n")        
            sys_answer = __continuous_dialogue(question)
            print(" sys_answer:",sys_answer)       