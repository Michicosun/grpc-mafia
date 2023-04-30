package client

type GamePhase uint64

const (
	Init  GamePhase = 0
	Day   GamePhase = 1
	Night GamePhase = 2
)

var GameState = &gameState{}

type gameState struct {
	phase GamePhase
}

func (g *gameState) Init() {
	g.phase = Init
}
