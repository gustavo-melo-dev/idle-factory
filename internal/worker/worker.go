// Package worker defines the contract for an entity to perform some work.
package worker

import (
	"context"
	"crypto/rand"
)

// WorkerStats represents the information that every worker should keep.
type WorkerStats struct {
	ID string
}

func NewStats() *WorkerStats {
	return &WorkerStats{
		ID: rand.Text(),
	}
}

// WorkMessage represents the most basic message structure that workers need to send to the broker.
type WorkMessage struct {
	WorkerID string `json:"worker_id"`
	Message  string `json:"message"`
}

// Worker represents the minimum interface that any worker must implement to be able to do some work.
type Worker interface {
	Work(ctx context.Context, exchange, routingKey, queueName string) error
}
