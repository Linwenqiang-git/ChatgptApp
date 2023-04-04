package main

import (
	"bufio"
	"flag"
	"log"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8899", "http service address")
var clientDone = make(chan struct{})
var serverDone = make(chan struct{})
var lastHealthTime time.Time

func main() {
	flag.Parse()
	log.SetFlags(0)

	serviceUrl := url.URL{Scheme: "ws", Host: *addr, Path: "/openai"}
	conn, _, err := websocket.DefaultDialer.Dial(serviceUrl.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
		return
	}
	conn.SetPongHandler(func(appdata string) error {
		lastHealthTime = time.Now()
		//log.Printf("recv pong")
		return nil
	})
	defer func() {
		log.Printf("exit from %s", serviceUrl.String())
		conn.Close()
	}()
	log.Printf("connecting to %s", serviceUrl.String())
	// 创建一个 WaitGroup 对象，用于协调 goroutine 同步
	lastHealthTime = time.Now()
	var wg sync.WaitGroup
	wg.Add(2)
	//接收消息
	go receiveMsg(conn, &wg)
	//心跳检测
	go healthCheck(conn, &wg)
	//读取用户输入
	readInput(conn)
	// 等待所有 goroutine 结束
	wg.Wait()
	releaseAllChannel()
}

func releaseAllChannel() {
	_, ok := <-clientDone
	if ok {
		close(clientDone)
	}
	_, ok = <-serverDone
	if ok {
		close(serverDone)
	}
}

// 接收消息
func receiveMsg(conn *websocket.Conn, wg *sync.WaitGroup) {
	defer close(serverDone)
	defer wg.Done()
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Conn has been closed", err)
				// 等待一段时间后重试
				time.Sleep(time.Second * 3)
				continue
			}
			log.Println("chatgpt error:", err)
			return
		}
		switch messageType {
		case websocket.TextMessage:
			// 处理文本消息
			log.Printf("\n 【chatgpt】: %s \n", message)
		case websocket.BinaryMessage:
			// 处理二进制消息
			log.Printf("Recv binary msg: %v\n", message)
		default:
			log.Printf("Recv unknown message type: %d\n", messageType)
		}
	}
}

func healthCheck(conn *websocket.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	for {
		select {
		// 定期发送心跳包
		case <-ticker.C:
			err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second))
			if err != nil {
				log.Println("health check error:", err)
				return
			}
			// 检查心跳超时
			if time.Since(lastHealthTime) > time.Second*10 {
				log.Println("health check timeout")
				return
			}
		// 接收关闭信号
		case <-clientDone: //chan被关闭时返回
			println(" client has Done")
			return
		case <-serverDone:
			println(" server has Done")
			return
		}
	}
}

func readInput(conn *websocket.Conn) {
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {

		line := input.Text()
		// 输入bye时 结束
		if line == "bye" {
			close(clientDone)
			break
		}
		log.Printf("\n 【you】: %s", line)
		err := conn.WriteMessage(websocket.TextMessage, []byte(line))
		if err != nil {
			log.Println("send to server err:", err)
			continue
		}
	}
}
