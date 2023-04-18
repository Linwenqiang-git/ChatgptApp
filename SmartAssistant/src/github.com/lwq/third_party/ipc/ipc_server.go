package ipc

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	. "github.com/lwq/configs"
	. "github.com/lwq/third_party/ipc/dto"
	. "github.com/lwq/utils/convert"
	. "github.com/lwq/utils/event"
	. "github.com/lwq/utils/process_manager"
	provider "github.com/lwq/utils/wire"
)

type IpcState int

const (
	UnInit IpcState = iota
	Run
	Disposed
)

/*
进程间通信服务端，用于创建管道，负责进程间消息的处理
*/
type IpcServer struct {
	lock              sync.RWMutex
	conn              net.Conn
	state             IpcState
	ln                net.Listener
	OnMsgReveiveEvent *Event
	wg                sync.WaitGroup
	config            Configure
}

func NewIpcServer() *IpcServer {
	config, err := provider.GetConfigure()
	if err != nil {
		panic(err)
	}
	return &IpcServer{
		state:             UnInit,
		OnMsgReveiveEvent: NewEvent(),
		config:            config,
	}
}

func (s *IpcServer) Start() error {
	if s.state == Run {
		return nil
	} else if s.state == Disposed {
		return errors.New("ipc-service is disposed")
	} else {
		s.lock.Lock()
		err := s.createConn()
		if err != nil {
			return err
		}
		s.lock.Unlock()

		log.Println("ipc server created.")
		return nil
	}
}
func (s *IpcServer) Stop() {
	if s.state != Run {
		return
	} else {
		s.lock.Lock()
		if s.state != Run {
			return
		} else {
			s.setState(Disposed)
		}
		err := s.ln.Close()
		if err != nil {
			log.Println("close ipcserver ln err:", err)
		}
		if s.conn != nil {
			s.conn.Close()
			//这里会导致python进程退出
			log.Println("close ipcserver conn...")
		}
		s.OnMsgReveiveEvent.RemoveAllEvent()
		s.lock.Unlock()
		s.wg.Wait()
		log.Println("ipcserver has stoped... all go routine has exist")
	}
}
func (s *IpcServer) SendRequest(request IpcRequest) error {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return err
	}
	s.lock.Lock() //防止多用户并发发送消息时，head和body不能匹配的情况
	err = s.writeByte(IntToBytes(len(jsonData)))
	if err != nil {
		s.lock.Unlock()
		return err
	}
	err = s.writeByte(jsonData)
	if err != nil {
		s.lock.Unlock()
		return err
	}
	s.lock.Unlock()
	if request.IsExit {
		log.Println("come in stop ipcserver ")
		s.Stop()
	}
	return nil
}

/*私有方法*/
func (s *IpcServer) writeByte(buf []byte) error {
	if s.conn == nil {
		return errors.New("ipcserver conn is nil")
	}
	n, err := s.conn.Write(buf)
	if err != nil {
		return err
	}
	fmt.Println("write data res:", n)
	return nil
}
func (s *IpcServer) readByte(count int) ([]byte, error) {
	buf := make([]byte, count)
	_, err := s.conn.Read(buf)
	return buf, err
}
func (s *IpcServer) killListener(port int) {
	processManager := NewProcesManager()
	processManager.Kill(port)
}
func (s *IpcServer) createListen() (net.Listener, error) {
	port, err := s.config.IpcServerSetting.GetPort()
	if err == nil {
		s.killListener(port)
	}
	protocol, address := s.config.IpcServerSetting.GetConnectInfo()
	return net.Listen(protocol, address)
}

func (s *IpcServer) createConn() error {
	ln, err := s.createListen()
	if err != nil {
		return err
	}
	s.ln = ln
	s.setState(Run)
	go s.waitForConnect()
	return nil
}
func (s *IpcServer) waitForConnect() {
	defer log.Println("waitForConnect go routine exit.")
	defer s.wg.Done()
	s.wg.Add(1)
	for s.state == Run {
		// 等待客户端连接
		log.Println("waiting for connecting...")
		conn, err := s.ln.Accept()
		if err != nil {
			log.Println(err)
			if s.state != Disposed {
				continue
			} else {
				return
			}
		}
		s.conn = conn
		fmt.Println("客户端已连接:", conn.RemoteAddr())
		// 处理连接
		go s.handleConnection()
		//发送openai key
		request := OpenaiKeyIpcRequest(s.config.OpenaiSetting.GetOpenaiSetting())
		err = s.SendRequest(request)
		if err != nil {
			panic(err)
		}
	}
}

func (s *IpcServer) handleConnection() {
	defer log.Println("handleConnection go routine exit.")
	defer s.wg.Done()
	s.wg.Add(1)
	for s.state == Run {
		response := s.readResponse()
		//publish onMsgReveive event
		s.OnMsgReveiveEvent.Invoke(response)
	}
}

func (s *IpcServer) readResponse() IpcResponse {
	buf, err := s.readByte(4)
	if err != nil {
		return IpcResponseError(err, "Read data length from pipe Fatal")
	}
	contentLength := BytesToInt(buf)
	buf, err = s.readByte(int(contentLength))
	if err != nil {
		return IpcResponseError(err, "Read data from pipe Fatal")
	}
	var response IpcResponse
	if err := json.Unmarshal(buf, &response); err != nil {
		return IpcResponseError(err, "Unmarshal data err")
	}
	return response
}
func (s *IpcServer) setState(state IpcState) {
	s.state = state
}
