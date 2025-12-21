// Package machine implements a worker that produces resources.
package machine

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gustavo-melo-dev/idle-factory/internal/broker"
	"github.com/gustavo-melo-dev/idle-factory/internal/resource"
	"github.com/gustavo-melo-dev/idle-factory/internal/worker"
	amqp "github.com/rabbitmq/amqp091-go"
)

// MachineType represents the type of machine
type MachineType int

const (
	MiningDrill MachineType = iota
	Furnace
	ScienceLab
)

func (t MachineType) String() string {
	switch t {
	case MiningDrill:
		return "Mining Drill"
	case Furnace:
		return "Furnace"
	case ScienceLab:
		return "Science Lab"
	default:
		return "Unknown Worker Type"
	}
}

// CanProduce reports whether that machine type can produce the given resource
func (t MachineType) CanProduce(res resource.ResourceType) bool {
	switch res {
	case resource.IronOre:
		return t == MiningDrill
	case resource.IronPlate:
		return t == Furnace
	case resource.RedScience:
		return t == ScienceLab
	default:
		return false
	}
}

// Machine represents a worker that produces resources.
type Machine struct {
	*worker.WorkerStats
	Type               MachineType
	CyclesComplete     int
	TotalOutput        float64
	ResultResourceType resource.ResourceType
	ResourcePerTick    float64
	ActionDescription  string // "mined", "smelted", "researched", etc.
	Ticker             *time.Ticker
}

// New returns a new Machine instance.
func New(
	machineType MachineType,
	resultResourceType resource.ResourceType,
	resourcePerTick float64,
	actionItPerforms string,
	interval time.Duration,
) (*Machine, error) {
	if !machineType.CanProduce(resultResourceType) {
		return nil, MachineShouldNotProduceResourceError{
			MachineType:  machineType,
			ResourceType: resultResourceType,
		}
	}
	ticker := time.NewTicker(interval)
	return &Machine{
		worker.NewStats(),
		machineType,
		0,
		0,
		resultResourceType,
		resourcePerTick,
		actionItPerforms,
		ticker,
	}, nil
}

// MachineMessage represents the message structure that machines will send to the broker
type MachineMessage struct {
	*worker.WorkMessage
	MachineType        MachineType           `json:"machine_type"`
	ResultResourceType resource.ResourceType `json:"result_resource_type"`
	Amount             float64               `json:"amount"`
}

// EncodeMessage encodes a MachineMessage into a JSON byte slice.
func EncodeMessage(msg *MachineMessage) []byte {
	mm, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("failed to marshal machine message: %v", err)
		return nil
	}
	return mm
}

// DecodeMessage decodes the body of an amqp.Delivery into a MachineMessage.
func DecodeMessage(d amqp.Delivery) *MachineMessage {
	var msg MachineMessage
	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		log.Fatalf("failed to unmarshal delivery body: %v", err)
		return nil
	}
	return &msg
}

func (m Machine) Message() *MachineMessage {
	return &MachineMessage{
		&worker.WorkMessage{
			WorkerID: m.ID,
			Message: fmt.Sprintf(
				"%s %s %s %.2f %s",
				m.Type.String(),
				m.ID,
				m.ActionDescription,
				m.ResourcePerTick,
				m.ResultResourceType.String(),
			),
		},
		m.Type,
		m.ResultResourceType,
		m.ResourcePerTick,
	}
}

func (m Machine) Work(ctx context.Context, exchange, routingKey, queueName string) error {
	defer m.Ticker.Stop()

	conn := broker.Connect()
	defer conn.Close()

	ch := broker.GetChannel(conn)
	defer ch.Close()

	broker.DeclareTopicQueue(ch, exchange, routingKey, queueName)

	for {
		select {
		case <-m.Ticker.C:
			body := EncodeMessage(m.Message())
			broker.PublishMessage(ctx, ch, exchange, routingKey, body)
			m.CyclesComplete++
			m.TotalOutput += m.ResourcePerTick
		case <-ctx.Done():
			return nil
		}
	}
}
