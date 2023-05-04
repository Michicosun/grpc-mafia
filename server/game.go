package server

import (
	"fmt"
	mafia "grpc-mafia/server/proto"
	"math/rand"
	"sync"
)

type Player struct {
	send_chan chan *mafia.Event
}

type GameState uint32

const (
	Day          = 0
	Night        = 1
	PrepareState = 2
	WinMafia     = 3
	WinSheriffs  = 4
)

type Game struct {
	players_cnt uint32
	playersInfo map[string]Player
	playersRole map[string]mafia.Role

	alive_players map[string]struct{}
	ghosts        map[string]struct{}
	mafia         map[string]struct{}
	sheriffs      map[string]struct{}

	state            GameState
	need_events_from map[string]struct{}
	kill_votes       []string
	check_votes      []string

	mtx     sync.Mutex
	actions sync.Cond
}

/////////////////////////////////////////////
// Control methods
/////////////////////////////////////////////

func (g *Game) Join(name string) (<-chan *mafia.Event, bool) {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	stream := make(chan *mafia.Event, 2)

	g.playersInfo[name] = Player{
		send_chan: stream,
	}

	g.alive_players[name] = struct{}{}

	g.sendMsgToAll(fmt.Sprintf("connected new player: %s | %d/%d", name, len(g.alive_players), g.players_cnt))

	is_started := len(g.alive_players) == int(g.players_cnt)
	if is_started {
		go g.Start()
	}

	return stream, is_started
}

func (g *Game) DoNothing(from string) {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	delete(g.need_events_from, from)

	if len(g.need_events_from) == 0 {
		g.actions.Signal()
	}
}

func (g *Game) Vote(from string, to string) {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	delete(g.need_events_from, from)

	if g.state == Day {
		g.kill_votes = append(g.kill_votes, to)
	} else if g.isMafia(from) {
		g.kill_votes = append(g.kill_votes, to)
	} else if g.isSheriff(from) {
		g.check_votes = append(g.check_votes, to)
	} else {
		panic("broken state")
	}

	if len(g.need_events_from) == 0 {
		g.actions.Signal()
	}
}

func (g *Game) Publish(mafia_name string) {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	g.sendToGroup(g.alive_players, &mafia.Event{
		Type: mafia.EventType_Publish,
		Data: &mafia.Event_Publish_{
			Publish: &mafia.Event_Publish{
				MafiaName: mafia_name,
			},
		},
	})
}

func (g *Game) Disconnect(player string) {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	// disconnect == death
	_, is_alive := g.alive_players[player]
	if is_alive {
		g.sendMsgToAll(fmt.Sprintf("%s disconnected", player))
		g.kill(player)
	}
	delete(g.ghosts, player)

	// remove from wait list
	delete(g.need_events_from, player)

	if len(g.need_events_from) == 0 {
		g.actions.Signal()
	}
}

/////////////////////////////////////////////
// Send Message methods
/////////////////////////////////////////////

func mapToArray(grp map[string]struct{}) []string {
	keys := make([]string, 0)

	for k := range grp {
		keys = append(keys, k)
	}

	return keys
}

func (g *Game) getRoleGroup(player string) []string {
	if g.isMafia(player) {
		return mapToArray(g.mafia)
	} else if g.isSheriff(player) {
		return mapToArray(g.sheriffs)
	} else {
		return make([]string, 0)
	}
}

func (g *Game) sendToGroup(grp map[string]struct{}, msg *mafia.Event) {
	for player := range grp {
		g.playersInfo[player].send_chan <- msg
	}
}

func (g *Game) sendStartGame() {
	for player := range g.alive_players {
		g.playersInfo[player].send_chan <- &mafia.Event{
			Type: mafia.EventType_GameStart,
			Data: &mafia.Event_GameStart_{
				GameStart: &mafia.Event_GameStart{
					Role:    g.playersRole[player],
					Players: mapToArray(g.alive_players),
					Group:   g.getRoleGroup(player),
				},
			},
		}
	}
}

