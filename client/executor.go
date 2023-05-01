package client

import (
	"fmt"
	"strings"
)

func Executor(in string) {
	in = strings.TrimSpace(in)

	blocks := strings.Split(in, " ")

	switch blocks[0] {
	case "connect":
		if len(blocks) != 2 {
			fmt.Println("need to provide login as parameter")
		} else {
			Game.Init(blocks[1])
		}
	case "message":
		if len(blocks) != 3 {
			fmt.Println("need to provide which group is this message for")
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
		}
		if err := GrpcConnect.SendVote(Game.Name, blocks[1]); err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			Game.Stop()
		} else {
			Game.ChangeState(Waiting)
		}
	case "nothing":
		if err := GrpcConnect.SendDoNothing(); err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			Game.Stop()
		} else {
			Game.ChangeState(Waiting)
		}
	}
}
