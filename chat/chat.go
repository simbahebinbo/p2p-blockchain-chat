package main

import (
	"blocks"
	"encoding/hex"
	"fmt"
)

func main() {
	blockchain, err := blocks.NewBlockchain()
	if err != nil {
		print(err.Error())
	}

	for i := 0; i < 10; i++ {
		blockchain.AddBlock(fmt.Sprintf("Block: %d", i))
	}

	for i := 0; i < len(blockchain.Blocks); i++ {
		fmt.Println(blockchain.Blocks[i].Data)
		fmt.Printf(" - parent  > %s\n", hex.EncodeToString(blockchain.Blocks[i].ParentHash[:]))
		fmt.Printf(" - current > %s\n", hex.EncodeToString(blockchain.Blocks[i].Hash[:]))
		fmt.Println("-------------------------------------")
	}

	_, err = blockchain.CheckBlocks()

	if err != nil {
		print(err.Error())
		return
	}

	println("VALID BLOCKCHAIN")
}
