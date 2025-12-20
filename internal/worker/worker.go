// Package worker defines the base interface that all game structures must implement along with help structs.
// It provides the common contract that each building that does some work must comply
package worker

import (
	"context"
	"crypto/rand"
	"time"
)

// WorkerStats represents the information that every worker should keep
type WorkerStats struct {
	ID             string
	WorkerType     WorkerType
	CyclesComplete int
	TotalOutput    float32
	OutputType     OutputType
	OutputPerTick  float32
	Ticker         *time.Ticker
}

func NewWorkerStats(workerType WorkerType, outputType OutputType, outputPerTick float32, interval time.Duration) (*WorkerStats, error) {
	ticker := time.NewTicker(interval)

	if !workerType.CanProduce(outputType) {
		return nil, &IncompatibleOutputToWorkerTypeError{WorkerType: workerType, OutputType: outputType}
	}

	return &WorkerStats{
		ID:             rand.Text(),
		WorkerType:     workerType,
		CyclesComplete: 0,
		TotalOutput:    0,
		OutputType:     outputType,
		OutputPerTick:  outputPerTick,
		Ticker:         ticker,
	}, nil
}

// Worker represents the minimum interface that any worker must implement to be able to do some work
type Worker interface {
	Work(ctx context.Context) error
	Stats() *WorkerStats
}
