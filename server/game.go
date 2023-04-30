package server

import (
	"fmt"
	mafia "grpc-mafia/server/proto"
	"sync"
)

type Player struct {
	send_chan chan *mafia.Event
}

type GameState uint32

const (
	Day          = 0
	Night        = 1
	PrepareState = 3
	WinMafia     = 4
	WinSheriffs  = 5
)

type Game struct {
	players_cnt   uint32
	playersInfo   map[string]Player
	alive_players map[string]struct{}
	ghosts        map[string]struct{}
	mafia         map[string]struct{}
	sheriffs      map[string]struct{}

	state       GameState
	need_events int
	kill_votes  []string
	check_votes []string

	mtx     sync.Mutex
	actions sync.Cond
}

func (g *Game) Join(name string) (<-chan *mafia.Event, bool) {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	stream := make(chan *mafia.Event, 1)

	g.playersInfo[name] = Player{
		send_chan: stream,
	}

	g.alive_players[name] = struct{}{}

	g.msgToAll(fmt.Sprintf("connected new player: %s | %d/%d", name, len(g.alive_players), g.players_cnt))

	is_started := len(g.alive_players) == int(g.players_cnt)
	if is_started {
		go g.Start()
	}

	return stream, is_started
}

func (g *Game) DoNothing() {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	g.need_events -= 1

	if g.need_events == 0 {
		g.actions.Signal()
	}
}

func (g *Game) Vote(from string, to string) {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	g.need_events -= 1

	if g.isMafia(from) {
		g.kill_votes = append(g.kill_votes, to)
	} else if g.isSheriff(from) {
		g.check_votes = append(g.check_votes, to)
	} else {
		g.kill_votes = append(g.kill_votes, to)
	}

	if g.need_events == 0 {
		g.actions.Signal()
	}
}

func (g *Game) msgToGroup(grp map[string]struct{}, msg string) {
	for player := range grp {
		g.playersInfo[player].send_chan <- &mafia.Event{
			Type: mafia.EventType_SystemMessage,
			Data: &mafia.Event_Message{
				Message: &mafia.Event_SystemMessage{
					Text: msg,
				},
			},
		}
	}
}

func (g *Game) msgToAll(msg string) {
	g.msgToGroup(g.alive_players, msg)
	g.msgToGroup(g.ghosts, msg)
	g.msgToGroup(g.mafia, msg)
	g.msgToGroup(g.sheriffs, msg)
}

func (g *Game) requestVotes(grp map[string]struct{}) {
	for player := range grp {
		g.playersInfo[player].send_chan <- &mafia.Event{
			Type: mafia.EventType_VoteRequest,
			Data: &mafia.Event_VoteRequest_{
				VoteRequest: &mafia.Event_VoteRequest{},
			},
		}
	}
}

func (g *Game) sendCheckResult(to_check string) {
	is_mafia := g.isMafia(to_check)

	for player := range g.sheriffs {
		g.playersInfo[player].send_chan <- &mafia.Event{
			Type: mafia.EventType_MafiaCheckResponse,
			Data: &mafia.Event_MafiaCheckResponse_{
				MafiaCheckResponse: &mafia.Event_MafiaCheckResponse{
					Name:    to_check,
					IsMafia: is_mafia,
				},
			},
		}
	}
}

func (g *Game) Start() {
	// prepare phase

	g.mtx.Lock()
	defer g.mtx.Unlock()

	g.need_events = len(g.alive_players)

	g.msgToAll("starting game")
	g.requestVotes(g.alive_players)
	g.actions.Wait()

	for {
		g.msgToAll("day started")
		g.day()

		if g.isEnded() {
			g.end()
			return
		}

		g.msgToAll("night started")
		g.night()

		if g.isEnded() {
			g.end()
			return
		}
	}
}

func (g *Game) getMostRecent(arr []string) string {
	cnt := make(map[string]uint32)

	max_cnt := 0
	most_recent := ""
	for _, elem := range arr {
		cnt[elem] += 1

		if cnt[elem] > uint32(max_cnt) {
			most_recent = elem
		}
	}

	return most_recent
}

func (g *Game) day() {
	g.kill_votes = make([]string, 0)
	g.need_events = len(g.alive_players)

	g.msgToGroup(g.alive_players, "vote who is the mafia")
	g.requestVotes(g.alive_players)
	g.actions.Wait()

	to_kill := g.getMostRecent(g.kill_votes)
	g.kill(to_kill)

	g.msgToAll(fmt.Sprintf("%s was killed", to_kill))
	g.changeState()
}

func (g *Game) night() {
	g.kill_votes = make([]string, 0)
	g.check_votes = make([]string, 0)
	g.need_events = len(g.mafia) + len(g.sheriffs)

	g.msgToGroup(g.mafia, "vote who to kill")
	g.requestVotes(g.mafia)
	g.msgToGroup(g.sheriffs, "vote who to check")
	g.requestVotes(g.sheriffs)
	g.actions.Wait()

	to_kill := g.getMostRecent(g.kill_votes)
	g.kill(to_kill)

	to_check := g.getMostRecent(g.check_votes)
	g.sendCheckResult(to_check)

	g.msgToAll(fmt.Sprintf("%s was killed", to_kill))
	g.changeState()
}

func (g *Game) isMafia(player string) bool {
	_, ok := g.mafia[player]
	return ok
}

func (g *Game) isSheriff(player string) bool {
	_, ok := g.sheriffs[player]
	return ok
}

func (g *Game) kill(player string) {
	delete(g.alive_players, player)
	g.ghosts[player] = struct{}{}

	if g.isMafia(player) {
		delete(g.mafia, player)
	}

	if g.isSheriff(player) {
		delete(g.mafia, player)
	}
}

func (g *Game) isEnded() bool {
	return g.state == WinMafia || g.state == WinSheriffs
}

func (g *Game) changeState() {
	if len(g.mafia) == 0 {
		g.state = WinSheriffs
	} else if 2*len(g.mafia) > len(g.alive_players) {
		g.state = WinMafia
	} else {
		g.state = (g.state + 1) % 2
	}
}

func (g *Game) end() {
	// TODO: check the winner group
	g.msgToAll("game ended")
}

func NewGame(player_cnt uint32) *Game {
	game := &Game{
		players_cnt:   player_cnt,
		playersInfo:   make(map[string]Player),
		alive_players: make(map[string]struct{}),
		ghosts:        make(map[string]struct{}),
		mafia:         make(map[string]struct{}),
		sheriffs:      make(map[string]struct{}),

		state:       PrepareState,
		need_events: int(player_cnt),
		kill_votes:  make([]string, 0),
		check_votes: make([]string, 0),
	}

	game.actions = sync.Cond{
		L: &game.mtx,
	}

	return game
}
