package types

// 
const (
	MaximumDepth uint = 1024
)

type Byte32 [32]byte


type Stack struct {
	Data []byte
}
type Memory struct {
	Data []Byte32
}

type Storage struct {
	Data map[Byte32]Byte32
}
