package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitCredentials struct {
	User string
	Pass string
	Host string
	Port string
}

type TaskQueue struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    *amqp.Queue
}

func createConnectionString(creds *RabbitCredentials) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", creds.User, creds.Pass, creds.Host, creds.Port)
}

func (tq *TaskQueue) GetTaskChan() <-chan amqp.Delivery {
	msgs, err := tq.ch.Consume(
		tq.q.Name, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		panic(err)
	}

	return msgs
}

func (tq *TaskQueue) SubmitTask(task interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	bytes, err := json.Marshal(task)
	if err != nil {
		return err
	}

	if err := tq.ch.PublishWithContext(ctx,
		"",        // exchange
		tq.q.Name, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        bytes,
		}); err != nil {
		return err
	}

	return nil
}

func NewTaskQueue(creds RabbitCredentials) *TaskQueue {
	tq := TaskQueue{}

	conn, err := amqp.Dial(createConnectionString(&creds))
	if err != nil {
		panic(err) // can't work without task queue
	}
	tq.conn = conn

	ch, err := tq.conn.Channel()
	if err != nil {
		panic(err) // can't work without task queue
	}
	tq.ch = ch

	q, err := tq.ch.QueueDeclare(
		"render-tasks", // name
		false,          // durable
		false,          // delete when unused
		true,           // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		panic(err) // can't work without task queue
	}

	tq.q = &q

	return &tq
}
