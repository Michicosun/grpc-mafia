package server

import (
	mafia "grpc-mafia/server/proto"
	"sync"
)

type GameStorage struct {
	games   []*Game
	pending *Game

	mtx sync.Mutex
}

func (gs *GameStorage) JoinGame(name string) (<-chan *mafia.Event, *Game) {
	gs.mtx.Lock()
	defer gs.mtx.Unlock()

	game := gs.pending
	stream, is_started := game.Join(name)

	player_cnt := 4 // TODO fetch environment
	if is_started {
		gs.games = append(gs.games, gs.pending)
		gs.pending = NewGame(uint32(player_cnt))
	}

	return stream, game
}

func MakeGameStorage() *GameStorage {
	player_cnt := 4 // TODO fetch environment

	return &GameStorage{
		games:   make([]*Game, 0),
		pending: NewGame(uint32(player_cnt)),
	}
}
