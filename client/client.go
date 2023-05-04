package client

import (
	game "grpc-mafia/client/game"
	bot "grpc-mafia/client/game/bot"
	human "grpc-mafia/client/game/human"
)

func MakeClient(use_bot bool) game.IInteractor {
	if use_bot {
		return makeBotClient()
	}

	return makeHumanClient()
}

func makeHumanClient() game.IInteractor {
	game.Session.Interactor = human.MakeHumanInteractor()
	return game.Session.Interactor
}

func makeBotClient() game.IInteractor {
	game.Session.Interactor = bot.MakeBotInteractor()
	return game.Session.Interactor
}
