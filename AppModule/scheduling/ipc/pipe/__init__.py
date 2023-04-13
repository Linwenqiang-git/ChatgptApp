import sys

def create_new_pipe():
    if is_win32():
        from .pipeclient_win import PipeClientWin        
        return  PipeClientWin()
    if is_macos():
        from .pipeclient_mac import PipeClientMac        
        return  PipeClientMac()
    if is_linux():
        from .pipeclient_linux import PipeClientLinux        
        return PipeClientLinux()
    return None

def is_win32():
    return (sys.platform == "win32")

def is_macos():
    return (sys.platform == "darwin")

def is_linux():
    return (sys.platform == "linux")