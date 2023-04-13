from enum import IntEnum, unique

@unique
class PipeErrorCode(IntEnum):
    Timeout = 108 # 操作超时
