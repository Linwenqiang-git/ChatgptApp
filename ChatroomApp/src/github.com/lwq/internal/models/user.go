package models

import (
	"log"

	"github.com/gorilla/websocket"
	. "github.com/lwq/internal/handle"
	. "github.com/lwq/pkg/utils/event"
)

type User struct {
	account          string
	addr             string
	wsConn           *websocket.Conn
	sendChan         chan []byte
	UserOfflineEvent *Event
}

func CreatUser(conn *websocket.Conn) *User {
	var user = &User{
		addr:             conn.RemoteAddr().String(),
		wsConn:           conn,
		sendChan:         make(chan []byte),
		UserOfflineEvent: NewEvent(),
	}
	return user
}

// 账号登录
func Login(userName string) {

}
func (user *User) GetUserAddr() string {
	return user.addr
}
func (user *User) GetUserConn() *websocket.Conn {
	return user.wsConn
}

// 接收消息
func (user *User) recvMessage() {
	defer user.offline()
	for {
		_, byteMsg, err := user.wsConn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Recv [%s] msg:%s", user.addr, byteMsg)
		//chat with gpt
		msg, err := HandleWsMessgae(user.account, string(byteMsg))
		if err != nil {
			msg = err.Error()
		}
		user.sendChan <- []byte(msg)
	}
}

// 发送消息
func (user *User) sendMessage() {
	defer user.offline()
	for {
		buf := <-user.sendChan
		err := user.wsConn.WriteMessage(1, buf)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Send [%s] msg:%s", user.addr, buf)
	}
}

func (user *User) Online() {
	log.Printf("[%s] Online", user.addr)
	go user.recvMessage()
	go user.sendMessage()
}

func (user *User) offline() {
	user.wsConn.Close()
	log.Printf("[%s] OffLine", user.addr)
	//publish offline event
	user.UserOfflineEvent.Invoke(user)
}
