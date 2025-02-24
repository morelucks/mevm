package types

import (
	"fmt"
)
//..... stack contructor ...//
func NewStack() *Stack {
	return &Stack{
		Data: make([]byte, 0),
	}
}

//..... memory contructor ...//
func NewMemory() *Memory {
	return &Memory{make([]Byte32, 0)}
}

//..... storage contructor ...//
func NewStorage() *Storage {
	return &Storage{make(map[Byte32]Byte32, 0)}
}


func (m *Memory) Store(offset byte, value Byte32) {
	index := int(offset) / 32

	if index >= len(m.Data) {
		newMemoryData := make([]Byte32, index+1)
		copy(newMemoryData, m.Data)
		m.Data = newMemoryData
	}

	m.Data[index] = value
}

func (m *Memory) Load(offset byte) Byte32 {
	index := int(offset) / 32

	if index > len(m.Data)-1 {
		panic("invalid memory location")
	}
	return m.Data[index]
}
func (s *Storage) Store(key Byte32, value Byte32) {
	s.Data[key] = value
}

func (s *Storage) Load(key Byte32) (Byte32, error) {
	data, err := s.Data[key]
	if !err {
		return Byte32{}, fmt.Errorf("Error: invalid storage key")
	}

	return data, nil
}

//.... push method ...//

func (s *Stack) Push(data byte) {
	if len(s.Data) >= int(MaximumDepth) {
		panic("stack overflow")
	}

	s.Data = append(s.Data, data)

}
//.... pop method ...//

func (s *Stack) Pop() byte {
	lastIndex := len(s.Data) - 1
	lastItem := s.Data[lastIndex]
	s.Data = append(s.Data[:lastIndex], s.Data[lastIndex+1:]...)

	return lastItem
}