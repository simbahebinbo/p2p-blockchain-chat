package blocks

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

type Block struct {
	Timestamp  []byte
	ParentHash []byte
	Hash       []byte
	Data       Chat
}

type Chat struct {
	Message string
	Owner   string
}

func (block *Block) AddHash() {
	headers := bytes.Join([][]byte{block.Timestamp, block.ParentHash, []byte(block.Data.Owner), []byte(block.Data.Message)}, []byte{})
	hash := sha256.Sum256(headers)
	block.Hash = hash[:]
}

type Blockchain struct {
	Blocks []*Block
}

func AddBlock(data Chat, parentHash []byte) (*Block, error) {
	timeNow, err := time.Now().MarshalBinary()
	if err != nil {
		return nil, err
	}
	block := &Block{timeNow, parentHash[:], []byte{}, data}
	block.AddHash()

	return block, nil
}

func AddGenesisBlock() (*Block, error) {

	block, err := AddBlock(Chat{"Genesis Block", "root"}, []byte{})
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

func (blockchain *Blockchain) AddBlock(data Chat) {

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

func (blockchain *Blockchain) PrintBlocks() {
	for _, block := range blockchain.Blocks {
		var unmarshaledTime time.Time

		err := unmarshaledTime.UnmarshalBinary(block.Timestamp)
		if err != nil {
			log.Fatal(err.Error())
		}

		timeString := unmarshaledTime.Format(time.RFC3339)
		fmt.Println("===================")
		fmt.Printf("Created At: %s\n", timeString)
		fmt.Printf("Parent:     %x\n", block.ParentHash)
		fmt.Printf("Hash:       %x\n", block.Hash)
		fmt.Printf("Data:    %s\n", block.Data)
		fmt.Println("===================")
	}
}

func SerializeBlock(block *Block) ([]byte, error) {
	var buf bytes.Buffer

	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(block)

	if err != nil {
		// log.Fatal("Failed to encode block: ", err.Error())
		return nil, err
	}

	return buf.Bytes(), nil
}

func DeserializeBlock(data []byte) (*Block, error) {

	block := &Block{}

	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)

	err := decoder.Decode(block)

	if err != nil {
		return nil, err
	}

	return block, nil

}

func (b *Blockchain) WriteBlock(filename string, block *Block) error {
	blockBytes, err := SerializeBlock(block)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	_, err = f.Write(blockBytes)

	if err != nil {
		return err
	}

	return err
}

func (b *Blockchain) ReadBlockchain(filename string) error {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	for i := 0; i < len(fileData)-88; i += 88 {
		block, err := DeserializeBlock(fileData[i : i+88])
		if err != nil {
			return err
		}
		b.Blocks = append(b.Blocks, block)
		println("got block")
	}

	b.CheckBlocks()

	return nil

}
