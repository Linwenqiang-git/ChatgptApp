from .pipeclient_linux import PipeClientLinux

class PipeClientMac(PipeClientLinux):
    def __init__(self):
        super(PipeClientLinux, self).__init__()