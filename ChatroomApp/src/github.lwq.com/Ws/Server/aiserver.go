package server

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	. "github.lwq.com/Ws/User"

	"github.com/gorilla/websocket"
)

type Server struct {
	onlineUserMap map[string]*User
	userMapLock   sync.RWMutex
}

func CreateServer() *Server {
	return &Server{
		onlineUserMap: make(map[string]*User),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (server *Server) handler(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	server.createNewUser(wsConn)
}

func (server *Server) createNewUser(conn *websocket.Conn) {
	var user = NewUser(conn)
	user.Online()
	server.addOnlineUserMap(user)
	//注册User下线事件
	user.UserOfflineEvent.AddEventHandler("UserOffline", server.deleteOnlineUserMap)

}

func (server *Server) addOnlineUserMap(user *User) {
	server.userMapLock.Lock()
	server.onlineUserMap[user.GetUserAddr()] = user
	server.userMapLock.Unlock()
}

func (server *Server) deleteOnlineUserMap(data interface{}) {
	user := data.(*User)
	server.userMapLock.Lock()
	delete(server.onlineUserMap, user.GetUserAddr())
	server.userMapLock.Unlock()
}

func (server *Server) Start() {
	fmt.Println("Service is Listening:8899")
	http.HandleFunc("/openai", server.handler)
	http.ListenAndServe(":8899", nil)
}
