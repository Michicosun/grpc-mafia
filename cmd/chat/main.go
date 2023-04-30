package main

import (
	"os"

	"grpc-mafia/chat"
	"grpc-mafia/logger"

	zlog "github.com/rs/zerolog/log"
)

func main() {
	logger.Init()

	coord, err := chat.MakeCoordinator()
	if err != nil {
		zlog.Fatal().Err(err).Msg("create coordinator")
		os.Exit(1)
	}

	zlog.Info().Msg("chat coordinator started")

	coord.Listen()
}
