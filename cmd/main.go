package main

import (
	"fmt"

	types "github.com/morelucks/minievm/typess"
)

func main() {
	fmt.Println("=== MiniEVM - Milestone 1: Core Arithmetic & Stack Operations ===\n")

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

	fmt.Println("\n=== Milestone 1 Complete! ===")
	fmt.Println("You've successfully implemented core EVM arithmetic and stack operations!")
	fmt.Println("Next: Add gas system and execution limits (Milestone 2)")
}
