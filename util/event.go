package util

type Event struct {
	ch chan struct{}
}

func (e *Event) Signal() {
	e.ch <- struct{}{}
}

func (e *Event) Wait() {
	<-e.ch
}

func CreateEvent() Event {
	return Event{
		ch: make(chan struct{}, 1),
	}
}