func (g *Game) sendMsgToGroup(grp map[string]struct{}, msg string) {
	g.sendToGroup(grp, &mafia.Event{
		Type: mafia.EventType_SystemMessage,
		Data: &mafia.Event_Message{
			Message: &mafia.Event_SystemMessage{
				Text: msg,
			},
		},
	})
}

func (g *Game) sendMsgToAll(msg string) {
	g.sendMsgToGroup(g.alive_players, msg)
	g.sendMsgToGroup(g.ghosts, msg)
}

func (g *Game) requestVotes(grp map[string]struct{}) {
	g.sendToGroup(grp, &mafia.Event{
		Type: mafia.EventType_VoteRequest,
		Data: &mafia.Event_VoteRequest_{
			VoteRequest: &mafia.Event_VoteRequest{},
		},
	})
}

func (g *Game) sendCheckResult(to_check string) {
	is_mafia := g.isMafia(to_check)

	g.sendToGroup(g.sheriffs, &mafia.Event{
		Type: mafia.EventType_MafiaCheckResponse,
		Data: &mafia.Event_MafiaCheckResponse_{
			MafiaCheckResponse: &mafia.Event_MafiaCheckResponse{
				Name:    to_check,
				IsMafia: is_mafia,
			},
		},
	})
}

func (g *Game) sendDeathMessageToGroup(grp map[string]struct{}, killed string) {
	g.sendToGroup(grp, &mafia.Event{
		Type: mafia.EventType_Death,
		Data: &mafia.Event_Death_{
			Death: &mafia.Event_Death{
				Name: killed,
			},
		},
	})
}

func (g *Game) sendDeathMessageToAll(killed string) {
	g.sendDeathMessageToGroup(g.alive_players, killed)
	g.sendDeathMessageToGroup(g.ghosts, killed)
}

func (g *Game) sendGameEndToGrp(grp map[string]struct{}, text string) {
	g.sendToGroup(grp, &mafia.Event{
		Type: mafia.EventType_GameEnd,
		Data: &mafia.Event_GameEnd_{
			GameEnd: &mafia.Event_GameEnd{
				Text: text,
			},
		},
	})
}

func (g *Game) sendGameEndToAll(text string) {
	g.sendGameEndToGrp(g.alive_players, text)
	g.sendGameEndToGrp(g.ghosts, text)
}

/////////////////////////////////////////////
// Game methods
/////////////////////////////////////////////

func randRole() mafia.Role {
	return mafia.Role(rand.Uint32() % 3)
}

func (g *Game) waitEventsFrom(grp map[string]struct{}) {
	for player := range grp {
		g.need_events_from[player] = struct{}{}
	}
}

func (g *Game) assignRoles() {
	max_group_cnt := (g.players_cnt - 1) / 2

	civilians := make([]string, 0)

	for player := range g.alive_players {
		role := randRole()

		if role == mafia.Role_Mafia && len(g.mafia) < int(max_group_cnt) {
			g.mafia[player] = struct{}{}
			g.playersRole[player] = mafia.Role_Mafia
		} else {
			role = mafia.Role_Sheriff
		}

		if role == mafia.Role_Sheriff && len(g.sheriffs) < int(max_group_cnt) {
			g.sheriffs[player] = struct{}{}
			g.playersRole[player] = mafia.Role_Sheriff
		} else {
			role = mafia.Role_Civilian
		}

		if role == mafia.Role_Civilian {
			civilians = append(civilians, player)
		}
	}

	if len(g.mafia) == 0 {
		player := civilians[0]
		civilians = civilians[1:]
		g.mafia[player] = struct{}{}
		g.playersRole[player] = mafia.Role_Mafia
	}

	if len(g.sheriffs) == 0 {
		player := civilians[0]
		g.sheriffs[player] = struct{}{}
		g.playersRole[player] = mafia.Role_Sheriff
	}
}

