package chat

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func (c *rabbitConnection) GetMsgsChat() (<-chan amqp.Delivery, error) {
	msgs, err := c.ch.Consume(
		c.q.Name, // queue
		"",       // consumer
		true,     // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
