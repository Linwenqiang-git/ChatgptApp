package main

import (
	. "github.com/lwq/internal/server"
)

func main() {
	CreateWebSocketServer().Start()
	select {}
}
