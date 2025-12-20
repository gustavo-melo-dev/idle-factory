package main

import (
	"log"

	"github.com/gustavo-melo-dev/idle-factory/internal/broker"
)

func main() {
	exchange, routingKey, queueName := "work", "drill", "work.drill"
	conn := broker.Connect()
	defer conn.Close()

	ch := broker.GetChannel(conn)
	defer ch.Close()

	q := broker.DeclareTopicQueue(ch, exchange, routingKey, queueName)

	msgs := broker.ConsumeMessages(ch, q.Name)

	blocker := make(chan struct{})

	go func() {
		for msg := range msgs {
			workMsg := broker.DecodeDeliveryMessage(msg)
			log.Print(workMsg.Message)
			msg.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-blocker
}
