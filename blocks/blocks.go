package blocks

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"
	"unsafe"
)

type Block struct {
	Timestamp  []byte
	ParentHash []byte
	Hash       []byte
	Data       string
}

func (block *Block) AddHash() {
	headers := bytes.Join([][]byte{block.Timestamp, block.ParentHash, []byte(block.Data)}, []byte{})
	hash := sha256.Sum256(headers)
	block.Hash = hash[:]
}

type Blockchain struct {
	Blocks []*Block
}

func AddBlock(data string, parentHash []byte) (*Block, error) {
	timeNow, err := time.Now().MarshalBinary()
	if err != nil {
		return nil, err
	}
	block := &Block{timeNow, parentHash[:], []byte{}, data}
	block.AddHash()

	return block, nil
}

func AddGenesisBlock() (*Block, error) {
	block, err := AddBlock("Genesis Block", []byte{})
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

func (blockchain *Blockchain) CheckBlocks() error {

	for i := 0; i < len(blockchain.Blocks); i++ {
		if i == 0 {
			continue
		}

		parent := blockchain.Blocks[i-1]
		child := blockchain.Blocks[i]

		if !bytes.Equal(parent.Hash, child.ParentHash) {
			errString := fmt.Sprintf("Hashes do not match! \n PARENT: %s> %d \n CHILD: %s> %d", parent.Data, parent.Hash, child.Data, child.Hash)
			return errors.New(errString)
		}

	}

	return nil
}

func SerializeBlock(block *Block) []byte {
	var buf bytes.Buffer

	buf.Write(block.Timestamp)     // 24
	buf.Write(block.ParentHash[:]) // 32
	buf.Write(block.Hash[:])       // 32
	buf.Write([]byte(block.Data))

	return buf.Bytes()
}

func DeserializeBlock(data []byte) *Block {

	block := &Block{}

	timestamp_s := unsafe.Sizeof(block.Timestamp)
	parenthash_s := unsafe.Sizeof(block.ParentHash)
	hash_s := unsafe.Sizeof(block.Hash)
	data_s := unsafe.Sizeof(block.Data)

	// block.Timestamp = data[:timestamp_s-1]
	// print("sliced timestamp")
	// copy(block.ParentHash[:], data[timestamp_s:timestamp_s+hash_s-1])
	// print("sliced parent")
	// copy(block.Hash[:], data[timestamp_s+hash_s:timestamp_s+hash_s+hash_s-1])
	// print("sliced hash")
	// block.Data = string(data[timestamp_s+hash_s+hash_s : timestamp_s+hash_s+hash_s+data_s])
	// print("sliced data")

	var buf = bytes.NewBuffer(data)

	timestamp_b := make([]byte, timestamp_s)
	buf.Read(timestamp_b)

	parenthash_b := make([]byte, parenthash_s)
	buf.Read(parenthash_b)

	hash_b := make([]byte, hash_s)
	buf.Read(hash_b)

	data_b := make([]byte, data_s)
	buf.Read(data_b)

	block.Timestamp = timestamp_b
	block.ParentHash = []byte(parenthash_b)
	block.Hash = []byte(hash_b)
	block.Data = string(data_b)

	return block
}
