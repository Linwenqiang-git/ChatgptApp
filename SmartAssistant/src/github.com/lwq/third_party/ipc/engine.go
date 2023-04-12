package ipc

import (
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"log"

	. "github.com/lwq/third_party/ipc/dto"
	. "github.com/lwq/utils/event"
)

var client *PyEngine
var once sync.Once

/*
python引擎，负责启动确保python进程正常运行
*/
type EngineStatus int

const (
	Init EngineStatus = iota
	Running
	Stop
)

type PyEngine struct {
	pyInterpreterPath string
	cmd               *exec.Cmd
	lock              sync.RWMutex
	engineStatus      EngineStatus
	restartEngineChan chan int
	exitMonitorChan   chan struct{}
	ipcServer         *IpcServer
	OnMsgReveiveEvent *Event
	wg                sync.WaitGroup
}

func NewEngine() *PyEngine {
	once.Do(func() {
		interpreterPath := "././Apps/app_interpreter.py"
		absPyInterpreterPath, err := filepath.Abs(interpreterPath)
		if err != nil {
			panic(err)
		}
		absPyInterpreterPath = strings.Replace(absPyInterpreterPath, "SmartAssistant\\src\\github.com\\lwq\\cmd\\", "", -1)
		client = &PyEngine{
			pyInterpreterPath: absPyInterpreterPath,
			engineStatus:      Init,
			restartEngineChan: make(chan int),
			exitMonitorChan:   make(chan struct{}),
			ipcServer:         NewIpcServer(),
			OnMsgReveiveEvent: NewEvent(),
		}
		go client.startEngine()
	})
	return client
}

func (p *PyEngine) SendRequest(request IpcRequest) error {
	return p.ipcServer.SendRequest(request)
}

func (p *PyEngine) Stop() {
	p.setStatus(Stop)
	close(p.exitMonitorChan)
	//关闭ipc_server相关连接
	//向子进程发送退出消息
	request := ExitIpcRequest()
	err := p.SendRequest(request)
	if err != nil {
		log.Println("engine stop error1:", err.Error())
		panic(err)
	}
	p.wg.Wait()
	log.Println("all go routine has exit")
}

/*私有方法*/
func (p *PyEngine) startEngine() {
	defer log.Println("startEngine go routine exit")
	defer p.wg.Done()
	p.wg.Add(1)
	err := p.ipcServer.Start()
	if err != nil {
		panic(err)
	}
	p.ipcServer.OnMsgReveiveEvent.AddEventHandler("OnMsgReveive", p.recvMessage)
	cmd := exec.Command("python", p.pyInterpreterPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // 设置子进程组ID，以便终止整个进程组
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	p.cmd = cmd
	p.setStatus(Running)
	go p.monitorEngine()
	log.Println("python engine created.")

	err = cmd.Wait()
	if err != nil {
		if p.engineStatus == Running {
			//restart engine
			log.Println("engine abnormal exit:", err, " restarting...")
			p.restartEngineChan <- 1
		}
	}
	if p.engineStatus == Stop {
		//exit
		log.Println("engine has exit[3]...")
		close(p.restartEngineChan)
	}
}
func (p *PyEngine) monitorEngine() {
	defer log.Println("monitorEngine go routine end...")
	defer p.wg.Done()
	p.wg.Add(1)
	for {
		select {
		case <-p.exitMonitorChan:
			return
		case isRestart := <-p.restartEngineChan:
			//检测python进程是否正常运行
			if isRestart == 1 {
				go p.startEngine()
			}
			return
		}
	}
}
func (p *PyEngine) recvMessage(data interface{}) {
	response := data.(IpcResponse)
	//publish event
	p.OnMsgReveiveEvent.Invoke(response)
}
func (p *PyEngine) setStatus(engineStatus EngineStatus) {
	p.lock.Lock()
	p.engineStatus = engineStatus
	p.lock.Unlock()
}
