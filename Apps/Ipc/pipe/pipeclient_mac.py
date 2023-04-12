from .pipeclient_nix import PipeClientNix

class PipeClientMac(PipeClientNix):
    def __init__(self):
        super(PipeClientNix, self).__init__()