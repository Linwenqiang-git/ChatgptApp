import socket
import sys
from Ipc import robot
from Utils.logger import logger

def main():
    try:        
        robot.connect()
        logger.info("connect to server...")
    except Exception as e:
        logger.error(f"connect to server err:{e}")
        pass

    while True:
        # 接收响应        
        try:   
            logger.info("waiting info...")
            request = robot.recv()        
            if request == None:
                logger.info(f"receive empty info")             
                continue
            logger.info(f"receive info：{request}")         
            if request['IsExit'] == True:
                break
            logger.info(request)
            requestId = request['Id']
            module = request['Module']
            # 根据请求模块调用
            message = "这是python端返回的消息"
            response = {
                "ResponseId":requestId,
                "Code":200,
                "Message":message,
                "Error":None,
                "ErrorMsg":""
            }        
            robot.send(response)        
        except socket.error as e:
            logger.error(f"socket error:{e}")
            break
        except Exception as e:
            #read msg is None
            logger.error(f"handle message err:{e}")
            break
    try:
        # 关闭连接
        robot.close()
        logger.info("exit...bye")
    except Exception as e:
        logger.error("exit riobot err:",e)
    finally:
        sys.exit(0)

if __name__ == "__main__":
    main()