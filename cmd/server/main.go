package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"grpc-mafia/logger"
	"grpc-mafia/server"

	"google.golang.org/grpc"

	mafia "grpc-mafia/server/proto"

	zlog "github.com/rs/zerolog/log"
)

func main() {
	logger.Init()

	port := os.Getenv("PORT")
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
