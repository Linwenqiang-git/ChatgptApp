from .embedding_app.helpcenter_qa.main import main as HelpCenterQA
from .embedding_app.demand_classification.main import main as DemandClassification
from apps.api.api_keys import set_openai_key
from utils.logger import logger

_moduleDic = {
    1:HelpCenterQA,
    2:DemandClassification,    
}

def process_request(request :dict) -> dict:        
    requestId = request.get('Id',None)
    module = request.get('Module',-1)
    request_msg = request.get('Message',"")
    response = {
        "ResponseId":requestId,
        "Code":200,
        "Message":"",
        "Error":None,
        "ErrorMsg":""
    }   
    if module == -1:
        response["Code"] = 500
        response["ErrorMsg"] = "不支持的模式"    
    elif module == 0:
        #set openai key
        set_openai_key(request_msg)
        logger.info("set openai key success")
        pass
    elif request_msg == "":
        response["Code"] = 500
        response["ErrorMsg"] = "消息不能为空"    
    else:
        try:            
            main = _moduleDic[module]
            response["Message"] = main(request_msg)                 
        except Exception as e:   
            response["Code"] = 500
            response["ErrorMsg"] = "处理消息异常：" + str(e)    
    return response
