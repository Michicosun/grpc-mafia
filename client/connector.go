package client

import (
	"context"
	"fmt"
	mafia "grpc-mafia/server/proto"
	"io"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var GrpcConnect = &grpcConnect{}

type grpcConnect struct {
	conn   *grpc.ClientConn
	client mafia.MafiaServiceClient
	stream mafia.MafiaService_FindGameClient
}

func (gc *grpcConnect) SendInit(name string) error {
	fmt.Println("sending Init Game action")

	return gc.stream.Send(&mafia.Action{
		Type: mafia.ActionType_Init,
		Data: &mafia.Action_Init_{
			Init: &mafia.Action_Init{
				Name: name,
			},
		},
	})
}

func (gc *grpcConnect) SendVote(from string, to string) error {
	fmt.Println("sending Vote action")

	return gc.stream.Send(&mafia.Action{
		Type: mafia.ActionType_Vote,
		Data: &mafia.Action_Vote_{
			Vote: &mafia.Action_Vote{
				From: from,
				Name: to,
			},
		},
	})
}

func (gc *grpcConnect) SendDoNothing() error {
	fmt.Println("sending DoNothing action")

	return gc.stream.Send(&mafia.Action{
		Type: mafia.ActionType_DoNothing,
		Data: &mafia.Action_DoNothing_{
			DoNothing: &mafia.Action_DoNothing{},
		},
	})
}

func (gc *grpcConnect) Close() {
	gc.conn.Close()
}

func (gc *grpcConnect) Init(host string, port string) error {
	target := fmt.Sprintf("%s:%s", host, port)

	fmt.Printf("connecting to grpc server: %s\n", target)

	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return errors.Wrap(err, "dial")
	}

	gc.conn = conn
	gc.client = mafia.NewMafiaServiceClient(conn)

	return nil
}

func (gc *grpcConnect) CreateStream() error {
	stream, err := gc.client.FindGame(context.Background())
	if err != nil {
		return err
	}

	gc.stream = stream
	return nil
}

func (gc *grpcConnect) CloseStream() error {
	return gc.stream.CloseSend()
}

func (gc *grpcConnect) StartListening() error {
	for {
		if Game.State == Undefined {
			Printer.PrintLine("log", "game stopped -> stop listening")
			return nil
		}

		event, err := gc.stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		Printer.PrintLine("log", fmt.Sprintf("read event: %d", int32(event.Type)))

		switch event.Type {
		case mafia.EventType_GameStart:
			Game.HandleGameStart(event.GetGameStart())
		case mafia.EventType_VoteRequest:
			Game.HandleVoteRequest(event.GetVoteRequest())
		case mafia.EventType_SystemMessage:
			Game.HandleSystemMessage(event.GetMessage())
		case mafia.EventType_MafiaCheckResponse:
			Game.HandleMafiaCheckResponse(event.GetMafiaCheckResponse())
		case mafia.EventType_Death:
			Game.HandleDeath(event.GetDeath())
		case mafia.EventType_GameEnd:
			Game.HandleGameEnd(event.GetGameEnd())
		}
	}
}
