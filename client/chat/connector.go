package chat

import (
	"fmt"
	mafia "grpc-mafia/server/proto"
	"grpc-mafia/util"

	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitConnection = &rabbitConnection{}

type Message struct {
	Text string `json:"text"`
	User string `json:"user"`
}

type RabbitCredentials struct {
	User string
	Pass string
	Host string
	Port string
}

type rabbitConnection struct {
	creds RabbitCredentials

	exchs []string
	conn  *amqp.Connection
	ch    *amqp.Channel
	q     *amqp.Queue
}

func (c *rabbitConnection) createConnectionString() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", c.creds.User, c.creds.Pass, c.creds.Host, c.creds.Port)
}

func (c *rabbitConnection) createExchanges(session_id string, role mafia.Role) error {
	c.exchs[mafia.Role_Civilian] = session_id
	if err := c.ch.ExchangeDeclare(
		session_id, // name
		"fanout",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	); err != nil {
		return err
	}

	if role != mafia.Role_Civilian {
		c.exchs[role] = util.CreateExchangeName(session_id, role.String())
		if err := c.ch.ExchangeDeclare(
			util.CreateExchangeName(session_id, role.String()), // name
			"fanout", // type
			true,     // durable
			false,    // auto-deleted
			false,    // internal
			false,    // no-wait
			nil,      // arguments
		); err != nil {
			return err
		}
	}

	return nil
}

func (c *rabbitConnection) createQueue(session_id string, name string) error {
	q, err := c.ch.QueueDeclare(
		util.CreateQueueName(session_id, name), // name
		false,                                  // durable
		false,                                  // delete when unused
		true,                                   // exclusive
		false,                                  // no-wait
		nil,                                    // arguments
	)
	if err != nil {
		return err
	}

	c.q = &q

	return nil
}

func (c *rabbitConnection) bindQueue(session_id string, role mafia.Role) error {
	if err := c.ch.QueueBind(
		c.q.Name,   // queue name
		"",         // routing key
		session_id, // exchange
		false,
		nil,
	); err != nil {
		return err
	}

	if role != mafia.Role_Civilian {
		if err := c.ch.QueueBind(
			c.q.Name, // queue name
			"",       // routing key
			util.CreateExchangeName(session_id, role.String()), // exchange
			false,
			nil,
		); err != nil {
			return err
		}
	}

	return nil
}

func (c *rabbitConnection) Establish(session_id string, role mafia.Role, name string) error {
	c.exchs = make([]string, 3)

	conn, err := amqp.Dial(c.createConnectionString())
	if err != nil {
		return err
	}
	c.conn = conn

	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	c.ch = ch

	if err := c.createExchanges(session_id, role); err != nil {
		return err
	}

	if err := c.createQueue(session_id, name); err != nil {
		return err
	}

	return c.bindQueue(session_id, role)
}

func (c *rabbitConnection) Close() error {
	if _, err := c.ch.QueueDelete(c.q.Name, false, false, false); err != nil {
		return err
	}

	c.ch.Close()
	c.conn.Close()

	return nil
}

func (c *rabbitConnection) Init(creds RabbitCredentials) {
	c.creds = creds
}
