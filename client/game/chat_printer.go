package game

import (
	"encoding/json"
	"grpc-mafia/client/chat"
)

var ChatPrinter = &chatPrinter{}

type chatPrinter struct {
	stop_ch chan struct{}
}

func (cp *chatPrinter) Start() {
	go func() {
		cp.stop_ch = make(chan struct{}, 1)

		msgs, err := chat.RabbitConnection.GetMsgsChat()
		if err != nil {
			Session.StopWithError(err)
			return
		}

		for {
			select {
			case raw_msg := <-msgs:
				msg := chat.Message{}

				if err := json.Unmarshal(raw_msg.Body, &msg); err != nil {
					Session.StopWithError(err)
					return
				}

				PrintLine(msg.User, msg.Text, Session.Interactor)

			case <-cp.stop_ch:
				return
			}
		}
	}()
}

func (cp *chatPrinter) Stop() {
	cp.stop_ch <- struct{}{}
}
