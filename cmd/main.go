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
	
	fmt.Println("Popped:", stack.Pop())

	
}
