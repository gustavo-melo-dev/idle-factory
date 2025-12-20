package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel: %v", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"work",  // name
		"topic", // type
		true,    // durable
		false,   // auto-delete
		false,   // internal
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		log.Fatalf("failed to declare exchange work: %v", err)
	}

	q, err := ch.QueueDeclare(
		"work.drill",
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("failed to declare queue work.drill: %v", err)
	}

	err = ch.QueueBind(
		q.Name,  // queue
		"drill", // routing key
		"work",  // exchange
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		log.Fatalf(
			"failed to bind queue work.drill to exchange work: %v", err,
		)
	}

	msgs, err := ch.Consume(
		"work.drill", // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		log.Fatalf("failed to start consuming messages: %v", err)
	}

	blocker := make(chan struct{})

	go func() {
		for d := range msgs {
			log.Print(string(d.Body))
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-blocker
}
