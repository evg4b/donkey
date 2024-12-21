package store

type EventType int

const (
	FileProcessing EventType = iota
	FileProcessed
	MemoryCleared
)

type Event struct {
	Type           EventType
	InputFileName  string
	OutputFileName string
}

func (s *Event) HasDifferentOutput() bool {
	return s.InputFileName != s.OutputFileName
}
