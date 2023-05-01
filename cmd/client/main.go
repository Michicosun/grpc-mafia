package main

import (
	"flag"
	"fmt"
	"grpc-mafia/client"
	"os"
)

func main() {
	var grpc_host string
	var grpc_port string

	flag.StringVar(&grpc_host, "grpc_host", "localhost", "host of grpc game service")
	flag.StringVar(&grpc_port, "grpc_port", "9000", "port of grpc game service")

	flag.Parse()

	client.Parser.Init()
	client.Printer.Init()
	client.Game.Init()

	if err := client.GrpcConnect.Init(grpc_host, grpc_port); err != nil {
		fmt.Printf("ERROR: %e\n", err)
		os.Exit(1)
	}

	client.Parser.Run()
}
