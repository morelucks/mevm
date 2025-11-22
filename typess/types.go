package types

import "fmt"

const (
	MaximumDepth uint = 1024
	WordSize     int  = 32 // 256 bits = 32 bytes
)

type Byte32 [32]byte

// EVM uses 256-bit words (32 bytes)
type Word [32]byte

// Stack holds 256-bit values
type Stack struct {
	Data []Word
}

type Memory struct {
	Data []Byte32
}

type Storage struct {
	Data map[Byte32]Byte32
}

// VM holds code, program counter (PC), gas tracking, stack, and memory
// Implements machine state (μ) from Yellow Paper Section 9
type VM struct {
	Code     []byte
	PC       uint64  // μ_pc - Program counter
	Gas      uint64  // μ_g - Remaining gas
	GasLimit uint64  // Maximum gas allowed
	Stack    *Stack  // μ_s - Stack contents
	Memory   *Memory // μ_m - Memory contents
}

// EVM Opcodes
const (
	// Stack operations
	STOP       = 0x00
	ADD        = 0x01
	MUL        = 0x02
	SUB        = 0x03
	DIV        = 0x04
	SDIV       = 0x05
	MOD        = 0x06
	SMOD       = 0x07
	ADDMOD     = 0x08
	MULMOD     = 0x09
	EXP        = 0x0a
	SIGNEXTEND = 0x0b

	// Comparison operations
	LT     = 0x10
	GT     = 0x11
	SLT    = 0x12
	SGT    = 0x13
	EQ     = 0x14
	ISZERO = 0x15

	// Bitwise operations
	AND  = 0x16
	OR   = 0x17
	XOR  = 0x18
	NOT  = 0x19
	BYTE = 0x1a
	SHL  = 0x1b
	SHR  = 0x1c
	SAR  = 0x1d

	// Stack manipulation
	POP    = 0x50
	PUSH1  = 0x60
	PUSH2  = 0x61
	PUSH3  = 0x62
	PUSH4  = 0x63
	PUSH5  = 0x64
	PUSH6  = 0x65
	PUSH7  = 0x66
	PUSH8  = 0x67
	PUSH9  = 0x68
	PUSH10 = 0x69
	PUSH11 = 0x6a
	PUSH12 = 0x6b
	PUSH13 = 0x6c
	PUSH14 = 0x6d
	PUSH15 = 0x6e
	PUSH16 = 0x6f
	PUSH17 = 0x70
	PUSH18 = 0x71
	PUSH19 = 0x72
	PUSH20 = 0x73
	PUSH21 = 0x74
	PUSH22 = 0x75
	PUSH23 = 0x76
	PUSH24 = 0x77
	PUSH25 = 0x78
	PUSH26 = 0x79
	PUSH27 = 0x7a
	PUSH28 = 0x7b
	PUSH29 = 0x7c
	PUSH30 = 0x7d
	PUSH31 = 0x7e
	PUSH32 = 0x7f

	// Duplication operations
	DUP1  = 0x80
	DUP2  = 0x81
	DUP3  = 0x82
	DUP4  = 0x83
	DUP5  = 0x84
	DUP6  = 0x85
	DUP7  = 0x86
	DUP8  = 0x87
	DUP9  = 0x88
	DUP10 = 0x89
	DUP11 = 0x8a
	DUP12 = 0x8b
	DUP13 = 0x8c
	DUP14 = 0x8d
	DUP15 = 0x8e
	DUP16 = 0x8f

	// Exchange operations
	SWAP1  = 0x90
	SWAP2  = 0x91
	SWAP3  = 0x92
	SWAP4  = 0x93
	SWAP5  = 0x94
	SWAP6  = 0x95
	SWAP7  = 0x96
	SWAP8  = 0x97
	SWAP9  = 0x98
	SWAP10 = 0x99
	SWAP11 = 0x9a
	SWAP12 = 0x9b
	SWAP13 = 0x9c
	SWAP14 = 0x9d
	SWAP15 = 0x9e
	SWAP16 = 0x9f
)

// Gas cost constants (Istanbul fork - pre-Berlin)
const (
	GasZero         uint64 = 0
	GasBase         uint64 = 2
	GasVeryLow      uint64 = 3
	GasLow          uint64 = 5
	GasMid          uint64 = 8
	GasHigh         uint64 = 10
	GasExtStep      uint64 = 20
	GasExtCode      uint64 = 700
	GasBalance      uint64 = 400
	GasSLoad        uint64 = 200
	GasSStore       uint64 = 20000
	GasSStoreReset  uint64 = 5000
	GasSStoreNoop   uint64 = 200
	GasCall         uint64 = 700
	GasCreate       uint64 = 32000
	GasMemory       uint64 = 3 // Per word (32 bytes)
	GasLog          uint64 = 375
	GasLogTopic     uint64 = 375
	GasLogData      uint64 = 8
	GasExp          uint64 = 10
	GasExpByte      uint64 = 50
	GasSHA3         uint64 = 30
	GasSHA3Word     uint64 = 6
	GasCopy         uint64 = 3
	GasCopyWord     uint64 = 3
	GasJumpDest     uint64 = 1
	GasSelfDestruct uint64 = 5000
)

// OutOfGasError represents when execution runs out of gas
type OutOfGasError struct {
	Required  uint64
	Remaining uint64
}

func (e *OutOfGasError) Error() string {
	return fmt.Sprintf("out of gas: required %d, remaining %d", e.Required, e.Remaining)
}
