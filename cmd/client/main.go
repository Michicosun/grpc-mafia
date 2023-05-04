package main

import (
	"fmt"
	client "grpc-mafia/client"
	game "grpc-mafia/client/game"
	grpc "grpc-mafia/client/grpc"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	grpc_host := os.Getenv("SRV_HOST")
	grpc_port := os.Getenv("SRV_PORT")
	_, use_bot := os.LookupEnv("USE_BOT")

	if err := grpc.Connection.Init(grpc_host, grpc_port); err != nil {
		fmt.Printf("ERROR: %e\n", err)
		os.Exit(1)
	}

	game.Session.Init()
	client := client.MakeClient(use_bot)

	client.Run()
}

// USE_BOT=yes SRV_HOST=localhost SRV_PORT=9000 go run client/main.go
