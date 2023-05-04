package game

import (
	"fmt"
	mafia "grpc-mafia/server/proto"

	"grpc-mafia/client/grpc"
	"io"
)

func startListening(interactor IInteractor) error {
	stream := grpc.Connection.GetStream()

	for {
		if Session.GetState() == Undefined {
			PrintLine("log", "game stopped -> stop listening", interactor)
			return nil
		}

		event, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		PrintLine("log", fmt.Sprintf("read event: %s", event.GetType().String()), interactor)

		switch event.GetType() {
		case mafia.EventType_GameStart:
			Session.HandleGameStart(event.GetGameStart())
		case mafia.EventType_VoteRequest:
			Session.HandleVoteRequest(event.GetVoteRequest())
		case mafia.EventType_SystemMessage:
			Session.HandleSystemMessage(event.GetMessage())
		case mafia.EventType_MafiaCheckResponse:
			Session.HandleMafiaCheckResponse(event.GetMafiaCheckResponse())
		case mafia.EventType_Publish:
			Session.HandlePublish(event.GetPublish())
		case mafia.EventType_Death:
			Session.HandleDeath(event.GetDeath())
		case mafia.EventType_GameEnd:
			Session.HandleGameEnd(event.GetGameEnd())
		}
	}
}
