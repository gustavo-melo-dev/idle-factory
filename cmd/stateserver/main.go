package main

import (
	"log"

	"github.com/gustavo-melo-dev/idle-factory/internal/broker"
	"github.com/gustavo-melo-dev/idle-factory/internal/machine"
	"github.com/gustavo-melo-dev/idle-factory/internal/state"
)

func main() {
	state := state.New()
	exchange := "work"
	routingKeys := []string{"drill", "furnace", "lab"}

	for _, rk := range routingKeys {
		listenToMachineWorking(state, exchange, rk)
	}

	blocker := make(chan struct{})
	<-blocker
}

func listenToMachineWorking(state *state.State, exchange string, routingKey string) {
	conn := broker.Connect()
	ch := broker.GetChannel(conn)

	queueName := exchange + "." + routingKey
	q := broker.DeclareTopicQueue(ch, exchange, routingKey, queueName)

	go func() {
		msgs := broker.ConsumeMessages(ch, q.Name)
		for msg := range msgs {
			decoded := machine.DecodeMessage(msg)
			state.UpdateResource(decoded.ResultResourceType, decoded.Amount)
			msg.Ack(false)
			log.Print(state)
		}
	}()
}
