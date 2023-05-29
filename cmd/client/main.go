package main

import (
	"fmt"
	client "grpc-mafia/client"
	game "grpc-mafia/client/game"
	grpc "grpc-mafia/client/grpc"
	"grpc-mafia/util"
	"math/rand"
	"os"
	"time"
)

const (
	DEFAULT_COORD_CHAT_HOST = "localhost"
	DEFAULT_COORD_CHAT_PORT = "20500"

	DEFAULT_SERVER_HOST = "localhost"
	DEFAULT_SERVER_PORT = "9000"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	grpc_host := util.GetEnvWithDefault("SRV_HOST", DEFAULT_SERVER_HOST)
	grpc_port := util.GetEnvWithDefault("SRV_PORT", DEFAULT_SERVER_PORT)

	if err := grpc.Connection.Init(grpc_host, grpc_port); err != nil {
		fmt.Printf("ERROR: %e\n", err)
		os.Exit(1)
	}

	game.Session.Init()

	_, use_bot := os.LookupEnv("USE_BOT")
	client := client.MakeClient(use_bot)

	client.Run()
}

// USE_BOT=yes SRV_HOST=localhost SRV_PORT=9000 go run client/main.go
