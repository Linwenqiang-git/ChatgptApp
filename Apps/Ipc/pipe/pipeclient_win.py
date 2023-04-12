from .pipeclient_nix import PipeClientNix
import socket

class PipeClientWin(PipeClientNix):
    def __init__(self):
        self._sock = None
        self._address = ("127.0.0.1",2245)
    def connect(self):        
        self._sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._sock.connect(self._address)