package types

import (
	// "fmt"
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