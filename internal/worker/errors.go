package worker

import "fmt"

// IncompatibleOutputToWorkerTypeError is returned when a worker is assigned an output type it cannot produce
type IncompatibleOutputToWorkerTypeError struct {
	WorkerType
	OutputType
}

func (e IncompatibleOutputToWorkerTypeError) Error() string {
	return fmt.Sprintf("Worker of type %s should not produce %s", e.WorkerType, e.OutputType)
}
