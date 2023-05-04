package main

import (
	"fmt"
	"grpc-mafia/chat"
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

func InitChatConnector() {
	coord_host := util.GetEnvWithDefault("CHAT_COORD_HOST", DEFAULT_COORD_CHAT_HOST)
	coord_port := util.GetEnvWithDefault("CHAT_COORD_PORT", DEFAULT_COORD_CHAT_PORT)

	if err := chat.Connector.Init(coord_host, coord_port); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	InitChatConnector()

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
