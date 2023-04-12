from .pipeclient_nix import PipeClientNix

class PipeClientLinux(PipeClientNix):
    def __init__(self):
        super(PipeClientNix, self).__init__()
