package server

import (
	"context"
	"fmt"
	"grpc-mafia/server/tracker_client_impl"
	"grpc-mafia/util"
	"net/http"
	"time"

	"github.com/Khan/genqlient/graphql"
	zlog "github.com/rs/zerolog/log"
)

var TrackerClient = &trackerClient{}

type trackerClient struct {
	client graphql.Client
}

func (tc *trackerClient) CreateRound(g *Game) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	round_info := tracker_client_impl.RoundInfo{
		Id:         g.session_id,
		State:      tracker_client_impl.RoundStateRunning,
		Started_at: g.game_start_timestamp.Format(time.RFC3339),
	}

	player_infos := make([]tracker_client_impl.PlayerInfo, 0)
	for player := range g.alive_players {
		player_infos = append(player_infos, tracker_client_impl.PlayerInfo{
			Login: player,
			Role:  g.playersRole[player].String(),
		})
	}

	zlog.Info().Time("start", g.game_start_timestamp).Str("id", g.session_id).Msg("creating game state in tracker")

	_, err := tracker_client_impl.CreateRound(ctx, tc.client, round_info, player_infos)
	if err != nil {
		zlog.Error().Err(err).Msg("create round in tracker")
	}
}

func (tc *trackerClient) UpdatePlayersInfo(g *Game) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	player_statuses := make([]tracker_client_impl.PlayerStatus, 0)
	for player := range g.playersRole {
		_, alive := g.alive_players[player]
		player_statuses = append(player_statuses, tracker_client_impl.PlayerStatus{
			Login: player,
			Alive: alive,
		})
	}

	state := tracker_client_impl.RoundStateRunning
	if g.state == WinMafia {
		state = tracker_client_impl.RoundStateWinMafia
	} else if g.state == WinSheriffs {
		state = tracker_client_impl.RoundStateWinSheriffs
	}

	zlog.Info().Str("new_state", string(state)).Interface("player_statuses", player_statuses).Msg("updating game state in tracker")

	_, err := tracker_client_impl.UpdateRound(ctx, tc.client, g.session_id, state, player_statuses)
	if err != nil {
		zlog.Error().Err(err).Msg("update round info")
	}
}

func (tc *trackerClient) Init() {
	host := util.GetEnvWithDefault("TRACKER_HOST", "localhost")
	port := util.GetEnvWithDefault("TRACKER_PORT", "9090")

	endpoint := fmt.Sprintf("http://%s:%s/graphql", host, port)

	tc.client = graphql.NewClient(endpoint, http.DefaultClient)
}
