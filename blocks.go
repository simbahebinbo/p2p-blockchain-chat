package blocks

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"
)

type Block struct {
	Timestamp  []byte
	ParentHash [32]byte
	Hash       [32]byte
	Data       string
}

func (block *Block) AddHash() {
	headers := bytes.Join([][]byte{block.Timestamp, block.ParentHash[:], []byte(block.Data)}, []byte{})
	hash := sha256.Sum256(headers)
	block.Hash = hash
}

type Blockchain struct {
	Blocks []*Block
}

func AddBlock(data string, parentHash [32]byte) (*Block, error) {
	timeNow, err := time.Now().MarshalBinary()
	if err != nil {
		return nil, err
	}
	block := &Block{timeNow, parentHash, [32]byte{}, data}
	block.AddHash()

	return block, nil
}

func AddGenesisBlock() (*Block, error) {
	block, err := AddBlock("Genesis Block", [32]byte{})
	if err != nil {
		return nil, err
	}

	return block, nil
}

func NewBlockchain() (*Blockchain, error) {

	blockchain := &Blockchain{}

	genesisBlock, err := AddGenesisBlock()
	if err != nil {
		return nil, err
	}

	blockchain.Blocks = append(blockchain.Blocks, genesisBlock)

	return blockchain, nil
}

func (blockchain *Blockchain) AddBlock(data string) {

	block, err := AddBlock(data, blockchain.Blocks[len(blockchain.Blocks)-1].Hash)

	if err != nil {
		print(err.Error())
	}

	blockchain.Blocks = append(blockchain.Blocks, block)
}

func (blockchain *Blockchain) CheckBlocks() (bool, error) {

	for i := 0; i < len(blockchain.Blocks); i++ {
		if i == 0 {
			continue
		}

		parent := blockchain.Blocks[i-1]
		child := blockchain.Blocks[i]

		if parent.Hash != child.ParentHash {
			errString := fmt.Sprintf("Hashes do not match! \n PARENT: %s> %d \n CHILD: %s> %d", parent.Data, parent.Hash, child.Data, child.Hash)
			return false, errors.New(errString)
		}

	}

	return true, nil
}
