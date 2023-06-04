package game

import (
	"encoding/json"
	"fmt"
	chat "grpc-mafia/client/chat"
	"grpc-mafia/client/game"
	"grpc-mafia/client/grpc"
	"grpc-mafia/client/tracker_client"
	mafia "grpc-mafia/server/proto"
	"strconv"
	"strings"
)

func PrettyPrintJson(obj any) {
	jsonData, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		fmt.Println("Error generating JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}

func (hi *humanInteractor) Executor(in string) {
	in = strings.TrimSpace(in)
	blocks := strings.Split(in, " ")

	switch blocks[0] {
	case "":
		return
	case "login":
		if !AllowLogin() {
			fmt.Println("login is not allowed")
			return
		} else if len(blocks) != 2 {
			fmt.Println("need to provide login as parameter")
			return
		} else {
			hi.SetLogin(blocks[1])
		}
	case "add-comment":
		if !AllowTrackerQuery() {
			fmt.Println("tracker query is not allowe right now")
			return
		} else if len(blocks) < 3 {
			fmt.Println("you need to provide round_id and comment text")
			return
		} else if len(hi.login) == 0 {
			fmt.Println("login not specified")
			return
		} else {
			text := strings.Join(blocks[2:], " ")

			resp, err := tracker_client.TrackerClient.AddComment(blocks[1], hi.login, text)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			PrettyPrintJson(resp.GetAddComment())
		}
	case "get-round-info":
		if !AllowTrackerQuery() {
			fmt.Println("tracker query is not allowe right now")
			return
		} else if len(blocks) != 2 {
			fmt.Println("you need to provide round_id")
			return
		} else {
			resp, err := tracker_client.TrackerClient.GetRoundInfo(blocks[1])
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			PrettyPrintJson(resp.GetGetRoundInfo())
		}
	case "list-rounds":
		if !AllowTrackerQuery() {
			fmt.Println("tracker query is not allowe right now")
			return
		} else if len(blocks) != 3 {
			fmt.Println("you need to provide n: count of rounds and comment state: RUNNING/WIN_MAFIA/WIN_SHERIFFS")
			return
		} else {
			n, err := strconv.Atoi(blocks[1])
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			resp, err := tracker_client.TrackerClient.ListRounds(n, blocks[2])
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			PrettyPrintJson(resp.GetListRounds())
		}
	case "connect":
		if !AllowConnect() {
			fmt.Println("connect is not allowed")
			return
		} else if len(blocks) != 1 {
			fmt.Println("parameters are not expected")
			return
		} else if len(hi.login) == 0 {
			fmt.Println("login not specified")
			return
		} else {
			fmt.Println("using login:", hi.login)
			game.Session.Start(hi.login)
		}
	case "message":
		if !AllowMessage() {
			fmt.Println("sending messages is not allowed")
			return
		} else if len(blocks) < 3 {
			fmt.Println("need to provide which group is this message for and message text")
			return
		} else {
			msg := strings.Join(blocks[2:], " ")
			switch blocks[1] {
			case "all":
				if !AllowMessageByRole(mafia.Role_Civilian) {
					fmt.Println("sending messages is not allowed")
					return
				} else {
					if err := chat.RabbitConnection.SendMessage(game.Session.Name, msg, int32(mafia.Role_Civilian)); err != nil {
						game.Session.StopWithError(err)
					}
				}
			default:
				if game.Session.Role == mafia.Role_Civilian || !AllowMessageByRole(game.Session.Role) {
					fmt.Println("sending messages is not allowed")
					return
				} else {
					if err := chat.RabbitConnection.SendMessage(game.Session.Name, msg, int32(game.Session.Role)); err != nil {
						game.Session.StopWithError(err)
					}
				}
			}
		}
	case "vote":
		if !AllowVote() {
			fmt.Println("vote is not allowed")
			return
		} else if len(blocks) != 2 {
			fmt.Println("need to specify who you are voting for")
			return
		} else {
			game.Session.ChangeState(game.Waiting, false)
			if err := grpc.Connection.SendVote(game.Session.Name, blocks[1]); err != nil {
				game.Session.StopWithError(err)
			}
		}
	case "nothing":
		if !AllowNothing() {
			fmt.Println("nothing is not allowed")
			return
		} else {
			game.Session.ChangeState(game.Waiting, false)
			if err := grpc.Connection.SendDoNothing(game.Session.Name); err != nil {
				game.Session.StopWithError(err)
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
