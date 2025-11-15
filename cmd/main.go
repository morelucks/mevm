package main

import (
	"fmt"

	types "github.com/morelucks/minievm/typess"
)

func main() {
	fmt.Println(" MiniEVM ")
	fmt.Println()

	// === Program Counter (PC) Demo ===
	code := []byte{types.PUSH1, 0x2a, types.STOP}
	vm := types.NewVM(code, 10000) // Initialize with 10000 gas
	fmt.Println("PC Demo:")
	fmt.Printf("  Code size: %d bytes\n", vm.CodeSize())
	fmt.Printf("  PC start: %d\n", vm.GetPC())
	op := vm.Fetch()
	fmt.Printf("  Fetched opcode: 0x%x, PC now: %d\n", op, vm.GetPC())
	if op >= types.PUSH1 && op <= types.PUSH32 && vm.HasMore() {
		// For PUSH1, read one immediate byte
		immediate := vm.Fetch()
		fmt.Printf("  PUSH immediate: 0x%x, PC now: %d\n", immediate, vm.GetPC())
	}
	// Peek next (should be STOP) without advancing
	if vm.HasMore() {
		next := vm.PeekByte(0)
		fmt.Printf("  Peek next opcode (no advance): 0x%x, PC still: %d\n", next, vm.GetPC())
	}
	// Advance to STOP and fetch it
	if vm.HasMore() {
		stop := vm.Fetch()
		fmt.Printf("  Fetched opcode: 0x%x (STOP), PC now: %d\n", stop, vm.GetPC())
	}
	fmt.Println()

	// === Gas Model Demo ===
	fmt.Println("=== Gas Model v1 (Pre-Berlin) Demo ===")
	fmt.Println()

	// Demo 1: Basic gas consumption
	fmt.Println("1. Basic Gas Consumption:")
	vm2 := types.NewVM(code, 1000)
	fmt.Printf("  Initial gas: %d\n", vm2.GetGas())
	fmt.Printf("  Gas limit: %d\n", vm2.GetGasLimit())

	// Consume gas for PUSH1
	gasCost := types.GetOpcodeGasCost(types.PUSH1)
	fmt.Printf("  PUSH1 gas cost: %d\n", gasCost)
	err := vm2.ConsumeGas(gasCost)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Gas after PUSH1: %d\n", vm2.GetGas())
	}

	// Consume gas for STOP
	gasCost = types.GetOpcodeGasCost(types.STOP)
	fmt.Printf("  STOP gas cost: %d\n", gasCost)
	err = vm2.ConsumeGas(gasCost)
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	} else {
		fmt.Printf("  Gas after STOP: %d\n", vm2.GetGas())
	}
	fmt.Println()

	// Demo 2: Out of gas scenario
	fmt.Println("2. Out of Gas Scenario:")
	vm3 := types.NewVM(code, 1) // Very low gas limit
	fmt.Printf("  Initial gas: %d\n", vm3.GetGas())
	gasCost = types.GetOpcodeGasCost(types.PUSH1)
	fmt.Printf("  Attempting to consume %d gas for PUSH1...\n", gasCost)
	err = vm3.ConsumeGas(gasCost)
	if err != nil {
		fmt.Printf("  Error (expected): %v\n", err)
	} else {
		fmt.Printf("  Gas remaining: %d\n", vm3.GetGas())
	}
	fmt.Println()

	// Demo 3: Gas costs for different opcodes
	fmt.Println("3. Gas Costs for Different Opcodes:")
	opcodes := []struct {
		name   string
		opcode byte
	}{
		{"STOP", types.STOP},
		{"ADD", types.ADD},
		{"MUL", types.MUL},
		{"EXP", types.EXP},
		{"PUSH1", types.PUSH1},
		{"DUP1", types.DUP1},
		{"SWAP1", types.SWAP1},
	}
	for _, op := range opcodes {
		cost := types.GetOpcodeGasCost(op.opcode)
		fmt.Printf("  %s: %d gas\n", op.name, cost)
	}
	fmt.Println()

	// Demo 4: Memory expansion gas
	fmt.Println("4. Memory Expansion Gas Calculation:")
	oldSize := uint64(0)
	newSize := uint64(64) // 64 bytes = 2 words
	expansionGas := types.MemoryExpansionGas(oldSize, newSize)
	fmt.Printf("  Expanding memory from %d to %d bytes\n", oldSize, newSize)
	fmt.Printf("  Gas cost: %d\n", expansionGas)

	oldSize = 32
	newSize = 128 // 128 bytes = 4 words
	expansionGas = types.MemoryExpansionGas(oldSize, newSize)
	fmt.Printf("  Expanding memory from %d to %d bytes\n", oldSize, newSize)
	fmt.Printf("  Gas cost: %d\n", expansionGas)
	fmt.Println()

	// Demo 5: Gas refund
	fmt.Println("5. Gas Refund (capped at gasLimit/2):")
	vm4 := types.NewVM(code, 10000)
	vm4.ConsumeGas(3000) // Use 3000 gas
	fmt.Printf("  Gas after consuming 3000: %d\n", vm4.GetGas())
	vm4.RefundGas(1000) // Refund 1000
	fmt.Printf("  Gas after refunding 1000: %d\n", vm4.GetGas())

	// Show refund cap
	vm5 := types.NewVM(code, 10000)
	vm5.ConsumeGas(5000) // Use half the gas
	fmt.Printf("  Gas after consuming 5000: %d\n", vm5.GetGas())
	vm5.RefundGas(3000) // Try to refund 3000 (but cap is 5000, already at max)
	fmt.Printf("  Gas after attempting to refund 3000 (capped): %d\n", vm5.GetGas())
	fmt.Println()

	// === Original Stack Operations Demo ===

	// Create a new stack
	stack := types.NewStack()

	fmt.Println("1. Basic Stack Operations:")
	fmt.Println("Pushing values: 5, 10, 15")
	stack.Push(types.NewWord(5))
	stack.Push(types.NewWord(10))
	stack.Push(types.NewWord(15))
	stack.Print()

	fmt.Println("\n2. Arithmetic Operations:")

	// Addition: 10 + 15 = 25
	fmt.Println("Adding top two values (10 + 15):")
	stack.Add()
	stack.Print()

	// Multiplication: 5 * 25 = 125
	fmt.Println("Multiplying top two values (5 * 25):")
	stack.Mul()
	stack.Print()

	// Subtraction: 125 - 5 = 120
	fmt.Println("Subtracting (125 - 5):")
	stack.Push(types.NewWord(5))
	stack.Sub()
	stack.Print()

	// Division: 120 / 4 = 30
	fmt.Println("Dividing (120 / 4):")
	stack.Push(types.NewWord(4))
	stack.Div()
	stack.Print()

	fmt.Println("\n3. Comparison Operations:")

	// Push some values for comparison
	stack.Push(types.NewWord(20))
	stack.Push(types.NewWord(30))
	stack.Print()

	// Less than: 20 < 30 = 1 (true)
	fmt.Println("Less than (20 < 30):")
	stack.Lt()
	stack.Print()

	// Greater than: 30 > 20 = 1 (true)
	fmt.Println("Greater than (30 > 20):")
	stack.Push(types.NewWord(20))
	stack.Push(types.NewWord(30))
	stack.Gt()
	stack.Print()

	// Equal: 30 == 30 = 1 (true)
	fmt.Println("Equal (30 == 30):")
	stack.Push(types.NewWord(30))
	stack.Push(types.NewWord(30))
	stack.Eq()
	stack.Print()

	// Is zero: 0 = 1 (true), 1 = 0 (false)
	fmt.Println("Is zero (0):")
	stack.Push(types.NewWord(0))
	stack.IsZero()
	stack.Print()

	fmt.Println("Is zero (1):")
	stack.Push(types.NewWord(1))
	stack.IsZero()
	stack.Print()

	fmt.Println("\n4. Bitwise Operations:")

	// Clear stack and add some values
	stack = types.NewStack()
	stack.Push(types.NewWord(0xFF)) // 255 in binary: 11111111
	stack.Push(types.NewWord(0x0F)) // 15 in binary:  00001111
	stack.Print()

	// AND: 255 & 15 = 15
	fmt.Println("AND (255 & 15):")
	stack.And()
	stack.Print()

	// OR: 15 | 240 = 255
	fmt.Println("OR (15 | 240):")
	stack.Push(types.NewWord(240))
	stack.Or()
	stack.Print()

	// XOR: 255 ^ 15 = 240
	fmt.Println("XOR (255 ^ 15):")
	stack.Push(types.NewWord(15))
	stack.Xor()
	stack.Print()

	// NOT: ~255 = 0xFFFFFF00... (32 bytes)
	fmt.Println("NOT (~255):")
	stack.Push(types.NewWord(255))
	stack.Not()
	stack.Print()

	fmt.Println("\n5. Stack Manipulation:")

	// Clear and setup
	stack = types.NewStack()
	stack.Push(types.NewWord(1))
	stack.Push(types.NewWord(2))
	stack.Push(types.NewWord(3))
	stack.Push(types.NewWord(4))
	stack.Print()

	// DUP1: Duplicate top element
	fmt.Println("DUP1 (duplicate top element):")
	stack.Dup(0)
	stack.Print()

	// DUP2: Duplicate second element
	fmt.Println("DUP2 (duplicate second element):")
	stack.Dup(1)
	stack.Print()

	// SWAP1: Swap top two elements
	fmt.Println("SWAP1 (swap top two elements):")
	stack.Swap(0)
	stack.Print()

	// SWAP2: Swap top and third elements
	fmt.Println("SWAP2 (swap top and third elements):")
	stack.Swap(1)
	stack.Print()

	fmt.Println("\n6. Advanced Operations:")

	// Modulo: 17 % 5 = 2
	fmt.Println("Modulo (17 % 5):")
	stack = types.NewStack()
	stack.Push(types.NewWord(17))
	stack.Push(types.NewWord(5))
	stack.Mod()
	stack.Print()

	// Exponentiation: 2^8 = 256
	fmt.Println("Exponentiation (2^8):")
	stack = types.NewStack()
	stack.Push(types.NewWord(2))
	stack.Push(types.NewWord(8))
	stack.Exp()
	stack.Print()

	// AddMod: (10 + 20) % 7 = 2
	fmt.Println("AddMod ((10 + 20) % 7):")
	stack = types.NewStack()
	stack.Push(types.NewWord(10))
	stack.Push(types.NewWord(20))
	stack.Push(types.NewWord(7))
	stack.AddMod()
	stack.Print()

	// MulMod: (3 * 4) % 5 = 2
	fmt.Println("MulMod ((3 * 4) % 5):")
	stack = types.NewStack()
	stack.Push(types.NewWord(3))
	stack.Push(types.NewWord(4))
	stack.Push(types.NewWord(5))
	stack.MulMod()
	stack.Print()
}
