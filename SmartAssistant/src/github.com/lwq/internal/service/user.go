package service

import (
	"encoding/binary"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	. "github.com/lwq/internal/handle"
	. "github.com/lwq/internal/shared/consts"
	. "github.com/lwq/third_party/ipc"
	. "github.com/lwq/third_party/ipc/dto"
	. "github.com/lwq/utils/event"
)

type User struct {
	account          string
	addr             string
	wsConn           *websocket.Conn
	appModule        AppModule
	sendChan         chan []byte
	healthCheckChan  chan []byte
	UserOfflineEvent *Event
	requestIdList    []uuid.UUID
	engine           *PyEngine
}

func NewUser(conn *websocket.Conn, account string) *User {
	var user = &User{
		account:          account,
		addr:             conn.RemoteAddr().String(),
		wsConn:           conn,
		sendChan:         make(chan []byte),
		healthCheckChan:  make(chan []byte),
		requestIdList:    make([]uuid.UUID, 0),
		UserOfflineEvent: NewEvent(),
		engine:           NewEngine(),
	}
	user.engine.OnMsgReveiveEvent.AddEventHandler("OnMsgReveive", user.recvEngineMessage)
	return user
}

func (user *User) GetUserAddr() string {
	return user.addr
}
func (user *User) GetUserName() string {
	return user.account
}
func (user *User) GetUserConn() *websocket.Conn {
	return user.wsConn
}

// 接收消息
func (user *User) recvMessage() {
	defer user.offline()
	for {
		messageType, byteMsg, err := user.wsConn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Connection closed: %v\n", err)
			} else if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Conn has been closed", err)
			} else {
				log.Printf("Read error: %v\n", err)
			}
			return
		}
		switch messageType {
		case websocket.TextMessage:
			// 处理文本消息
			//chat with gpt
			log.Printf("Recv [%s] text msg: %s\n", user.account, string(byteMsg))
			if user.appModule == LiveChat {
				go user.chatWithGpt(user.account, string(byteMsg))
			} else {
				go user.sendMsgToApps(string(byteMsg))
			}
		case websocket.BinaryMessage:
			// 处理二进制消息
			option := AppModule(binary.LittleEndian.Uint32(byteMsg))
			user.appModule = option
			log.Printf("Recv [%s] byte msg,choose module: %v\n", user.account, option)
		default:
			log.Printf("Recv [%s] unknown message type: %d\n", user.account, messageType)
		}
	}
}
func (user *User) recvEngineMessage(data interface{}) {
	rep := data.(IpcResponse)
	log.Println("user reveMsg:", rep)
	if rep.Code == 500 {
		user.sendChan <- []byte("engine err:" + rep.ErrorMsg)
		return
	}
	for i, id := range user.requestIdList {
		if id == rep.ResponseId {
			user.sendChan <- []byte(rep.Message)
			user.requestIdList = append(user.requestIdList[:i], user.requestIdList[i+1:]...)
			user.sendChan <- []byte("end")
			break
		}
	}
}

func (user *User) chatWithGpt(userName string, message string) {
	err := HandleWsMessgae(user.account, user.sendChan, message)
	if err != nil {
		log.Println(err.Error())
	}
}

func (user *User) sendMsgToApps(message string) {
	request := IpcRequest{
		Id:      uuid.New(),
		Module:  user.appModule,
		Message: message,
	}
	user.requestIdList = append(user.requestIdList, request.Id)
	user.engine.SendRequest(request)
}

// 发送消息
func (user *User) sendMessage() {
	defer user.offline()
	for {
		select {
		case buf := <-user.sendChan:
			err := user.wsConn.WriteMessage(websocket.TextMessage, buf)
			if err != nil {
				log.Println("Send Msg Error：", err)
				return
			}
		case <-user.healthCheckChan:
			err := user.wsConn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(time.Second))
			if err != nil {
				log.Println("Send Pong Error：", err)
				return
			}
		}
	}
}
func (user *User) healthCheck(appdata string) error {
	user.healthCheckChan <- []byte("1")
	return nil
}

func (user *User) Online() {
	log.Printf("[%s] Online", user.account)
	user.wsConn.SetPingHandler(user.healthCheck)
	go user.recvMessage()
	go user.sendMessage()
}

func (user *User) offline() {
	log.Printf("[%s] OffLine", user.account)
	user.wsConn.Close()
	//publish offline event
	user.UserOfflineEvent.Invoke(user)
}
