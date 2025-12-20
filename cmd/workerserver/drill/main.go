package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gustavo-melo-dev/idle-factory/internal/worker"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MiningDrill struct {
	*worker.WorkerStats
}

func (m MiningDrill) Work(ctx context.Context) error {
	defer m.Ticker.Stop()
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("failed to dial RabbitMQ: %v", err)
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

	for {
		select {
		case <-m.Ticker.C:
			output := fmt.Sprintf(
				"%s %s just mined %.2f %s",
				m.WorkerType.String(),
				m.ID,
				m.OutputPerTick,
				m.OutputType.String(),
			)

			ch.PublishWithContext(
				ctx,
				"work",  // exchange
				"drill", // routing key
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        []byte(output),
				},
			)
		case <-ctx.Done():
			return nil
		}
	}
}

func (m MiningDrill) Stats() *worker.WorkerStats {
	return m.WorkerStats
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
	drill := New(time.Second)

	blocker := make(chan struct{})

	ctx := context.Background()
	go drill.Work(ctx)

	<-blocker
}
