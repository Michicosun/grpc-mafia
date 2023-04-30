package server

import (
	mafia "grpc-mafia/server/proto"

	zlog "github.com/rs/zerolog/log"
)

type GameServer struct {
	mafia.UnimplementedMafiaServiceServer
}

func (s *GameServer) FindGame(stream mafia.MafiaService_FindGameServer) error {
	for {
		action, err := stream.Recv()
		if err != nil {
			return err
		}

		switch action.GetType() {
		case mafia.ActionType_Init:
			zlog.Info().Msg("get init message")
			stream.Send(&mafia.Event{
				Type: mafia.EventType_SystemMessage,
				Data: &mafia.Event_Message{
					Message: &mafia.Event_SystemMessage{
						Text: "init",
					},
				},
			})

		case mafia.ActionType_Vote:
			zlog.Info().Msg("get vote message")
			stream.Send(&mafia.Event{
				Type: mafia.EventType_SystemMessage,
				Data: &mafia.Event_Message{
					Message: &mafia.Event_SystemMessage{
						Text: "vote",
					},
				},
			})

		case mafia.ActionType_DoNothing:
			zlog.Info().Msg("get do_nothing message")
			stream.Send(&mafia.Event{
				Type: mafia.EventType_SystemMessage,
				Data: &mafia.Event_Message{
					Message: &mafia.Event_SystemMessage{
						Text: "do nothing",
					},
				},
			})
		}
	}
}
