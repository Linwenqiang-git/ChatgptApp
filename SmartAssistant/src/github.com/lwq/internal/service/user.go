package service

import (
	"encoding/binary"
	"log"
	"sync"
	"time"
	"unsafe"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	. "github.com/lwq/internal/handle"
	. "github.com/lwq/internal/shared/consts"
	. "github.com/lwq/third_party/ipc"
	. "github.com/lwq/third_party/ipc/dto"
	. "github.com/lwq/utils/event"
)

var onMsgReveiveEventName string = "OnMsgReveive"

type User struct {
	account          string
	addr             string
	wsConn           *websocket.Conn
	appModule        AppModule
	sendChan         chan []byte
	offLineChan      chan struct{}
	UserOfflineEvent *Event
	requestIdList    []uuid.UUID
	wg               sync.WaitGroup
	engine           *PyEngine
}

func NewUser(conn *websocket.Conn, account string) *User {
	var user = &User{
		account:          account,
		addr:             conn.RemoteAddr().String(),
		wsConn:           conn,
		sendChan:         make(chan []byte),
		offLineChan:      make(chan struct{}),
		requestIdList:    make([]uuid.UUID, 0),
		UserOfflineEvent: NewEvent(),
		engine:           NewEngine(),
	}
	user.engine.OnMsgReveiveEvent.AddEventHandler(onMsgReveiveEventName, user.onRecvEngineMessageHandle)
	return user
}

// 公共方法
func (user *User) GetUserAddr() string {
	return user.addr
}
func (user *User) GetUserName() string {
	return user.account
}
func (user *User) GetUserConn() *websocket.Conn {
	return user.wsConn
}
func (user *User) Online() {
	log.Printf("[%s] Online", user.account)
	user.wsConn.SetPingHandler(user.healthCheck)
	go user.recvMessage()
	go user.sendMessage()
	user.wg.Add(2)
}

// 接收消息
func (user *User) recvMessage() {
	defer user.offline()
	defer user.wg.Done()
	for {
		messageType, byteMsg, err := user.wsConn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) || websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Connection closed: %v\n", err)
			} else {
				log.Printf("Read error: %v\n", err)
			}
			return
		}
		switch messageType {
		case websocket.TextMessage:
			// 处理文本消息
			log.Printf("Recv [%s] text msg: %s\n", user.account, string(byteMsg))
			if user.appModule == LiveChat {
				//chat with gpt
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

func (user *User) chatWithGpt(userName string, message string) {
	user.wg.Add(1)
	defer user.wg.Done()
	//读取消息异常以后也不会写入
	err := HandleWsMessgae(user.account, user.sendChan, message)
	if err != nil {
		log.Println(err.Error())
	}
}

func (user *User) sendMsgToApps(message string) {
	user.wg.Add(1)
	defer user.wg.Done()
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
	defer user.wg.Done()
	for {
		select {
		case buf, ok := <-user.sendChan:
			if !ok {
				return
			}
			err := user.wsConn.WriteMessage(websocket.TextMessage, buf)
			if err != nil {
				log.Println("Send Msg Error：", err)
			}
		case <-user.offLineChan:
			close(user.sendChan)
			return
		}
	}
}
func (user *User) healthCheck(appdata string) error {
	err := user.wsConn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(time.Second))
	if err != nil {
		log.Printf("Send %s Pong Error：%e", user.account, err)
		return err
	}
	return nil
}

func (user *User) offline() {
	log.Printf("[%s] OffLine", user.account)
	//remove handle
	f := user.onRecvEngineMessageHandle
	user.engine.OnMsgReveiveEvent.RemoveEventHandler(onMsgReveiveEventName, uintptr(unsafe.Pointer(&f)))
	//publish offline event
	user.UserOfflineEvent.Invoke(user)
	user.wsConn.Close()
	close(user.offLineChan)
	user.wg.Wait()
}

// 事件移除以后，讲不会写入消息
func (user *User) onRecvEngineMessageHandle(data interface{}) {
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
