package main

import (
	. "github.com/lwq/internal/service"
)

func main() {
	CreateWebSocketServer().Start()
	select {}
}
