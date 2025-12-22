package main

import (
	"log"

	"github.com/gustavo-melo-dev/idle-factory/internal/broker"
	"github.com/gustavo-melo-dev/idle-factory/internal/machine"
	"github.com/gustavo-melo-dev/idle-factory/internal/state"
)

func main() {
	state := state.New()
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
			machineMsg := machine.DecodeMessage(msg)
			state.UpdateResource(machineMsg.ResultResourceType, machineMsg.Amount)
			log.Print(state)
			msg.Ack(false)
		}
	}()

	log.Print(state)
	<-blocker
}
