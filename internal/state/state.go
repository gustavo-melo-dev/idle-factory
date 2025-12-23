package state

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/gustavo-melo-dev/idle-factory/internal/machine"
	"github.com/gustavo-melo-dev/idle-factory/internal/resource"
)

type State struct {
	Machines  map[machine.MachineType][]string
	Resources map[resource.ResourceType]float64
	mu        sync.RWMutex // state is updated concurrently
}

func (s *State) String() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result strings.Builder
	// ANSI escape codes explanation:
	//   - \033 is the ESC character (ASCII 27). It signals the start of an escape sequence.
	//   - [2J is the command to clear the entire screen.
	//   - [H moves the cursor to the home position (top-left corner, row 1, column 1).
	// Combining these codes effectively clears the terminal screen and resets the cursor position.
	result.WriteString("\033[2J\033[H")
	result.WriteString("Current State:\n")

	// maps in Go do not guarantee order, so we sort the resource types for consistent output
	var resourceTypes []int
	for resType := range s.Resources {
		resourceTypes = append(resourceTypes, int(resType))
	}
	sort.Ints(resourceTypes)

	result.WriteString("Resources:\n")
	for resourceType := range resourceTypes {
		// convert from int back to ResourceType
		resourceType := resource.ResourceType(resourceType)

		result.WriteString("- ")
		result.WriteString(resourceType.String())
		result.WriteString(": ")
		result.WriteString(fmt.Sprintf("%.2f", s.Resources[resourceType]))
		result.WriteString("\n")
	}

	var machineTypes []int
	for mType := range s.Machines {
		machineTypes = append(machineTypes, int(mType))
	}
	sort.Ints(machineTypes)

	result.WriteString("Machines:\n")
	for mType, machines := range s.Machines {
		// convert from int back to MachineType
		mType := machine.MachineType(mType)

		result.WriteString("- ")
		result.WriteString(mType.String())
		result.WriteString(": ")
		result.WriteString(fmt.Sprintf("%d units\n", len(machines)))
	}
	return result.String()
}

func New() *State {
	return &State{
		Machines:  make(map[machine.MachineType][]string),
		Resources: make(map[resource.ResourceType]float64),
	}
}
func (s *State) UpdateMachine(mType machine.MachineType, id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Machines[mType] = append(s.Machines[mType], id)
}

func (s *State) UpdateResource(resType resource.ResourceType, amount float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Resources[resType] += amount
}
