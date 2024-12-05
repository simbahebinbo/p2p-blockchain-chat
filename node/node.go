package node

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

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

func (n *Node) ConnectToPeer(addr string) {
	port := n.ServerPort + len(n.Clients) + 1
	client := Client{connectionAddress: addr, port: port}
	n.Clients = append(n.Clients, client)

	connectionStatusChannel := make(chan int)

	go func() {
		connectionStatusChannel <- client.Connect(n.ServerPort)
	}()

	n.KnownNodes[addr] = <-connectionStatusChannel
}

func (n *Node) ConnectToKnownNodes() {
	for addr := range n.KnownNodes {
		if n.KnownNodes[addr] == 0 {
			n.ConnectToPeer(addr)
		}
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

		if _, exists := n.KnownNodes[address]; !exists {
			n.AddPeer(conn, address)
		}

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
			Logger.Error("Error deserializing block: %s\n", err.Error(), nil)
			return
		}

		n.Blockchain.Blocks = append(n.Blockchain.Blocks, block)

		data := []byte("NEWBLOCKCONF:\n")
		_, err = conn.Write(data)

		if err != nil {
			log.Printf("Error sending data to %s\n", conn.LocalAddr().String())
			log.Fatal(err.Error())
		}

		Console.Write("[%s]>%s\n", block.Data.Owner, block.Data.Message)
	}
	conn.Close()

}

func (n *Node) AddPeer(conn net.Conn, peerAddress string) {
	n.PeersLock.Lock()

	if _, exists := n.KnownNodes[peerAddress]; !exists {
		n.KnownNodes[peerAddress] = 0
		Console.Write("[localhost:%d] Added %s\n", n.ServerPort, peerAddress)
		n.PeersLock.Unlock()
	} else {
		n.PeersLock.Unlock()
	}

	// if val, _ := n.KnownNodes[peerAddress]; val == 0 {
	// 	n.ConnectToPeer(peerAddress)
	// }
}

func (n *Node) AddChat() {
	input := bufio.NewReader(os.Stdin)

	line, err := input.ReadString('\n')

	if err != nil {
		Logger.Error("Error reading input: %s", err.Error(), nil)
	}

	chat := blocks.Chat{Message: line[:len(line)-1], Owner: n.Name}

	n.Blockchain.AddBlock(chat)

	currentBlock := n.Blockchain.Blocks[len(n.Blockchain.Blocks)-1]

	// n.Blockchain.WriteBlock("blockchain", currentBlock)
	n.ShareBlock(currentBlock)

	Console.Write("")
}

func (n *Node) ShareBlock(block *blocks.Block) {

	for _, client := range n.Clients {
		conn, err := net.Dial("tcp", client.connectionAddress)
		if err != nil {
			errMsg := fmt.Sprintf("localhost:%d couldn't share block with:  %s %s\n", client.port, client.connectionAddress, err.Error())
			Logger.Error(errMsg)
		} else {

			blockBytes, err := blocks.SerializeBlock(block)
			if err != nil {
				Logger.Error("Failed to serialize block")
			}

			var messageBuffer [1024]byte
			messageType := []byte("NEWBLOCK:")
			copy(messageBuffer[:32], messageType)
			copy(messageBuffer[32:], blockBytes)

			_, err = conn.Write(messageBuffer[:])
			if err != nil {
				Logger.Error("Failed to send block")
			}

		}
		conn.Close()
	}
}

func (n *Node) Start() {
	Console.Clear()
	Console.Write("========= %s ========= localhost:%d =========\n", n.Name, n.ServerPort)

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
	n.ConnectToKnownNodes()

	// ask for current chain

	// validate current chain against own chain

	// chat
	for {
		n.AddChat()
	}
}
