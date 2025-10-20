package types

import (
	"fmt"
	"math/big"
)

// Constructor functions
func NewStack() *Stack {
	return &Stack{
		Data: make([]Word, 0),
	}
}

func NewMemory() *Memory {
	return &Memory{make([]Byte32, 0)}
}

func NewStorage() *Storage {
	return &Storage{make(map[Byte32]Byte32, 0)}
}

// Helper function to convert Word to big.Int
func (w Word) ToBigInt() *big.Int {
	return new(big.Int).SetBytes(w[:])
}

// Helper function to convert big.Int to Word
func BigIntToWord(val *big.Int) Word {
	var w Word
	bytes := val.Bytes()
	if len(bytes) > 32 {
		// Truncate if too large
		copy(w[:], bytes[len(bytes)-32:])
	} else {
		// Pad with zeros if too small
		copy(w[32-len(bytes):], bytes)
	}
	return w
}

// Helper function to create Word from uint64
func NewWord(val uint64) Word {
	return BigIntToWord(new(big.Int).SetUint64(val))
}

// Helper function to create Word from bytes
func NewWordFromBytes(data []byte) Word {
	var w Word
	if len(data) > 32 {
		copy(w[:], data[len(data)-32:])
	} else {
		copy(w[32-len(data):], data)
	}
	return w
}

// Stack operations
func (s *Stack) Push(data Word) {
	if len(s.Data) >= int(MaximumDepth) {
		panic("stack overflow")
	}
	s.Data = append(s.Data, data)
}

func (s *Stack) Pop() Word {
	if len(s.Data) == 0 {
		panic("stack underflow")
	}
	lastIndex := len(s.Data) - 1
	lastItem := s.Data[lastIndex]
	s.Data = s.Data[:lastIndex]
	return lastItem
}

func (s *Stack) Peek() Word {
	if len(s.Data) == 0 {
		panic("stack underflow")
	}
	return s.Data[len(s.Data)-1]
}

func (s *Stack) PeekAt(index int) Word {
	if index >= len(s.Data) {
		panic("stack underflow")
	}
	return s.Data[len(s.Data)-1-index]
}

func (s *Stack) Dup(index int) {
	if index >= len(s.Data) {
		panic("stack underflow")
	}
	val := s.Data[len(s.Data)-1-index]
	s.Push(val)
}

func (s *Stack) Swap(index int) {
	if index >= len(s.Data) {
		panic("stack underflow")
	}
	top := len(s.Data) - 1
	target := top - index
	s.Data[top], s.Data[target] = s.Data[target], s.Data[top]
}

func (s *Stack) Size() int {
	return len(s.Data)
}

// Arithmetic operations
func (s *Stack) Add() {
	if len(s.Data) < 2 {
		panic("stack underflow")
	}
	a := s.Pop()
	b := s.Pop()
	result := new(big.Int).Add(a.ToBigInt(), b.ToBigInt())
	s.Push(BigIntToWord(result))
}

func (s *Stack) Sub() {
	if len(s.Data) < 2 {
		panic("stack underflow")
	}
	a := s.Pop()
	b := s.Pop()
	result := new(big.Int).Sub(a.ToBigInt(), b.ToBigInt())
	s.Push(BigIntToWord(result))
}

func (s *Stack) Mul() {
	if len(s.Data) < 2 {
		panic("stack underflow")
	}
	a := s.Pop()
	b := s.Pop()
	result := new(big.Int).Mul(a.ToBigInt(), b.ToBigInt())
	s.Push(BigIntToWord(result))
}

func (s *Stack) Div() {
	if len(s.Data) < 2 {
		panic("stack underflow")
	}
	a := s.Pop()
	b := s.Pop()

	if b.ToBigInt().Sign() == 0 {
		// Division by zero returns zero in EVM
		s.Push(NewWord(0))
		return
	}

	result := new(big.Int).Div(a.ToBigInt(), b.ToBigInt())
	s.Push(BigIntToWord(result))
}

func (s *Stack) Mod() {
	if len(s.Data) < 2 {
		panic("stack underflow")
	}
	a := s.Pop()
	b := s.Pop()

	if b.ToBigInt().Sign() == 0 {
		// Modulo by zero returns zero in EVM
		s.Push(NewWord(0))
		return
	}

	result := new(big.Int).Mod(a.ToBigInt(), b.ToBigInt())
	s.Push(BigIntToWord(result))
}

func (s *Stack) AddMod() {
	if len(s.Data) < 3 {
		panic("stack underflow")
	}
	a := s.Pop()
	b := s.Pop()
	n := s.Pop()

	if n.ToBigInt().Sign() == 0 {
		// Modulo by zero returns zero in EVM
		s.Push(NewWord(0))
		return
	}

	result := new(big.Int).Add(a.ToBigInt(), b.ToBigInt())
	result.Mod(result, n.ToBigInt())
	s.Push(BigIntToWord(result))
}

func (s *Stack) MulMod() {
	if len(s.Data) < 3 {
		panic("stack underflow")
	}
	a := s.Pop()
	b := s.Pop()
	n := s.Pop()

	if n.ToBigInt().Sign() == 0 {
		// Modulo by zero returns zero in EVM
		s.Push(NewWord(0))
		return
	}

	result := new(big.Int).Mul(a.ToBigInt(), b.ToBigInt())
	result.Mod(result, n.ToBigInt())
	s.Push(BigIntToWord(result))
}

