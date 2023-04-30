package main

import (
	"fmt"
	"grpc-mafia/chat"
)

func main() {
	chat.Connector.Init(2050, "127.0.0.1", 20050)

	go func() {
		for i := 0; i < 10000; i += 1 {
			chat.Connector.MakeBCast("new_group", chat.Message{
				From: "michicosun",
				Text: "hello!",
			})
		}
	}()

	for i := 0; i < 10000; i += 1 {
		msg1, _ := chat.Connector.RecvMessage()
		chat.Connector.RecvMessage()
		fmt.Println(msg1.From, msg1.Text)
	}

	// chat.Connector.CreateGroup("new_group", []string{"127.0.0.1:2050", "localhost:2050"})
}
