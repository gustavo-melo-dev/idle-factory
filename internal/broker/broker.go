// Package broker provides functions to interact with RabbitMQ message broker more easily.
package broker

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Connect establishes a connection to the RabbitMQ server.
func Connect() *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("failed to dial RabbitMQ: %v", err)
	}

	return conn
}

// GetChannel opens a channel on the given RabbitMQ connection.
func GetChannel(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel: %v", err)
	}
	return ch
}

// DeclareTopicQueue declares a topic exchange and a queue, then binds them together using a routing key and returns the queue.
func DeclareTopicQueue(ch *amqp.Channel, exchange string, routingKey string, queueName string) amqp.Queue {
	err := ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-delete
		false,    // internal
		false,    // no-wait
		nil,      // args
	)
	if err != nil {
		log.Fatalf("failed to declare exchange %s: %v", queueName, err)
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("failed to declare queue %s: %v", queueName, err)
	}

	err = ch.QueueBind(
		q.Name,     // queue
		routingKey, // routing key
		exchange,   // exchange
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Fatalf(
			"failed to bind queue %s to exchange %s: %v", queueName, exchange, err,
		)
	}

	return q
}

// PublishMessage publishes a message to the specified exchange with the given routing key.
func PublishMessage(ctx context.Context, ch *amqp.Channel, exchange string, routingKey string, body []byte) {
	err := ch.PublishWithContext(
		ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Fatalf("failed to publish message to exchange %s: %v", exchange, err)
	}
}

// ConsumeMessages sets up a consumer on the specified queue and returns a channel to receive messages.
func ConsumeMessages(ch *amqp.Channel, queueName string) <-chan amqp.Delivery {
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack // we acknowledge manually after processing
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("failed to start consuming messages: %v", err)
	}

	return msgs
}
