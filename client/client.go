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

	// runChatPrinter(interactor)

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

// func runChatPrinter(interactor game.IInteractor) {
// 	go func() {
// 		for {
// 			msg, err := chat.Connector.RecvMessage()
// 			if err != nil {
// 				game.PrintLine("ERROR", err.Error(), interactor)
// 			} else {
// 				game.PrintLine(msg.From, msg.Text, interactor)
// 			}
// 		}
// 	}()
// }
