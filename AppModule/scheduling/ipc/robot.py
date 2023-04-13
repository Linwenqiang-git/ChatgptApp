from .retry import Retry
import threading
from .pipe import create_new_pipe

_pipe = create_new_pipe()
_lock = threading.Lock()


def connect():    
    for _ in Retry(20, error_message='无法连接到Robot'):
        try:
            # 连接到IPC服务器
            _pipe.connect()
            break
        except:
            pass

def send(response) -> bool:
    with _lock:
        _pipe.write(response)


def recv():
    with _lock:
        request = _pipe.read()
        if request is None:
            raise Exception("read msg is None")        
        return request
    
def close():
    _pipe.close()