package main

import (
	"fmt"
	client "grpc-mafia/client"
	chat "grpc-mafia/client/chat"
	game "grpc-mafia/client/game"
	grpc "grpc-mafia/client/grpc"
	"grpc-mafia/client/tracker_client"
	util "grpc-mafia/util"
	"math/rand"
	"os"
	"time"
)

const (
	DEFAULT_RABBIT_HOST = "localhost"
	DEFAULT_RABBIT_PORT = "5672"
	DEFAULT_RABBIT_USER = "guest"
	DEFAULT_RABBIT_PASS = "guest"

	DEFAULT_SERVER_HOST = "localhost"
	DEFAULT_SERVER_PORT = "9000"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	rabbit_creds := chat.RabbitCredentials{
		Host: util.GetEnvWithDefault("RABBIT_HOST", DEFAULT_RABBIT_HOST),
		Port: util.GetEnvWithDefault("RABBIT_PORT", DEFAULT_RABBIT_PORT),
		User: util.GetEnvWithDefault("RABBIT_USER", DEFAULT_RABBIT_USER),
		Pass: util.GetEnvWithDefault("RABBIT_PASS", DEFAULT_RABBIT_PASS),
	}

	chat.RabbitConnection.Init(rabbit_creds)

	grpc_host := util.GetEnvWithDefault("SRV_HOST", DEFAULT_SERVER_HOST)
	grpc_port := util.GetEnvWithDefault("SRV_PORT", DEFAULT_SERVER_PORT)

	if err := grpc.Connection.Init(grpc_host, grpc_port); err != nil {
		fmt.Printf("ERROR: %e\n", err)
		os.Exit(1)
	}

	game.Session.Init()
	tracker_client.TrackerClient.Init()

	_, use_bot := os.LookupEnv("USE_BOT")
	client := client.MakeClient(use_bot)

	client.Run()
}

// USE_BOT=yes SRV_HOST=localhost SRV_PORT=9000 go run client/main.go
