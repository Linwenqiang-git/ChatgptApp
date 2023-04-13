from .embedding_app.helpcenter_qa.main import main as HelpCenterQA
from .embedding_app.demand_classification.main import main as DemandClassification
from utils.logger import logger

_moduleDic = {
    1:HelpCenterQA,
    2:DemandClassification,    
}

def sendRequet(request :dict) -> dict:        
    requestId = request['Id']
    module = request.get('Module',0)
    request_msg = request.get('Message',"")
    response = {
        "ResponseId":requestId,
        "Code":200,
        "Message":"",
        "Error":None,
        "ErrorMsg":""
    }   
    if module == 0:
        response["Code"] = 500
        response["ErrorMsg"] = "不支持的模式"    
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