func (g *Game) Start() {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	g.assignRoles()
	g.sendStartGame()

	// prepare phase
	g.waitEventsFrom(g.alive_players)

	g.sendMsgToAll("prepare phase -> send DoNothing to continue")
	g.actions.Wait()

	// Day -> Night -> Day -> ...
	g.state = Day

	// main loop
	for {
		g.sendMsgToAll("day started")
		g.day()

		if g.isEnded() {
			g.end()
			return
		}

		g.sendMsgToAll("night started")
		g.night()

		if g.isEnded() {
			g.end()
			return
		}
	}
}

func (g *Game) getMostFrequent(arr []string) (string, bool) {
	cnt := make(map[string]int)

	max_cnt := 0
	most_frequent := ""
	is_unique := false

	for _, elem := range arr {
		cnt[elem] += 1

		if cnt[elem] > max_cnt {
			max_cnt = cnt[elem]
			most_frequent = elem
			is_unique = true
		} else if cnt[elem] == max_cnt {
			// not unique variant
			is_unique = false
		}
	}

	return most_frequent, is_unique
}

func (g *Game) day() {
	g.kill_votes = make([]string, 0)

	g.waitEventsFrom(g.alive_players)

	g.sendMsgToGroup(g.alive_players, "vote who is the mafia")
	g.requestVotes(g.alive_players)
	g.actions.Wait()

	to_kill, is_unique := g.getMostFrequent(g.kill_votes)

	if is_unique {
		g.sendMsgToAll(fmt.Sprintf("%s was killed", to_kill))
		g.kill(to_kill)
	} else {
		g.sendMsgToAll("people did not come to a consensus")
	}

	g.changeState()
}

func (g *Game) night() {
	g.kill_votes = make([]string, 0)
	g.check_votes = make([]string, 0)

	g.waitEventsFrom(g.mafia)
	g.waitEventsFrom(g.sheriffs)

	g.sendMsgToGroup(g.mafia, "vote who to kill")
	g.requestVotes(g.mafia)
	g.sendMsgToGroup(g.sheriffs, "vote who to check")
	g.requestVotes(g.sheriffs)
	g.actions.Wait()

	to_kill, is_unique := g.getMostFrequent(g.kill_votes)
	if is_unique {
		g.sendMsgToAll(fmt.Sprintf("mafia killed %s", to_kill))
		g.kill(to_kill)
	} else {
		g.sendMsgToAll("mafia did not come to a consensus")
	}

	to_check, is_unique := g.getMostFrequent(g.check_votes)
	if is_unique {
		g.sendCheckResult(to_check)
	} else {
		g.sendMsgToAll("sheriffs did not come to a consensus")
	}

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
		delete(g.sheriffs, player)
	}

	g.sendDeathMessageToAll(player)
}

func (g *Game) isEnded() bool {
	return g.state == WinMafia || g.state == WinSheriffs
}

func (g *Game) changeState() {
	mafia_cnt := len(g.mafia)
	sheriffs_cnt := len(g.sheriffs)
	civilian_cnt := len(g.alive_players) - mafia_cnt - sheriffs_cnt

	if mafia_cnt == 0 {
		g.state = WinSheriffs
	} else if mafia_cnt >= civilian_cnt+sheriffs_cnt {
		g.state = WinMafia
	} else {
		g.state = (g.state + 1) % 2
	}
}

func (g *Game) end() {
	if g.state == WinMafia {
		g.sendGameEndToAll("game ended, mafia win")
	} else {
		g.sendGameEndToAll("game ended, sheriffs win")
	}
}

func NewGame(player_cnt uint32) *Game {
	game := &Game{
		players_cnt: player_cnt,
		playersInfo: make(map[string]Player),
		playersRole: make(map[string]mafia.Role),

		alive_players: make(map[string]struct{}),
		ghosts:        make(map[string]struct{}),
		mafia:         make(map[string]struct{}),
		sheriffs:      make(map[string]struct{}),

		state:            PrepareState,
		need_events_from: make(map[string]struct{}),
		kill_votes:       make([]string, 0),
		check_votes:      make([]string, 0),
	}

	game.actions = sync.Cond{
		L: &game.mtx,
	}

	return game
}
