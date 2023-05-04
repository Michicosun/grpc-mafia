package main

import (
	"os"

	"grpc-mafia/chat"
	"grpc-mafia/logger"
	"grpc-mafia/util"

	zlog "github.com/rs/zerolog/log"
)

const (
	DEFAULT_COORD_CHAT_PORT = "20500"
)

func main() {
	logger.Init()

	port := util.GetEnvWithDefault("PORT", DEFAULT_COORD_CHAT_PORT)

	coord, err := chat.MakeCoordinator(port)
	if err != nil {
		zlog.Fatal().Err(err).Msg("create coordinator")
		os.Exit(1)
	}

	zlog.Info().Msg("chat coordinator started")

	coord.Listen()
}
