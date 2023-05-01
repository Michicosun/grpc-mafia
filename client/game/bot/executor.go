package game

import (
	"sync"
)

var BotInteractor = &botInteractor{}

type botInteractor struct {
	event sync.Cond
}

func (b *botInteractor) Run() {

	// TODO: check Game and make move

	b.event.Wait()
}

func (b *botInteractor) Signal() {
	b.event.Signal()
}
