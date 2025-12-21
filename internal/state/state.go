package state

import (
	"fmt"
	"strings"

	"github.com/gustavo-melo-dev/idle-factory/internal/machine"
	"github.com/gustavo-melo-dev/idle-factory/internal/resource"
)

type State struct {
	Machines  map[machine.MachineType][]string
	Resources map[resource.ResourceType]float64
}

func (s *State) String() string {
	var result strings.Builder
	// ANSI escape codes explanation:
	//   - \033 is the ESC character (ASCII 27). It signals the start of an escape sequence.
	//   - [2J is the command to clear the entire screen.
	//   - [H moves the cursor to the home position (top-left corner, row 1, column 1).
	// Combining these codes effectively clears the terminal screen and resets the cursor position.
	result.WriteString("\033[2J\033[H")
	result.WriteString("Current State:\n")
	result.WriteString("Resources:\n")
	for resType, amount := range s.Resources {
		result.WriteString("- ")
		result.WriteString(resType.String())
		result.WriteString(": ")
		result.WriteString(fmt.Sprintf("%.2f", amount))
		result.WriteString("\n")
	}
	result.WriteString("Machines:\n")
	for mType, machines := range s.Machines {
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
	s.Machines[mType] = append(s.Machines[mType], id)
}

func (s *State) UpdateResource(resType resource.ResourceType, amount float64) {
	s.Resources[resType] += amount
}
