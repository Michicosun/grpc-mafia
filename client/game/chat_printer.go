package game

import (
	"encoding/json"
	"grpc-mafia/client/chat"

	amqp "github.com/rabbitmq/amqp091-go"
)

var ChatPrinter = &chatPrinter{}

type chatPrinter struct {
	msgs    <-chan amqp.Delivery
	stop_ch chan struct{}
}

func (cp *chatPrinter) Start() {
	go func() {
		cp.msgs = chat.RabbitConnection.GetMsgsChat()
		cp.stop_ch = make(chan struct{}, 1)

		for {
			select {
			case raw_msg := <-cp.msgs:
				msg := chat.Message{}

				if err := json.Unmarshal(raw_msg.Body, &msg); err != nil {
					panic(err)
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
