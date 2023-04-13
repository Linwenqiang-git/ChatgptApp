def ExceptionResponse(requestId,errMsg:str) ->dict:
    response = {
        "ResponseId":requestId,
        "Code":500,
        "Message":"",
        "Error":Exception(errMsg),
        "ErrorMsg":errMsg
    }   
    return response