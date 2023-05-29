package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"grpc-mafia/logger"
	"grpc-mafia/server"
	"grpc-mafia/util"

	"google.golang.org/grpc"

	mafia "grpc-mafia/server/proto"

	zlog "github.com/rs/zerolog/log"
)

const (
	DEFAULT_GRPC_PORT = "9000"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	logger.Init()

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
