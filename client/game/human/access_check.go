package game

import (
	"grpc-mafia/client/game"
	mafia "grpc-mafia/server/proto"
)

func AllowLogin() bool {
	return game.Session.GetState() == game.Undefined
}

func AllowConnect() bool {
	return game.Session.GetState() == game.Undefined
}

func AllowDisconnect() bool {
	return game.Session.GetState() != game.Undefined
}

func AllowNothing() bool {
	return game.Session.GetState() == game.PrepareState
}

func AllowPublish() bool {
	return game.Session.MafiaCheck
}

func AllowMessage() bool {
	return game.Session.GetState() != game.Undefined
}

func AllowMessageByRole(role mafia.Role) bool {
	if game.Session.GetState() == game.Ghost {
		return false
	}

	if role == mafia.Role_Civilian {
		return game.Session.GetTime() == mafia.Time_Day
	}

	return game.Session.GetTime() == mafia.Time_Night
}

func AllowVote() bool {
	return game.Session.GetState() == game.NeedVote
}
