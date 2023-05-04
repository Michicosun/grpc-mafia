package grpc

import (
	"context"
	"fmt"
	"grpc-mafia/chat"
	mafia "grpc-mafia/server/proto"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Connection = &connection{}

type connection struct {
	conn   *grpc.ClientConn
	client mafia.MafiaServiceClient
	stream mafia.MafiaService_FindGameClient
}

func (gc *connection) SendInit(name string) error {
	fmt.Println("sending Init Game action")

	return gc.stream.Send(&mafia.Action{
		Type: mafia.ActionType_Init,
		Data: &mafia.Action_Init_{
			Init: &mafia.Action_Init{
				Name:     name,
				ChatPort: chat.Connector.GetLocalPort(),
			},
		},
	})
}

func (gc *connection) SendVote(from string, to string) error {
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

func (gc *connection) SendDoNothing(from string) error {
	fmt.Println("sending DoNothing action")

	return gc.stream.Send(&mafia.Action{
		Type: mafia.ActionType_DoNothing,
		Data: &mafia.Action_DoNothing_{
			DoNothing: &mafia.Action_DoNothing{
				From: from,
			},
		},
	})
}

func (gc *connection) SendPublishRequest(mafia_name string) error {
	fmt.Println("sending PublishRequest action")

	return gc.stream.Send(&mafia.Action{
		Type: mafia.ActionType_PublishRequest,
		Data: &mafia.Action_PublishRequest_{
			PublishRequest: &mafia.Action_PublishRequest{
				MafiaName: mafia_name,
			},
		},
	})
}

func (gc *connection) Close() {
	gc.conn.Close()
}

func (gc *connection) Init(host string, port string) error {
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

func (gc *connection) CreateStream() error {
	stream, err := gc.client.FindGame(context.Background())
	if err != nil {
		return err
	}

	gc.stream = stream
	return nil
}

func (gc *connection) CloseStream() error {
	return gc.stream.CloseSend()
}

func (gc *connection) GetStream() mafia.MafiaService_FindGameClient {
	return gc.stream
}
