// Package resource defines the different types of resources that machines can produce.
package resource

// ResourceType represents the type of resource a machine can produce
type ResourceType int

const (
	IronOre ResourceType = iota
	IronPlate
	RedScience
)

func (t ResourceType) String() string {
	switch t {
	case IronOre:
		return "Iron Ore"
	case IronPlate:
		return "Iron Plate"
	case RedScience:
		return "Red Science"
	default:
		return "Unknown Resource Type"
	}
}
