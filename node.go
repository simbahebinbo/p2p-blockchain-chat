package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/elitracy/chat-blockchain/blocks"
)

type Node struct {
	Name       string
	ServerPort int

	KnownNodes map[string]int // 0 for disconnected, 1 for connected
	PeersLock  sync.Mutex

	Clients []Client

	Blockchain blocks.Blockchain
}

type Client struct {
	port              int
	connectionAddress string
}

func (c *Client) Connect(serverPort int) {
	for {
		conn, err := net.Dial("tcp", c.connectionAddress)
		if err != nil {
			fmt.Print("\033[2K\r")
			fmt.Printf("[localhost:%d] Retrying %s\n", c.port, c.connectionAddress)
			fmt.Printf("[YOU]>")
			time.Sleep(2 * time.Second)
			continue
		} else {

			var buffer [1024]byte
			messageType := []byte("CONN:")
			host := []byte(fmt.Sprintf("localhost:%d", serverPort))

			copy(buffer[:32], messageType)
			copy(buffer[32:], host)

			conn.Write(buffer[:])
			conn.Close()
		}
	}
}

func (n *Node) ConnectToNodes() {
	for addr := range n.KnownNodes {
		port := n.ServerPort + len(n.Clients) + 1
		client := Client{connectionAddress: addr, port: port}
		n.Clients = append(n.Clients, client)

		go client.Connect(n.ServerPort)
	}
}

func (n *Node) StartServer() {

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", n.ServerPort))

	if err != nil {
		log.Fatal(err.Error())
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}

		go n.HandleClient(conn)
	}
}

func (n *Node) HandleClient(conn net.Conn) {
	var buffer [1024]byte
	_, err := bufio.NewReader(conn).Read(buffer[:])

	messageType := string(buffer[:32])

	if err != nil {
		log.Printf("Failed to read message from: %s\n %s", conn.LocalAddr().String(), err.Error())
	}

	if strings.HasPrefix(messageType, "CONN:") {
		address := string(buffer[32:])
		n.AddPeer(conn, address)

		data := []byte("RECV:localhost:" + string(n.ServerPort) + "\n")
		_, err := conn.Write(data)

		if err != nil {
			log.Printf("Error sending data to %s\n", address)
			log.Fatal(err.Error())
		}

	} else if strings.HasPrefix(messageType, "NEWBLOCK:") {
		blockBytes := buffer[32:]

		block, err := blocks.DeserializeBlock(blockBytes)

		if err != nil {
			fmt.Printf("Error deserializing block: %s\n", err.Error())
			return
		}

		n.Blockchain.Blocks = append(n.Blockchain.Blocks, block)

		data := []byte("NEWBLOCKCONF:\n")
		_, err = conn.Write(data)

		if err != nil {
			log.Printf("Error sending data to %s\n", conn.LocalAddr().String())
			log.Fatal(err.Error())
		}

		fmt.Print("\033[2K\r")

		fmt.Printf("[%s]>%s\n", block.Data.Owner, block.Data.Message)
		fmt.Printf("[YOU]>")

	}
	conn.Close()

}

func (n *Node) AddPeer(conn net.Conn, peerAddress string) {
	n.PeersLock.Lock()

	if _, exists := n.KnownNodes[peerAddress]; !exists {
		n.KnownNodes[peerAddress] = 0
		fmt.Print("\033[2K\r")
		fmt.Printf("[localhost:%d] Added %s\n", n.ServerPort, peerAddress)
		fmt.Printf("[YOU]>")
		n.PeersLock.Unlock()
	} else {
		n.PeersLock.Unlock()
	}

}

func (n *Node) AddChat() {
	input := bufio.NewReader(os.Stdin)

	fmt.Print("[YOU]>")
	line, err := input.ReadString('\n')

	if err != nil {
		fmt.Printf("Error reading input: %s", err.Error())
	}

	chat := blocks.Chat{Message: line[:len(line)-1], Owner: n.Name}

	n.Blockchain.AddBlock(chat)

	currentBlock := n.Blockchain.Blocks[len(n.Blockchain.Blocks)-1]

	// n.Blockchain.WriteBlock("blockchain", currentBlock)
	n.ShareBlock(currentBlock)
}

func (n *Node) ShareBlock(block *blocks.Block) {

	for _, client := range n.Clients {
		conn, err := net.Dial("tcp", client.connectionAddress)
		if err != nil {
			fmt.Printf("localhost:%d couldn't share block with:  %s %s\n", client.port, client.connectionAddress, err.Error())
		} else {

			blockBytes, err := blocks.SerializeBlock(block)

			if err != nil {
				fmt.Printf("Failed to serialize block")
			}

			var messageBuffer [1024]byte

			messageType := []byte("NEWBLOCK:")
			copy(messageBuffer[:32], messageType)
			copy(messageBuffer[32:], blockBytes)

			_, err = conn.Write(messageBuffer[:])
			if err != nil {
				fmt.Printf("Failed to send block")
			}

		}
		conn.Close()
	}
}

func (n *Node) Start() {

	fmt.Print("\033[H\033[2J") // clear screen

	blockchain, err := blocks.NewBlockchain()

	if err != nil {
		log.Fatal(err)
	}

	n.Blockchain = *blockchain

	// read blockchain from file
	// n.Blockchain.ReadBlockchain("blockchain")
	// n.Blockchain.PrintBlocks()

	// connect to nodes
	go n.StartServer()
	n.ConnectToNodes()

	// ask for current chain

	// validate current chain against own chain

	// chat
	for {
		n.AddChat()
	}
}
