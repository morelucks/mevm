package main

import (
	"fmt"
	types "github.com/morelucks/minievm/typess"
)

func main() {
	// Create a new stack
	stack := types.NewStack()

	// Push elements
	stack.Push(22)
	
	fmt.Println("Popped:", stack.Pop())

	// Uncommenting the next line will cause a panic (stack underflow)
	// fmt.Println("Popped:", stack.Pop())
}
