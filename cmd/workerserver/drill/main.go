package main

import (
	"context"
	"log"
	"time"

	"github.com/gustavo-melo-dev/idle-factory/internal/machine"
	"github.com/gustavo-melo-dev/idle-factory/internal/resource"
	"github.com/gustavo-melo-dev/idle-factory/internal/worker"
)

func main() {
	var drill worker.Worker
	drill, err := machine.New(machine.MiningDrill, resource.IronOre, 0.1, "mined", time.Second*1)
	if err != nil {
		log.Fatalf("Failed to create a mining drill")
	}

	blocker := make(chan struct{})

	ctx := context.Background()
	go drill.Work(ctx, "work", "drill", "work.drill")

	<-blocker
}
