package main

import (
	"grpc-mafia/registry"
	"grpc-mafia/util"
	"strconv"
)

const (
	DEFAULT_DATA_FOLDER = "/tmp/grpc-mafia"
	DEFAULT_RENDERS_CNT = "10"
)

func main() {
	data_folder := util.GetEnvWithDefault("DATA_FOLDER", DEFAULT_DATA_FOLDER)
	renders_str := util.GetEnvWithDefault("RENDERS_CNT", DEFAULT_RENDERS_CNT)

	renders, err := strconv.Atoi(renders_str)
	if err != nil {
		panic(err)
	}

	registry.Server.Init(data_folder, renders)
	registry.Server.Run()
}
