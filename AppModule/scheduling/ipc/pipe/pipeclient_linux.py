import json
import struct

from .pipeclient_base import PipeClientBase
from utils.logger import logger

import socket

class PipeClientLinux(PipeClientBase):
    def __init__(self):
        self._sock = None
        self._address = "/tmp/robotsocket.sock"

    def connect(self):        
        self._sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
        self._sock.connect(self._address)        

    def _read_data(self, recv_size):
        try:
            body_part=self._sock.recv(recv_size)
            return body_part
        except:
            print('_read_data except')
            return None
        
    def _write_data(self, message):
        try:
            ret = self._sock.send(message)
            if (ret > 0):
                return True
            else:
                return False
        except Exception as e:
            logger.error(f'_write_data except:{e}')
            return False
        
    def write(self, obj):
        try:            
            content = json.dumps(obj)                
            content_bytes = content.encode('utf-8')            
            #以4字节写入流
            logger.info(len(content_bytes))
            buf = struct.pack('i', len(content_bytes))
            suc = self._write_data(buf)
            if not suc:
                return False
            suc = self._write_data(content_bytes)
            if not suc:
                return False
            return True            
        except Exception as e:
            logger.error(f"write data err:{e}")
            return False

    def read(self):
        size_bytes = self._read_data(4)
        if size_bytes is None:
            return None
        size, = struct.unpack('i', size_bytes)
        content_bytes = self._read_data(size)

        if content_bytes is None:
            return None
        content = content_bytes.decode('utf-8')
        return json.loads(content)

        
    def close(self):
        try:
            if self._sock != None:
                self._sock.close()
                self._sock = None
        except:
            self._sock = None
        
    


