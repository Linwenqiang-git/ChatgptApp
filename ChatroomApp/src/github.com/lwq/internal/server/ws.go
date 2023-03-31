package server

import (
	. "github.com/lwq/internal/models"
)

func CreateWebSocketServer() *Server {
	return &Server{
		onlineUserMap: make(map[string]*User),
	}
}
