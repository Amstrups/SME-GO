package types

const (
	PUBLIC = "PUBLIC"
	SECRET = "SECRET"
)

type T int

type Chan chan T
type ReadChan <-chan T
type WriteChan chan<- T

type ReadChannels map[string]ReadChan
type WriteChannels map[string]WriteChan

// utils
func MakeReadOnly(ch Chan) <-chan T {
	return ch
}

func MakeWriteOnly(ch Chan) chan<- T {
	return ch
}
