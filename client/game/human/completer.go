package game

import (
	"fmt"
	"grpc-mafia/client/game"
	mafia "grpc-mafia/server/proto"
	"strings"

	"github.com/c-bata/go-prompt"
)

func (hi *humanInteractor) Completer(in prompt.Document) []prompt.Suggest {
	hi.cur_buffer = in.Text

	s := stateToSuggestions[game.Session.State]
	args := strings.Split(in.TextBeforeCursor(), " ")

	if len(args) > 0 {
		switch args[0] {
		case "connect":
			s = []prompt.Suggest{}
		case "exit":
			s = []prompt.Suggest{}
		case "nothing":
			s = []prompt.Suggest{}
		case "message":
			s = makeMessageSuggestions()
		case "vote":
			s = makeAlivePlayersSuggestions()
		}
	}

	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func (hi *humanInteractor) exitChecker(in string, breakline bool) bool {
	if in == "exit" && breakline {
		fmt.Print("\r")
		return true
	}
	return false
}

var stateToSuggestions = map[game.GameState][]prompt.Suggest{
	game.Undefined: {
		{Text: "connect", Description: "find game"},
		{Text: "exit", Description: "close client"},
	},
	game.Waiting: {
		{Text: "message", Description: "send message"},
		{Text: "exit", Description: "close client"},
	},
	game.PrepareState: {
		{Text: "message", Description: "send message"},
		{Text: "nothing", Description: "pass this turn"},
		{Text: "exit", Description: "close client"},
	},
	game.NeedVote: {
		{Text: "message", Description: "send message"},
		{Text: "vote", Description: "vote in this turn"},
		{Text: "exit", Description: "close client"},
	},
	game.Ghost: {
		{Text: "disconnect", Description: "exit this game"},
		{Text: "exit", Description: "close client"},
	},
}

var msg_suggests = []prompt.Suggest{
	{Text: "all", Description: "send message to all players"},
}

func makeMessageSuggestions() []prompt.Suggest {
	msg_suggests := msg_suggests

	if game.Session.Role == mafia.Role_Mafia {
		msg_suggests = append(msg_suggests, prompt.Suggest{
			Text: "mafias", Description: "send to all mafias",
		})
	} else if game.Session.Role == mafia.Role_Sheriff {
		msg_suggests = append(msg_suggests, prompt.Suggest{
			Text: "sheriffs", Description: "send to all sheriffs",
		})
	}

	return msg_suggests
}

func makeAlivePlayersSuggestions() []prompt.Suggest {
	players := make([]prompt.Suggest, 0)

	for player := range game.Session.AlivePlayers {
		players = append(players, prompt.Suggest{
			Text: player, Description: "alive player",
		})
	}

	return players
}
