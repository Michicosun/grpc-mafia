package main

import (
	"grpc-mafia/registry"
	"grpc-mafia/registry/queue"
	"grpc-mafia/util"
	"strconv"
)

const (
	DEFAULT_RABBIT_HOST = "localhost"
	DEFAULT_RABBIT_PORT = "5672"
	DEFAULT_RABBIT_USER = "guest"
	DEFAULT_RABBIT_PASS = "guest"

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

	rabbit_creds := queue.RabbitCredentials{
		Host: util.GetEnvWithDefault("RABBIT_HOST", DEFAULT_RABBIT_HOST),
		Port: util.GetEnvWithDefault("RABBIT_PORT", DEFAULT_RABBIT_PORT),
		User: util.GetEnvWithDefault("RABBIT_USER", DEFAULT_RABBIT_USER),
		Pass: util.GetEnvWithDefault("RABBIT_PASS", DEFAULT_RABBIT_PASS),
	}

	server_cfg := registry.ServerConfig{
		DataFolder:  data_folder,
		Renders:     renders,
		RabbitCreds: rabbit_creds,
	}

	registry.Server.Init(server_cfg)
	registry.Server.Run()
}
