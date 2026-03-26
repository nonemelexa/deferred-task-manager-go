package queue

import (
	"encoding/json"
	"log"
	"os"

	"task-scheduler/internal/domain"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

// NewRabbitMQ establishes a connection to RabbitMQ, declares a queue, and returns a RabbitMQ struct for publishing and consuming tasks.
func NewRabbitMQ() *RabbitMQ {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	q, err := ch.QueueDeclare(
		"tasks",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	return &RabbitMQ{conn, ch, q}
}

// Publish sends a task to the RabbitMQ queue by marshaling it to JSON and publishing it to the declared queue.
func (r *RabbitMQ) Publish(task domain.Task) error {
	body, err := json.Marshal(task)
	if err != nil {
		return err
	}

	err = r.channel.Publish(
		"",
		r.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	return err
}

// Consume returns a channel that receives tasks from the RabbitMQ queue by consuming messages, unmarshaling them from JSON, and sending them to the output channel.
func (r *RabbitMQ) Consume() <-chan domain.Task {
	msgs, err := r.channel.Consume(
		r.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	out := make(chan domain.Task)

	go func() {
		for msg := range msgs {
			var task domain.Task
			json.Unmarshal(msg.Body, &task)
			out <- task
		}
	}()

	return out
}
