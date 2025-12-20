package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gustavo-melo-dev/idle-factory/internal/broker"
	"github.com/gustavo-melo-dev/idle-factory/internal/worker"
)

type MiningDrill struct {
	*worker.WorkerStats
}

func (m MiningDrill) Work(ctx context.Context) error {
	exchange, routingKey, queueName := "work", "drill", "work.drill"

	defer m.Ticker.Stop()

	conn := broker.Connect()
	defer conn.Close()

	ch := broker.GetChannel(conn)
	defer ch.Close()

	broker.DeclareTopicQueue(ch, exchange, routingKey, queueName)

	for {
		select {
		case <-m.Ticker.C:
			body := m.Message()
			broker.PublishMessage(ctx, ch, exchange, routingKey, body)
			m.CyclesComplete++
			m.TotalOutput += m.OutputPerTick
		case <-ctx.Done():
			return nil
		}
	}
}

func (m MiningDrill) Stats() *worker.WorkerStats {
	return m.WorkerStats
}

func (m MiningDrill) Message() *worker.WorkMessage {
	return &worker.WorkMessage{
		WorkerID:   m.ID,
		WorkerType: m.WorkerType,
		OutputType: m.OutputType,
		Amount:     m.OutputPerTick,
		Message:    fmt.Sprintf("%s %s just mined %.2f %s", m.WorkerType.String(), m.ID, m.OutputPerTick, m.OutputType.String()),
	}
}

func New(interval time.Duration) *MiningDrill {
	stats, err := worker.NewWorkerStats(worker.MiningDrill, worker.IronOre, 0.1, interval)
	if err != nil {
		log.Fatalf("failed to create worker stats: %v", err)
	}

	return &MiningDrill{
		stats,
	}
}

func main() {
	var drill worker.Worker
	drill = New(time.Second)

	blocker := make(chan struct{})

	ctx := context.Background()
	go drill.Work(ctx)

	<-blocker
}
