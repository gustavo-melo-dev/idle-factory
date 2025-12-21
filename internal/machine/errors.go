package machine

import (
	"fmt"

	"github.com/gustavo-melo-dev/idle-factory/internal/resource"
)

// MachineShouldNotProduceResourceError is returned when a machine is assigned a resource type it cannot produce
type MachineShouldNotProduceResourceError struct {
	MachineType
	resource.ResourceType
}

func (e MachineShouldNotProduceResourceError) Error() string {
	return fmt.Sprintf("Machine of type %s should not produce %s", e.MachineType, e.ResourceType)
}