func (s *Stack) Exp() {
	if len(s.Data) < 2 {
		panic("stack underflow")
	}
	base := s.Pop()
	exp := s.Pop()

	// Limit exponent to prevent excessive computation
	if exp.ToBigInt().BitLen() > 256 {
		panic("exponent too large")
	}

	result := new(big.Int).Exp(base.ToBigInt(), exp.ToBigInt(), nil)
	s.Push(BigIntToWord(result))
}

// Comparison operations
func (s *Stack) Lt() {
	if len(s.Data) < 2 {
		panic("stack underflow")
	}
	a := s.Pop()
	b := s.Pop()

	if a.ToBigInt().Cmp(b.ToBigInt()) < 0 {
		s.Push(NewWord(1))
	} else {
		s.Push(NewWord(0))
	}
}

func (s *Stack) Gt() {
	if len(s.Data) < 2 {
		panic("stack underflow")
	}
	a := s.Pop()
	b := s.Pop()

	if a.ToBigInt().Cmp(b.ToBigInt()) > 0 {
		s.Push(NewWord(1))
	} else {
		s.Push(NewWord(0))
	}
}

func (s *Stack) Eq() {
	if len(s.Data) < 2 {
		panic("stack underflow")
	}
	a := s.Pop()
	b := s.Pop()

	if a.ToBigInt().Cmp(b.ToBigInt()) == 0 {
		s.Push(NewWord(1))
	} else {
		s.Push(NewWord(0))
	}
}

func (s *Stack) IsZero() {
	if len(s.Data) < 1 {
		panic("stack underflow")
	}
	a := s.Pop()

	if a.ToBigInt().Sign() == 0 {
		s.Push(NewWord(1))
	} else {
		s.Push(NewWord(0))
	}
}

// Bitwise operations
func (s *Stack) And() {
	if len(s.Data) < 2 {
		panic("stack underflow")
	}
	a := s.Pop()
	b := s.Pop()

	var result Word
	for i := 0; i < 32; i++ {
		result[i] = a[i] & b[i]
	}
	s.Push(result)
}

func (s *Stack) Or() {
	if len(s.Data) < 2 {
		panic("stack underflow")
	}
	a := s.Pop()
	b := s.Pop()

	var result Word
	for i := 0; i < 32; i++ {
		result[i] = a[i] | b[i]
	}
	s.Push(result)
}

func (s *Stack) Xor() {
	if len(s.Data) < 2 {
		panic("stack underflow")
	}
	a := s.Pop()
	b := s.Pop()

	var result Word
	for i := 0; i < 32; i++ {
		result[i] = a[i] ^ b[i]
	}
	s.Push(result)
}

func (s *Stack) Not() {
	if len(s.Data) < 1 {
		panic("stack underflow")
	}
	a := s.Pop()

	var result Word
	for i := 0; i < 32; i++ {
		result[i] = ^a[i]
	}
	s.Push(result)
}

func (s *Stack) Byte() {
	if len(s.Data) < 2 {
		panic("stack underflow")
	}
	index := s.Pop()
	value := s.Pop()

	idx := int(index.ToBigInt().Uint64())
	if idx >= 32 {
		s.Push(NewWord(0))
		return
	}

	result := uint64(value[idx])
	s.Push(NewWord(result))
}

func (s *Stack) Shl() {
	if len(s.Data) < 2 {
		panic("stack underflow")
	}
	shift := s.Pop()
	value := s.Pop()

	shiftAmount := shift.ToBigInt().Uint64()
	if shiftAmount >= 256 {
		s.Push(NewWord(0))
		return
	}

	result := new(big.Int).Lsh(value.ToBigInt(), uint(shiftAmount))
	s.Push(BigIntToWord(result))
}

func (s *Stack) Shr() {
	if len(s.Data) < 2 {
		panic("stack underflow")
	}
	shift := s.Pop()
	value := s.Pop()

	shiftAmount := shift.ToBigInt().Uint64()
	if shiftAmount >= 256 {
		s.Push(NewWord(0))
		return
	}

	result := new(big.Int).Rsh(value.ToBigInt(), uint(shiftAmount))
	s.Push(BigIntToWord(result))
}

// Memory operations
func (m *Memory) Store(offset byte, value Word) {
	index := int(offset) / 32

	if index >= len(m.Data) {
		newMemoryData := make([]Byte32, index+1)
		copy(newMemoryData, m.Data)
		m.Data = newMemoryData
	}

	copy(m.Data[index][:], value[:])
}

func (m *Memory) Load(offset byte) Word {
	index := int(offset) / 32

	if index > len(m.Data)-1 {
		panic("invalid memory location")
	}

	var result Word
	copy(result[:], m.Data[index][:])
	return result
}

// Storage operations
func (s *Storage) Store(key Word, value Word) {
	var key32, value32 Byte32
	copy(key32[:], key[:])
	copy(value32[:], value[:])
	s.Data[key32] = value32
}

func (s *Storage) Load(key Word) Word {
	var key32 Byte32
	copy(key32[:], key[:])

	value32, exists := s.Data[key32]
	if !exists {
		// Return zero if key doesn't exist
		return NewWord(0)
	}

	var result Word
	copy(result[:], value32[:])
	return result
}

// Utility function to print stack contents
func (s *Stack) Print() {
	fmt.Printf("Stack (size: %d):\n", len(s.Data))
	for i := len(s.Data) - 1; i >= 0; i-- {
		fmt.Printf("  [%d]: %x\n", len(s.Data)-1-i, s.Data[i])
	}
}
