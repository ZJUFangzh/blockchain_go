package main

import (
	"fmt"
)

func main() {

	blockChain := NewBlockChain()
	blockChain.AddBlock("add 1 btc")
	blockChain.AddBlock("add 2 btc")

	for i, block := range blockChain.blocks {
		fmt.Printf("The %d block\n", i)
		fmt.Printf("Prehash is %x\n", block.PreBlockHash)
		fmt.Printf("Data is %s\n", block.Data)
		fmt.Printf("Hash is %x\n", block.Hash)
		fmt.Println()
	}
}
