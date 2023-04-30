package server

import (
	mafia "grpc-mafia/server/proto"

	zlog "github.com/rs/zerolog/log"
)

type GameServer struct {
	mafia.UnimplementedMafiaServiceServer
	gs *GameStorage
}

func (s *GameServer) FindGame(stream mafia.MafiaService_FindGameServer) error {
	action, err := stream.Recv()
	if err != nil {
		return err
	}

	events, game := s.gs.JoinGame(action.GetInit().Name)

	go s.GameListener(stream, events)
	s.GameWriter(stream, game)

	return nil
}

func (s *GameServer) GameListener(stream mafia.MafiaService_FindGameServer, events <-chan *mafia.Event) {
	for event := range events {
		zlog.Info().Int32("type", int32(event.Type)).Msg("sending event")
		stream.Send(event)
	}
}

func (s *GameServer) GameWriter(stream mafia.MafiaService_FindGameServer, game *Game) {
	for {
		action, err := stream.Recv()
		if err != nil {
			zlog.Error().Err(err).Msg("action receive")
			return
		}

		zlog.Info().Int32("type", int32(action.Type)).Msg("received action")

		switch action.GetType() {
		case mafia.ActionType_Init:
			zlog.Info().Msg("get init message")
			panic("unexpected message")

		case mafia.ActionType_Vote:
			zlog.Info().Msg("get vote message")

			from := action.GetVote().From
			to := action.GetVote().Name

			game.Vote(from, to)

		case mafia.ActionType_DoNothing:
			zlog.Info().Msg("get do_nothing message")
			game.DoNothing()
		}
	}
}

func MakeGameServer() *GameServer {
	return &GameServer{
		gs: MakeGameStorage(),
	}
}
