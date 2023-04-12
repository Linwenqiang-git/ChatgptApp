package service

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

type Server struct {
	onlineUserMap map[string]*User
	userMapLock   sync.RWMutex
}

// 将请求升级为 web socket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	account, _, _ := s.parseBasicAuth(auth)
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error：", err)
		return
	}
	var user = NewUser(wsConn, account)
	user.Online()
	s.addOnlineUserMap(user)
	//注册User下线事件
	user.UserOfflineEvent.AddEventHandler("UserOffline", s.deleteOnlineUserMap)
}

func (s *Server) parseBasicAuth(auth string) (string, string, error) {
	if !strings.HasPrefix(auth, "Basic ") {
		return "", "", fmt.Errorf("Invalid authorization header")
	}
	payload, err := base64.StdEncoding.DecodeString(auth[6:])
	if err != nil {
		return "", "", fmt.Errorf("Invalid authorization header")
	}
	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		return "", "", fmt.Errorf("Invalid authorization header")
	}
	return pair[0], pair[1], nil
}

func (s *Server) addOnlineUserMap(user *User) {
	s.userMapLock.Lock()
	s.onlineUserMap[user.GetUserAddr()] = user
	s.userMapLock.Unlock()
}

func (s *Server) deleteOnlineUserMap(data interface{}) {
	user := data.(*User)
	s.userMapLock.Lock()
	delete(s.onlineUserMap, user.GetUserAddr())
	s.userMapLock.Unlock()
}
func (s *Server) Start() {
	fmt.Println("Service is Listening:8899")
	http.HandleFunc("/openai", s.handler)
	http.ListenAndServe(":8899", nil)
}
