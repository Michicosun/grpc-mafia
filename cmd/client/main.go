package main

import (
	"flag"
	"fmt"
	client "grpc-mafia/client"
	game "grpc-mafia/client/game"
	grpc "grpc-mafia/client/grpc"
	"os"
)

func main() {
	var grpc_host string
	var grpc_port string

	flag.StringVar(&grpc_host, "grpc_host", "localhost", "host of grpc game service")
	flag.StringVar(&grpc_port, "grpc_port", "9000", "port of grpc game service")

	flag.Parse()

	if err := grpc.Connection.Init(grpc_host, grpc_port); err != nil {
		fmt.Printf("ERROR: %e\n", err)
		os.Exit(1)
	}

	game.Session.Init()
	client := client.MakeClient(false)

	client.Run()
}
