package main

import (
	"encoding/hex"
	"fmt"
	"github.com/elitracy/chat-blockchain/blocks"
	"log"
	"time"
)

func main() {
	blockchain, err := blocks.NewBlockchain()
	if err != nil {
		print(err.Error())
	}

	for i := 0; i < 1; i++ {
		blockchain.AddBlock(fmt.Sprintf("Block: %d", i))
	}

	// for i := 0; i < len(blockchain.Blocks); i++ {
	// 	fmt.Println(blockchain.Blocks[i].Data)
	// 	fmt.Printf(" - parent  > %s\n", hex.EncodeToString(blockchain.Blocks[i].ParentHash[:]))
	// 	fmt.Printf(" - current > %s\n", hex.EncodeToString(blockchain.Blocks[i].Hash[:]))
	// 	fmt.Println("-------------------------------------")
	// }

	err = blockchain.CheckBlocks()

	if err != nil {
		print(err.Error())
		return
	}

	println("VALID BLOCKCHAIN")

	for i := 0; i < len(blockchain.Blocks); i++ {
		encodedBlock := blocks.SerializeBlock(blockchain.Blocks[i])
		decodedBlock := blocks.DeserializeBlock(encodedBlock)

		var unmarshaledTime time.Time

		err := unmarshaledTime.UnmarshalBinary(decodedBlock.Timestamp)
		if err != nil {
			log.Fatal(err.Error())
		}

		timeString := unmarshaledTime.Format(time.RFC3339)

		fmt.Printf(" - timestamp > %s\n", timeString)
		fmt.Printf(" - Data      > %s\n", decodedBlock.Data)
		fmt.Printf(" - parent    > %s\n", hex.EncodeToString(decodedBlock.ParentHash[:]))
		fmt.Printf(" - current   > %s\n", hex.EncodeToString(decodedBlock.Hash[:]))
		fmt.Println("-------------------------------------")
	}

}
