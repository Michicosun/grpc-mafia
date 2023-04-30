package main

import (
	"grpc-mafia/client"
	"time"
)

func main() {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			client.Printer.PrintLine("Hello")
		}
	}()

	client.Parser.Run()
}
