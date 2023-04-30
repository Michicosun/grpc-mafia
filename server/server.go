package server

import (
	"google.golang.org/grpc/peer"

	mafia "grpc-mafia/server/proto"
	"io"

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

	p, _ := peer.FromContext(stream.Context())
	// TODO: check ok value

	zlog.Info().Str("addr", p.Addr.String()).Msg("get connection of user")

	events, game := s.gs.JoinGame(action.GetInit().Name)

	go func() {
		if err = s.GameListener(stream, events); err != nil {
			zlog.Error().Err(err).Msg("write to stream error")
		}
	}()

	if err = s.GameWriter(stream, game); err != nil {
		zlog.Error().Err(err).Msg("read from stream error")
		return err
	}

	return err
}

func (s *GameServer) GameListener(stream mafia.MafiaService_FindGameServer, events <-chan *mafia.Event) error {
	for event := range events {
		zlog.Info().Int32("type", int32(event.Type)).Msg("sending event")
		err := stream.Send(event)
		if err != nil {
			return err
		}
	}

	zlog.Info().Msg("game ended -> listener disabled")

	return nil
}

func (s *GameServer) GameWriter(stream mafia.MafiaService_FindGameServer, game *Game) error {
	for {
		if game.isEnded() {
			zlog.Info().Msg("game ended -> listener disabled")
			return nil
		}

		action, err := stream.Recv()
		if err == io.EOF {
			zlog.Info().Msg("client disconnected -> writer disabled")
			return nil
		}

		if err != nil {
			return err
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
