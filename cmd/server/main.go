package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"grpc-mafia/chat"
	"grpc-mafia/logger"
	"grpc-mafia/server"
	"grpc-mafia/util"

	"google.golang.org/grpc"

	mafia "grpc-mafia/server/proto"

	zlog "github.com/rs/zerolog/log"
)

const (
	DEFAULT_COORD_CHAT_HOST = "localhost"
	DEFAULT_COORD_CHAT_PORT = "20500"

	DEFAULT_GRPC_PORT = "9000"
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

	logger.Init()
	InitChatConnector()

	port := util.GetEnvWithDefault("PORT", DEFAULT_GRPC_PORT)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		zlog.Fatal().Err(err).Msg("failed to listen")
		log.Fatalf(": %v", err)
	}

	srv := grpc.NewServer()
	mafia.RegisterMafiaServiceServer(srv, server.MakeGameServer())

	zlog.Info().Str("port", port).Msg("started grpc server")

	log.Fatalln(srv.Serve(lis))
}
