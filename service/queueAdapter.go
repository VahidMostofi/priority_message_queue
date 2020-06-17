package main

import (
	"log"

	"github.com/streadway/amqp"
)

type QueueAdapter struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queues     map[string]amqp.Queue
	Done       chan bool
}

func NewQueueAdapter(amqpURL string) *QueueAdapter {
	queueAdapter := &QueueAdapter{}
	// Connect to rabbitMQ server
	conn, err := amqp.Dial(amqpURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	queueAdapter.Connection = conn

	// Open a chanel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	queueAdapter.Channel = ch

	queueAdapter.Queues = make(map[string]amqp.Queue)

	return queueAdapter
}

func (q *QueueAdapter) CreateQueue(queueName string) {
	var max uint8 = 9
	args := map[string]interface{}{"x-max-priority": max}
	// A queue to send to
	queue, err := q.Channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		args,      // arguments
	)
	failOnError(err, "Failed to declare a queue")
	q.Queues[queueName] = queue
}

func (q *QueueAdapter) Publish(body []byte, priority uint8, queueName string) {
	err := q.Channel.Publish("", q.Queues[queueName].Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
		Priority:    priority,
	})
	failOnError(err, "failed to publish to "+queueName)
}

func (q *QueueAdapter) Consume(queueName string, handler chan<- []byte) {
	msgs, err := q.Channel.Consume(
		q.Queues[queueName].Name, // queue
		"",                       // consumer
		true,                     // auto-ack
		false,                    // exclusive
		false,                    // no-local
		false,                    // no-wait
		nil,                      // args
	)
	failOnError(err, "Failed to register a consumer")

	q.Done = make(chan bool)

	go func() {
		for d := range msgs {
			handler <- d.Body
		}
	}()

	<-q.Done

}

func (q *QueueAdapter) Close() {
	if q.Done != nil {
		q.Done <- true
	}
	q.Channel.Close()
	q.Connection.Close()
}

// helper function to check the return value for each amqp call
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
