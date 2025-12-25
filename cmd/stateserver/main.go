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

	listenToMachineWorking(state, exchange)

	blocker := make(chan struct{})
	<-blocker
}

func listenToMachineWorking(state *state.State, exchange string) {
	conn := broker.Connect()
	ch := broker.GetChannel(conn)

	broker.DeclareTopicExchange(ch, exchange)
	q := broker.DeclareQueue(ch, exchange, "#", "work.queue")

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
