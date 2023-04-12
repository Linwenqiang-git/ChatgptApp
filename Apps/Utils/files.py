import chardet


def GetFileEncoding(file_name:str) -> str:
    with open(file_name, 'rb') as f:
        result = chardet.detect(f.read())    
        return result['encoding']