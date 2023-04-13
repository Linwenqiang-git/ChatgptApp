import abc
import typing

class PipeClientBase(metaclass=abc.ABCMeta):

    @abc.abstractmethod
    def connect(self, pipe_name) -> typing.NoReturn:
        '''
        打开并连接管道
        '''
    @abc.abstractmethod
    def write(self, data) -> typing.NoReturn:
        '''
        向管道写入数据
        * @param data, 需要写入的数据
        '''

    @abc.abstractmethod
    def read(self) -> typing.Any:
        '''
        从管道读取数据
        * @param return json 格式数据
        '''
    @abc.abstractmethod
    def close(self) -> typing.NoReturn:
        '''
        关闭管道

        '''
