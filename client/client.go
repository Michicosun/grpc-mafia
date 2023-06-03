package client

import (
	game "grpc-mafia/client/game"
	bot "grpc-mafia/client/game/bot"
	human "grpc-mafia/client/game/human"
)

func MakeClient(use_bot bool) game.IInteractor {
	var interactor game.IInteractor

	if use_bot {
		interactor = makeBotClient()
	} else {
		interactor = makeHumanClient()
	}

	return interactor
}

func makeHumanClient() game.IInteractor {
	game.Session.Interactor = human.MakeHumanInteractor()
	return game.Session.Interactor
}

func makeBotClient() game.IInteractor {
	game.Session.Interactor = bot.MakeBotInteractor()
	return game.Session.Interactor
}
