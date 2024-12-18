package store

type EventType int

const (
	FileProcessing EventType = iota
	FileProcessed
	MemoryCleared
)

type Event struct {
	Type     EventType
	FileName string
}
