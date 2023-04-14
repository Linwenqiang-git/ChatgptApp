package main

import (
	"bufio"
	"encoding/base64"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	. "github.com/lwq/internal/shared/consts"
)

var addr = flag.String("addr", "localhost:8899", "http service address")
var clientDone = make(chan struct{})
var serverDone = make(chan struct{})
var lastHealthTime time.Time

func main() {
	flag.Parse()
	log.SetFlags(0)
	conn, err := login()
	if err != nil {
		log.Print("login error", err)
		return
	}
	defer func() {
		log.Printf("exit from localhost:8899,bye...")
		conn.Close()
	}()
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

func login() (*websocket.Conn, error) {
	log.Print("Please your account:")
	input := bufio.NewScanner(os.Stdin)
	username := ""
	for input.Scan() {
		username = input.Text()
		if username != "" {
			break
		}
		log.Print("Account is not allowed to be empty")
	}
	password := ""
	auth := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password))))
	headers := http.Header{}
	headers.Set("Authorization", auth)
	serviceUrl := url.URL{Scheme: "ws", Host: *addr, Path: "/openai"}
	conn, _, err := websocket.DefaultDialer.Dial(serviceUrl.String(), headers)
	if err != nil {
		log.Fatal("dial:", err)
		return nil, err
	}
	conn.SetPongHandler(func(appdata string) error {
		lastHealthTime = time.Now()
		return nil
	})
	log.Println("connect to server:local:8899")
	welcome(conn)
	return conn, nil
}
func welcome(conn *websocket.Conn) AppModule {
	fmt.Println(`┌───────────────────────────────┐`)
	fmt.Println(`│                               │`)
	fmt.Println(`│      欢迎使用大G助手！       │`)
	fmt.Println(`│                               │`)
	fmt.Println(`│       请选择交互模式：        │`)
	fmt.Println(`│                               │`)
	fmt.Println(`│        1. 帮助中心           │`)
	fmt.Println(`│        2. 需求整理           │`)
	fmt.Println(`│        3. 在线聊天           │`)
	fmt.Println(`│                               │`)
	fmt.Println(`└───────────────────────────────┘`)

	optionList := GetModuleOption()
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		option, err := strconv.Atoi(input.Text())
		if err != nil {
			log.Print("Unrecognized pattern, please try again")
			continue
		}
		for _, validOption := range optionList {
			if AppModule(option) == validOption {
				optionBytes := make([]byte, 4)
				binary.LittleEndian.PutUint32(optionBytes, uint32(option))
				conn.WriteMessage(websocket.BinaryMessage, optionBytes)
				return validOption
			}
		}
		log.Print("Unrecognized pattern, please try again")
	}
	return LiveChat
}

func releaseAllChannel() {
	select {
	case <-clientDone:
	default:
		log.Print("release clientDone")
		close(clientDone)
	}
	select {
	case <-serverDone:
	default:
		log.Print("release serverDone")
		close(clientDone)
	}
}

// 接收消息
func receiveMsg(conn *websocket.Conn, wg *sync.WaitGroup) {
	defer close(serverDone)
	defer wg.Done()
	retryCount := 0
	continueMsg := false
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) && retryCount < 3 {
				log.Println("Conn has been closed", err)
				// 等待一段时间后重试
				time.Sleep(time.Second * 3)
				retryCount++
				continue
			}
			log.Println("chatgpt error:", err)
			return
		}
		retryCount = 0
		switch messageType {
		case websocket.TextMessage:
			if !continueMsg {
				log.Printf("\n 【chatgpt】: ")
				continueMsg = true
			}
			if string(message) == "end" {
				continueMsg = false
				log.Print("\n")
				continue
			}
			// 处理文本消息
			fmt.Fprintf(os.Stdout, "%s", message)
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
		select {
		case <-clientDone:
			return
		case <-serverDone:
			return
		default:
			// 继续执行
		}
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
