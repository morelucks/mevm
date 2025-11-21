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

func NewVM(code []byte, gasLimit uint64) *VM {
	// Keep an internal copy of code to avoid external mutation
	internal := make([]byte, len(code))
	copy(internal, code)
	return &VM{
		Code:     internal,
		PC:       0,
		Gas:      gasLimit,
		GasLimit: gasLimit,
	}
}

// Program Counter (PC) helpers
func (vm *VM) CodeSize() uint64 {
	return uint64(len(vm.Code))
}

func (vm *VM) HasMore() bool {
	return vm.PC < uint64(len(vm.Code))
}

func (vm *VM) GetPC() uint64 {
	return vm.PC
}

func (vm *VM) SetPC(pc uint64) {
	if pc > uint64(len(vm.Code)) {
		panic("pc out of bounds")
	}
	vm.PC = pc
}

func (vm *VM) IncPC(delta uint64) {
	newPC := vm.PC + delta
	if newPC > uint64(len(vm.Code)) {
		panic("pc out of bounds")
	}
	vm.PC = newPC
}

func (vm *VM) Fetch() byte {
	if vm.PC >= uint64(len(vm.Code)) {
		panic("pc out of bounds")
	}
	b := vm.Code[vm.PC]
	vm.PC++
	return b
}

func (vm *VM) PeekByte(offset uint64) byte {
	index := vm.PC + offset
	if index >= uint64(len(vm.Code)) {
		panic("pc out of bounds")
	}
	return vm.Code[index]
}

// Gas management methods
func (vm *VM) GetGas() uint64 {
	return vm.Gas
}

func (vm *VM) GetGasLimit() uint64 {
	return vm.GasLimit
}

func (vm *VM) ConsumeGas(amount uint64) error {
	if vm.Gas < amount {
		return &OutOfGasError{
			Required:  amount,
			Remaining: vm.Gas,
		}
	}
	vm.Gas -= amount
	return nil
}

func (vm *VM) RefundGas(amount uint64) {
	// Gas refunds are capped at gasLimit/2 (pre-Berlin)
	maxRefund := vm.GasLimit / 2
	currentRefund := vm.GasLimit - vm.Gas
	if currentRefund+amount > maxRefund {
		vm.Gas = vm.GasLimit - maxRefund
	} else {
		vm.Gas += amount
		if vm.Gas > vm.GasLimit {
			vm.Gas = vm.GasLimit
		}
	}
}

// MemoryExpansionGas calculates gas cost for memory expansion
// Formula: (newSize^2 - oldSize^2) / 512 + 3 * (newSize - oldSize)
// where sizes are in words (32 bytes)
func MemoryExpansionGas(oldSize, newSize uint64) uint64 {
	if newSize <= oldSize {
		return 0
	}

	// Convert to words (32 bytes per word)
	oldWords := (oldSize + 31) / 32
	newWords := (newSize + 31) / 32

	if newWords <= oldWords {
		return 0
	}

	// Quadratic cost: (newWords^2 - oldWords^2) / 512
	quadCost := (newWords*newWords - oldWords*oldWords) / 512

	// Linear cost: 3 * (newWords - oldWords)
	linearCost := 3 * (newWords - oldWords)

	return quadCost + linearCost
}

// Execute runs the bytecode interpreter loop
func (vm *VM) Execute() error {
	for vm.HasMore() {
		opcode := vm.Fetch()
		_ = opcode // my TODO
	}
	return nil
}

// GetOpcodeGasCost returns the base gas cost for an opcode (Istanbul fork)
// Note: Some opcodes have dynamic costs (EXP, SHA3, memory ops, etc.)
// that need additional calculation
func GetOpcodeGasCost(opcode byte) uint64 {
	switch opcode {
	case STOP:
		return GasZero
	case ADD, SUB, LT, GT, SLT, SGT, EQ, AND, OR, XOR, NOT, BYTE, SHL, SHR, SAR, ISZERO:
		return GasVeryLow
	case MUL, DIV, SDIV, MOD, SMOD, SIGNEXTEND:
		return GasLow
	case ADDMOD, MULMOD:
		return GasMid
	case EXP:
		return GasExp // Base cost, additional cost per byte
	case POP:
		return GasBase
	case PUSH1, PUSH2, PUSH3, PUSH4, PUSH5, PUSH6, PUSH7, PUSH8,
		PUSH9, PUSH10, PUSH11, PUSH12, PUSH13, PUSH14, PUSH15, PUSH16,
		PUSH17, PUSH18, PUSH19, PUSH20, PUSH21, PUSH22, PUSH23, PUSH24,
		PUSH25, PUSH26, PUSH27, PUSH28, PUSH29, PUSH30, PUSH31, PUSH32:
		return GasVeryLow
	case DUP1, DUP2, DUP3, DUP4, DUP5, DUP6, DUP7, DUP8,
		DUP9, DUP10, DUP11, DUP12, DUP13, DUP14, DUP15, DUP16:
		return GasVeryLow
	case SWAP1, SWAP2, SWAP3, SWAP4, SWAP5, SWAP6, SWAP7, SWAP8,
		SWAP9, SWAP10, SWAP11, SWAP12, SWAP13, SWAP14, SWAP15, SWAP16:
		return GasVeryLow
	default:
		// Unknown opcode - return high cost to discourage execution
		return GasHigh
	}
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
