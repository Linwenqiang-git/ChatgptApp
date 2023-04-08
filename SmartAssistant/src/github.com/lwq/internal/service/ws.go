package service

func CreateWebSocketServer() *Server {
	return &Server{
		onlineUserMap: make(map[string]*User),
	}
}
