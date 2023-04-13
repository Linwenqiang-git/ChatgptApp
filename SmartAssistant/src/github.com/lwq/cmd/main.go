package main

import (
	"log"
	"os"
	"os/signal"

	. "github.com/lwq/internal/service"
	. "github.com/lwq/third_party/ipc"
)

func main() {
	//程序启动时，就运行一个python引擎，单例模式确保全局一个引擎
	engine := NewEngine()

	exitChan := make(chan struct{})
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		<-signalChan // 接收到 SIGINT 信号时执行清理和退出操作
		log.Println("reveive stop signal")
		engine.Stop()
		close(exitChan)
	}()

	CreateWebSocketServer().Start()

	<-exitChan
}
