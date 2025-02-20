package main

import (
	"fmt"
	types "github.com/morelucks/minievm/typess"
)

func main() {
	// Create a new stack
	stack := types.NewStack()

	//... Push elements ..//
	stack.Push(22)
	stack.Push(54)
	stack.Push(26)

	fmt.Println("Popped:", stack.Pop())

	fmt.Println("Popped:", stack.Pop())
	fmt.Println("Popped:", stack.Pop())

}
