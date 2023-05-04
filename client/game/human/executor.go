package game

import (
	"fmt"
	"grpc-mafia/chat"
	"grpc-mafia/client/game"
	"grpc-mafia/client/grpc"
	"grpc-mafia/util"
	"strings"
)

func (hi *humanInteractor) Executor(in string) {
	in = strings.TrimSpace(in)
	blocks := strings.Split(in, " ")

	switch blocks[0] {
	case "":
		return
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
			chat.Connector.MakeBCast(util.ChatGroupName(game.Session.SessionId, "all"), chat.Message{
				From: game.Session.Name,
				Text: blocks[2],
			})
		default:
			chat.Connector.MakeBCast(util.ChatGroupName(game.Session.SessionId, game.Session.Role.String()), chat.Message{
				From: game.Session.Name,
				Text: blocks[2],
			})
		}
	case "vote":
		if len(blocks) != 2 {
			fmt.Println("need to specify who you are voting for")
			return
		}
		game.Session.ChangeState(game.Waiting, false)
		if err := grpc.Connection.SendVote(game.Session.Name, blocks[1]); err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			game.Session.Stop()
		}
	case "nothing":
		game.Session.ChangeState(game.Waiting, false)
		if err := grpc.Connection.SendDoNothing(game.Session.Name); err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			game.Session.Stop()
		}
	case "disconnect":
		game.Session.Stop()
	case "publish":
		grpc.Connection.SendPublishRequest(game.Session.MafiaName)
	case "exit":
		// execution will disrupted via exit function
	default:
		fmt.Println("unrecognized command")
	}
}
