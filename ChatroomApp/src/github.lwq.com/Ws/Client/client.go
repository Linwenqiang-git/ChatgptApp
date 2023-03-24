package main

import (
	"bufio"
	"flag"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8899", "http service address")
var clientDone = make(chan struct{})
var serverDone = make(chan struct{})

func main() {
	flag.Parse()
	log.SetFlags(0)

	serviceUrl := url.URL{Scheme: "ws", Host: *addr, Path: "/openai"}
	conn, _, err := websocket.DefaultDialer.Dial(serviceUrl.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	log.Printf("connecting to %s", serviceUrl.String())
	defer func() {
		conn.Close()
		log.Printf("exit from %s", serviceUrl.String())
	}()

	//接收消息
	go receiveMsg(conn, serverDone)
	readInput(conn, clientDone)
	for {
		select {
		case <-clientDone: //chan被关闭时返回
			return
		case <-serverDone:
			return
		}
	}
}
func receiveMsg(conn *websocket.Conn, serverDone chan struct{}) {
	defer close(serverDone)
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("chatgpt error:", err)
			return
		}
		log.Printf("\n chatgpt: %s", message)
	}
}
func readInput(conn *websocket.Conn, clientDone chan struct{}) {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		// 输入bye时 结束
		if line == "bye" {
			close(clientDone)
			break
		}
		err := conn.WriteMessage(websocket.TextMessage, []byte(line))
		if err != nil {
			log.Println("send to server err:", err)
			continue
		}
	}
}
