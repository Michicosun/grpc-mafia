package chat

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (c *rabbitConnection) SendMessage(from string, text string, to int32) error {
	msg := Message{
		User: from,
		Text: text,
	}

	raw_msg, err := json.Marshal(&msg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.ch.PublishWithContext(ctx,
		c.exchs[to], // exchange
		"",          // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        raw_msg,
		}); err != nil {
		return err
	}

	return nil
}
