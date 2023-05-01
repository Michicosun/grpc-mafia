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

	name := action.GetInit().Name
	zlog.Info().Str("name", name).Str("addr", p.Addr.String()).Msg("get connection of user")

	events, game := s.gs.JoinGame(name)

	go func() {
		if err = s.GameListener(stream, events); err != nil {
			zlog.Error().Err(err).Msg("write to stream error")
		}
	}()

	if err = s.GameWriter(stream, game); err != nil {
		zlog.Error().Err(err).Msg("read from stream error")
		game.Disconnect(name)
		return err
	}

	return err
}

func (s *GameServer) GameListener(stream mafia.MafiaService_FindGameServer, events <-chan *mafia.Event) error {
	for event := range events {
		zlog.Info().Str("type", event.GetType().String()).Msg("sending event")
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

		zlog.Info().Str("type", action.GetType().String()).Msg("received action")

		switch action.GetType() {
		case mafia.ActionType_Vote:
			from := action.GetVote().From
			to := action.GetVote().Name

			zlog.Info().Str("from", from).Str("to", to).Msg("get vote message")

			game.Vote(from, to)

		case mafia.ActionType_DoNothing:
			from := action.GetDoNothing().From

			zlog.Info().Str("from", from).Msg("get do_nothing message")

			game.DoNothing(from)

		case mafia.ActionType_PublishRequest:
			mafia_name := action.GetPublishRequest().MafiaName

			zlog.Info().Str("mafia_name", mafia_name).Msg("get publish_request message")

			game.Publish(mafia_name)

		default:
			panic("unexpected message")
		}
	}
}

func MakeGameServer() *GameServer {
	return &GameServer{
		gs: MakeGameStorage(),
	}
}
