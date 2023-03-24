package user

import (
	"log"

	"github.com/gorilla/websocket"
	. "github.lwq.com/Global"
	. "github.lwq.com/Utils/Event"
)

type User struct {
	addr             string
	wsConn           *websocket.Conn
	sendChan         chan []byte
	UserOfflineEvent *Event
}

func NewUser(conn *websocket.Conn) *User {
	var user = &User{
		addr:             conn.RemoteAddr().String(),
		wsConn:           conn,
		sendChan:         make(chan []byte),
		UserOfflineEvent: NewEvent(),
	}
	return user
}
func (user *User) GetUserAddr() string {
	return user.addr
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

//接收消息
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
		msg := CtClient.CreateChatCompletion(Ctx, string(byteMsg))
		user.sendChan <- []byte(msg)
	}
}

//发送消息
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
