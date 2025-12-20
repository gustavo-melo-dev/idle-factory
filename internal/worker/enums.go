package worker

// OutputType represents the type of output a worker can produce
type OutputType int

const (
	IronOre OutputType = iota
	IronPlate
	RedScience
)

func (t OutputType) String() string {
	switch t {
	case IronOre:
		return "Iron Ore"
	case IronPlate:
		return "Iron Plate"
	case RedScience:
		return "Red Science"
	default:
		return "Unknown Worker Output Type"
	}
}

// WorkerType represents the type of worker
type WorkerType int

const (
	MiningDrill WorkerType = iota
	Furnace
	ScienceLab
)

func (t WorkerType) String() string {
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

// CanProduce reports whether the worker type can produce the given output type
func (t WorkerType) CanProduce(out OutputType) bool {
	switch out {
	case IronOre:
		return t == MiningDrill
	case IronPlate:
		return t == Furnace
	case RedScience:
		return t == ScienceLab
	default:
		return false
	}
}
