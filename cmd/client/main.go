package main

import (
	"grpc-mafia/client"
	"time"
)

func main() {
	client.GameState.Init()
	client.Parser.Init()
	client.Printer.Init()

	go func() {
		for {
			time.Sleep(1 * time.Second)
			client.Printer.PrintLine("Hello")
		}
	}()

	client.Parser.Run()
}
