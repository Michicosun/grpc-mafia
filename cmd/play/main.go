package main

import (
	"grpc-mafia/chat"
)

func main() {
	chat.Connector.Init(2050, "127.0.0.1", 20050)
	chat.Connector.CreateGroup("new_group", []string{"127.0.0.1:2050", "localhost:2050"})
}
