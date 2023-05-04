package game

import (
	"fmt"
	"grpc-mafia/client/game"
	"grpc-mafia/client/grpc"
	"grpc-mafia/util"
	"math/rand"
	"time"
)

var BotInteractor = &botInteractor{}

type botInteractor struct {
	event util.Event
}

func (b *botInteractor) Run() {
	for {
		b.Step()
		b.event.Wait()
	}
}

func (b *botInteractor) Step() {
	time.Sleep(time.Millisecond * time.Duration(100+rand.Intn(100)))
	switch game.Session.GetState() {
	case game.Undefined:
		game.Session.Start(b.generateRandomName())

	case game.PrepareState:
		game.Session.ChangeState(game.Waiting, false)
		grpc.Connection.SendDoNothing(game.Session.Name)

	case game.NeedVote:
		if game.Session.MafiaCheck && rand.Intn(2) == 0 {
			grpc.Connection.SendPublishRequest(game.Session.MafiaName)
		}
		game.Session.ChangeState(game.Waiting, false)
		grpc.Connection.SendVote(game.Session.Name, b.getRandomAlive())
	}
}

func (b *botInteractor) generateRandomName() string {
	return fmt.Sprintf("bot_%d", rand.Int())
}

func (b *botInteractor) getRandomAlive() string {
	names := make([]string, 0)

	for name := range game.Session.AlivePlayers {
		if name != game.Session.Name {
			names = append(names, name)
		}
	}

	return names[rand.Intn(len(names))]
}

func (b *botInteractor) Signal() {
	b.event.Signal()
}

func (b *botInteractor) GetPrefixString() string {
	return ""
}

func (b *botInteractor) GetCurBuf() string {
	return ""
}

func MakeBotInteractor() *botInteractor {
	return &botInteractor{
		event: util.CreateEvent(),
	}
}
