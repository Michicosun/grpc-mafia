package client

import (
	mafia "grpc-mafia/server/proto"

	"github.com/c-bata/go-prompt"
)

func Completer(in prompt.Document) []prompt.Suggest {
	Parser.cur_buffer = in.Text

	s := stateToSuggestions[Game.State]
	word := in.GetWordBeforeCursor()

	if word == "message" {
		s = makeMessageSuggestions()
	} else if word == "vote" {
		s = makeAlivePlayersSuggestions()
	}

	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

var stateToSuggestions = map[GameState][]prompt.Suggest{
	Undefined: {
		{Text: "connect", Description: "find game"},
	},
	Waiting: {
		{Text: "message", Description: "send message"},
	},
	PrepareState: {
		{Text: "message", Description: "send message"},
		{Text: "nothing", Description: "pass this turn"},
	},
	NeedVote: {
		{Text: "message", Description: "send message"},
		{Text: "vote", Description: "vote in this turn"},
	},
}

var msg_suggests = []prompt.Suggest{
	{Text: "all", Description: "send message to all players"},
}

func makeMessageSuggestions() []prompt.Suggest {
	msg_suggests := msg_suggests

	if Game.Role == mafia.Role_Mafia {
		msg_suggests = append(msg_suggests, prompt.Suggest{
			Text: "mafias", Description: "send to all mafias",
		})
	} else if Game.Role == mafia.Role_Sheriff {
		msg_suggests = append(msg_suggests, prompt.Suggest{
			Text: "sheriffs", Description: "send to all sheriffs",
		})
	}

	return msg_suggests
}

func makeAlivePlayersSuggestions() []prompt.Suggest {
	players := make([]prompt.Suggest, 0)

	for player := range Game.AlivePlayers {
		players = append(players, prompt.Suggest{
			Text: player, Description: "alive player",
		})
	}

	return players
}

func exitChecker(in string, breakline bool) bool {
	return in == "exit" && breakline
}
