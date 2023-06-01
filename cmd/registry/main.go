package main

import (
	"grpc-mafia/registry"
	"grpc-mafia/util"
)

const (
	DEFAULT_DATA_FOLDER = "/tmp/grpc-mafia"
)

func main() {
	data_folder := util.GetEnvWithDefault("DATA_FOLDER", DEFAULT_DATA_FOLDER)

	registry.Server.Init(data_folder)
	registry.Server.Run()
}
