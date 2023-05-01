package game

import (
	"fmt"
	"grpc-mafia/client/game"
	"grpc-mafia/client/grpc"
	"strings"
)

func (hi *humanInteractor) Executor(in string) {
	in = strings.TrimSpace(in)
	blocks := strings.Split(in, " ")

	switch blocks[0] {
	case "connect":
		if len(blocks) != 2 {
			fmt.Println("need to provide login as parameter")
			return
		} else {
			game.Session.Start(blocks[1])
		}
	case "message":
		if len(blocks) != 3 {
			fmt.Println("need to provide which group is this message for")
			return
		}
		switch blocks[1] {
		case "all":
			fmt.Printf("send to all: %s\n", blocks[2])
		default:
			fmt.Printf("send to %s: %s\n", blocks[1], blocks[2])
		}
	case "vote":
		if len(blocks) != 2 {
			fmt.Println("need to specify who you are voting for")
			return
		}
		if err := grpc.Connection.SendVote(game.Session.Name, blocks[1]); err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			game.Session.Stop()
		} else {
			game.Session.ChangeState(game.Waiting)
		}
	case "nothing":
		if err := grpc.Connection.SendDoNothing(game.Session.Name); err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			game.Session.Stop()
		} else {
			game.Session.ChangeState(game.Waiting)
		}
	case "disconnect":
		game.Session.Stop()
	case "publish":
		grpc.Connection.SendPublishRequest(game.Session.MafiaName)
	default:
		fmt.Println("unrecognized command")
	}
}
