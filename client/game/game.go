package game

import (
	"fmt"
	"grpc-mafia/client/chat"
	"grpc-mafia/client/grpc"
	mafia "grpc-mafia/server/proto"
)

var Session = &session{}

type GameState int32

const (
	Undefined    = -1
	Waiting      = 0
	PrepareState = 1
	NeedVote     = 2
	Ghost        = 3
)

type session struct {
	SessionId string
	Name      string

	state        GameState
	time         mafia.Time
	Role         mafia.Role
	AlivePlayers map[string]struct{}
	Group        map[string]struct{}

	MafiaName  string
	MafiaCheck bool

	Interactor IInteractor
}

func (s *session) ChangeState(new_state GameState, use_signal bool) {
	s.state = new_state

	if use_signal && s.Interactor != nil {
		s.Interactor.Signal()
		RefreshLine(s.Interactor)
	}
}

func (s *session) GetState() GameState {
	return s.state
}

func (s *session) ChangeTime(received_time mafia.Time) {
	s.time = received_time
}

func (s *session) GetTime() mafia.Time {
	return s.time
}

func (s *session) ClearMafiaCheck() {
	s.MafiaCheck = false
	s.MafiaName = ""
}

func (s *session) Init() {
	s.Role = mafia.Role_Civilian
	s.ChangeState(Undefined, true)
}

func (s *session) Start(name string) error {
	s.clearGame()

	if err := grpc.Connection.CreateStream(); err != nil {
		return err
	}

	if err := grpc.Connection.SendInit(name); err != nil {
		return err
	}

	s.Name = name
	s.ChangeState(Waiting, false)
	go startListening(s.Interactor)

	return nil
}

func (s *session) Stop() {
	s.ChangeState(Undefined, true)
	grpc.Connection.CloseStream()

	ChatPrinter.Stop()
	if err := chat.RabbitConnection.Close(); err != nil {
		PrintLine("system", err.Error(), s.Interactor)
	}
}

func (s *session) StopWithError(err error) {
	fmt.Printf("ERROR: %s\n", err.Error())
	s.Stop()
}

func (s *session) HandleGameStart(e *mafia.Event_GameStart) {
	s.SessionId = e.GetSessionId()
	s.Role = e.Role
	s.AlivePlayers = make(map[string]struct{})
	s.Group = make(map[string]struct{})

	if err := chat.RabbitConnection.Establish(s.SessionId, s.Role, s.Name); err != nil {
		s.StopWithError(err)
	} else {
		ChatPrinter.Start()
	}

	for _, player := range e.Players {
		s.AlivePlayers[player] = struct{}{}
	}

	for _, player := range e.Group {
		s.Group[player] = struct{}{}
	}

	s.ChangeState(PrepareState, true)
}

func (s *session) HandleVoteRequest(e *mafia.Event_VoteRequest) {
	s.ChangeState(NeedVote, true)
}

func (s *session) HandleSystemMessage(e *mafia.Event_SystemMessage) {
	PrintLine("system", e.GetText(), s.Interactor)
}

func (s *session) HandleMafiaCheckResponse(e *mafia.Event_MafiaCheckResponse) {
	s.MafiaCheck = e.IsMafia
	s.MafiaName = e.Name
}

func (s *session) HandlePublish(e *mafia.Event_Publish) {
	PrintLine("system", fmt.Sprintf("sheriffs revealed mafia: %s", e.MafiaName), s.Interactor)

	// publish is oneshot
	s.ClearMafiaCheck()
}

func (s *session) HandleDeath(e *mafia.Event_Death) {
	delete(s.AlivePlayers, e.Name)

	if e.Name == s.Name {
		// user was killed
		s.ChangeState(Ghost, false)
	}
}

func (s *session) HandleGameEnd(e *mafia.Event_GameEnd) {
	PrintLine("system", e.GetText(), s.Interactor)
	s.Stop()
}

func (s *session) clearGame() {
	s.SessionId = ""
	s.Name = ""
	s.state = Undefined
	s.Role = mafia.Role_Civilian
	s.ClearMafiaCheck()
}
