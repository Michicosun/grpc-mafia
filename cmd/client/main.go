package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	mafia "grpc-mafia/server/proto"

	zlog "github.com/rs/zerolog/log"
)

func main() {
	// client.GameState.Init()
	// client.Parser.Init()
	// client.Printer.Init()

	// go func() {
	// 	for {
	// 		time.Sleep(1 * time.Second)
	// 		client.Printer.PrintLine("Hello")
	// 	}
	// }()

	// client.Parser.Run()

	conn, err := grpc.Dial(":9000", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		zlog.Fatal().Err(err).Msg("grpc dial")
	}
	defer conn.Close()

	stub := mafia.NewMafiaServiceClient(conn)

	client, err := stub.FindGame(context.Background())

	client.Send(&mafia.Action{
		Type: mafia.ActionType_Vote,
		Data: &mafia.Action_Vote_{
			Vote: &mafia.Action_Vote{
				Name: "michicosun",
			},
		},
	})

	if err != nil {
		zlog.Fatal().Err(err).Msg("find game")
		log.Fatalln(err)
	}

	response, _ := client.Recv()
	fmt.Println(response.GetMessage().Text)
}
