package main

import (
	"fmt"
	"grpc-mafia/chat"
)

func main() {
	chat.Connector.Init(2050, "127.0.0.1", 20050)

	go func() {
		for i := 0; i < 10000; i += 1 {
			chat.Connector.MakeBCast("new_group", "hello world!")
		}
	}()

	for i := 0; i < 10000; i += 1 {
		msg1, _ := chat.Connector.RecvMessage()
		msg2, _ := chat.Connector.RecvMessage()
		fmt.Println(i, msg1, msg2)
	}

	// chat.Connector.CreateGroup("new_group", []string{"127.0.0.1:2050", "localhost:2050"})
}
