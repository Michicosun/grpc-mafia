package client

import (
	game "grpc-mafia/client/game"
	human "grpc-mafia/client/game/human"
)

func MakeClient(use_bot bool) game.IInteractor {
	return makeHumanClient()
}

func makeHumanClient() game.IInteractor {
	game.Session.Interactor = human.MakeHumanInteractor()
	return game.Session.Interactor
}
