try:
    import socket
    import sys
    from scheduling.ipc import robot
    from utils.logger import logger
    from utils.ipc_response import ExceptionResponse
    from apps.app_manage import sendRequet
except Exception as e:
    logger.error(f"import module err:{e}")
    sys.exit(0)
    
def main():
    try:       
        robot.connect()
        logger.info("connect to server...")
    except Exception as e:
        logger.error(f"connect to server err:{e}")
        exit(0)

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
            response = sendRequet(request)  
            logger.info(f"response info：{response}")    
            robot.send(response) 
        except socket.error as e:
            logger.error(f"socket error:{e}")
            break
        except Exception as e:            
            if str(e) == "read msg is None":
                break
            else:
                logger.error(f"handle message err:{e}")
                res = ExceptionResponse(request['Id'],str(e))
                robot.send(res) 
    try:
        # 关闭连接
        robot.close()
        logger.info("exit...bye")
    except Exception as e:
        logger.error("exit robot err:",e)
    finally:
        sys.exit(0)

if __name__ == "__main__":
    main()