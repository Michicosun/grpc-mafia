package client

import (
	mafia "grpc-mafia/server/proto"
)

var Game = &game{}

type GameState int32

const (
	Undefined    = -1
	Waiting      = 0
	PrepareState = 1
	NeedVote     = 2
	Ghost        = 3
)

type game struct {
	Name string

	State        GameState
	Role         mafia.Role
	AlivePlayers map[string]struct{}
	Group        map[string]struct{}

	MafiaName  string
	MafiaCheck bool
}

func (g *game) ChangeState(new_state GameState) {
	g.State = new_state
}

func (g *game) Init() {
	g.State = Undefined
	g.Role = mafia.Role_Civilian
}

func (g *game) Start(name string) error {
	if err := GrpcConnect.CreateStream(); err != nil {
		return err
	}

	if err := GrpcConnect.SendInit(name); err != nil {
		return err
	}

	g.Name = name
	g.ChangeState(Waiting)
	go GrpcConnect.StartListening()

	return nil
}

func (g *game) Stop() {
	g.ChangeState(Undefined)
	GrpcConnect.CloseStream()
}

func (g *game) HandleGameStart(e *mafia.Event_GameStart) {
	g.Role = e.Role
	g.AlivePlayers = make(map[string]struct{})
	g.Group = make(map[string]struct{})

	for _, player := range e.Players {
		g.AlivePlayers[player] = struct{}{}
	}

	for _, player := range e.Group {
		g.Group[player] = struct{}{}
	}

	g.ChangeState(PrepareState)
}

func (g *game) HandleVoteRequest(e *mafia.Event_VoteRequest) {
	g.ChangeState(NeedVote)
}

func (g *game) HandleSystemMessage(e *mafia.Event_SystemMessage) {
	Printer.PrintLine("system", e.GetText())
}

func (g *game) HandleMafiaCheckResponse(e *mafia.Event_MafiaCheckResponse) {
	g.MafiaCheck = e.IsMafia
	g.MafiaName = e.Name
}

func (g *game) HandleDeath(e *mafia.Event_Death) {
	delete(g.AlivePlayers, e.Name)

	if e.Name == g.Name {
		// user was killed
		g.ChangeState(Ghost)
	}
}

func (g *game) HandleGameEnd(e *mafia.Event_GameEnd) {
	Printer.PrintLine("system", e.GetText())
	g.Stop()
}
