package server

import (
	mafia "grpc-mafia/server/proto"
	"grpc-mafia/util"
	"strconv"
	"sync"
)

type GameStorage struct {
	players_cnt uint32
	games       []*Game
	pending     *Game

	mtx sync.Mutex
}

func (gs *GameStorage) JoinGame(name string, connect string) (<-chan *mafia.Event, *Game) {
	gs.mtx.Lock()
	defer gs.mtx.Unlock()

	game := gs.pending
	stream, is_started := game.Join(name, connect)

	if is_started {
		gs.games = append(gs.games, gs.pending)
		gs.pending = NewGame(gs.players_cnt)
	}

	return stream, game
}

func MakeGameStorage() *GameStorage {
	players_cnt, err := strconv.Atoi(util.GetEnvWithDefault("PLAYERS_CNT", "4"))
	if err != nil {
		players_cnt = 4
	}

	return &GameStorage{
		players_cnt: uint32(players_cnt),
		games:       make([]*Game, 0),
		pending:     NewGame(uint32(players_cnt)),
	}
}
