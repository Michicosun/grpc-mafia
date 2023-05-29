package game

import (
	"fmt"
	chat "grpc-mafia/client/chat"
	"grpc-mafia/client/game"
	"grpc-mafia/client/grpc"
	mafia "grpc-mafia/server/proto"
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
		} else if !AllowConnect() {
			fmt.Println("connect is not allowed")
			return
		} else {
			game.Session.Start(blocks[1])
		}
	case "message":
		if len(blocks) != 3 {
			fmt.Println("need to provide which group is this message for")
			return
		} else {
			switch blocks[1] {
			case "all":
				if !AllowMessage(mafia.Role_Civilian) {
					fmt.Println("sending messages is not allowed")
					return
				} else {
					chat.RabbitConnection.SendMessage(game.Session.Name, blocks[2], int32(mafia.Role_Civilian))
				}
			default:
				if game.Session.Role == mafia.Role_Civilian || !AllowMessage(game.Session.Role) {
					fmt.Println("sending messages is not allowed")
					return
				} else {
					chat.RabbitConnection.SendMessage(game.Session.Name, blocks[2], int32(game.Session.Role))
				}
			}
		}
	case "vote":
		if len(blocks) != 2 {
			fmt.Println("need to specify who you are voting for")
			return
		} else if !AllowVote() {
			fmt.Println("vote is not allowed")
			return
		} else {
			game.Session.ChangeState(game.Waiting, false)
			if err := grpc.Connection.SendVote(game.Session.Name, blocks[1]); err != nil {
				fmt.Printf("ERROR: %s\n", err.Error())
				game.Session.Stop()
			}
		}
	case "nothing":
		if !AllowNothing() {
			fmt.Println("nothing is not allowed")
			return
		} else {
			game.Session.ChangeState(game.Waiting, false)
			if err := grpc.Connection.SendDoNothing(game.Session.Name); err != nil {
				fmt.Printf("ERROR: %s\n", err.Error())
				game.Session.Stop()
			}
		}
	case "disconnect":
		if !AllowDisconnect() {
			fmt.Println("disconnect is not allowed")
			return
		} else {
			game.Session.Stop()
		}
	case "publish":
		if !AllowPublish() {
			fmt.Println("publish is not allowed")
			return
		} else {
			grpc.Connection.SendPublishRequest(game.Session.MafiaName)
		}
	case "exit":
		// execution will disrupted via exit function
	default:
		fmt.Println("unrecognized command")
	}
}
