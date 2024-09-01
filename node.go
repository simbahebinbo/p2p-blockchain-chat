package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
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
			fmt.Printf("Couldn't connect localhost:%d -> %s %s\n", c.port, c.connectionAddress, err.Error())
			time.Sleep(2 * time.Second)
			continue
		} else {

			fmt.Fprintf(conn, "CONN:localhost:%d\n", serverPort)

			conn.Close()

			// fmt.Printf("%d: Holding connection to %s\n", c.port, c.connectionAddress)
			time.Sleep(2 * time.Second)
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

	fmt.Printf("%s connected to all known nodes\n", n.Name)
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
	msg, err := bufio.NewReader(conn).ReadString('\n')

	if err != nil {
		log.Fatal(err.Error())
	}

	if strings.HasPrefix(msg, "CONN:") {
		address := strings.TrimPrefix(msg, "CONN:")
		n.AddPeer(conn, address)

		data := []byte("RECV:localhost:" + string(n.ServerPort) + "\n")
		_, err := conn.Write(data)

		if err != nil {
			log.Printf("Error sending data to %s", address)
			log.Fatal(err.Error())
		}
		conn.Close()
	}

}

func (n *Node) AddPeer(conn net.Conn, peerAddress string) {
	n.PeersLock.Lock()

	if _, exists := n.KnownNodes[peerAddress]; !exists {
		n.KnownNodes[peerAddress] = 0
		fmt.Printf("%s added peer %s", n.Name, peerAddress)
		n.PeersLock.Unlock()
	} else {
		n.PeersLock.Unlock()
	}

}

func (n *Node) AddChat() {
	var chat string

	fmt.Print(">")
	fmt.Scan(&chat)
	n.Blockchain.AddBlock(chat)

	currentBlock := n.Blockchain.Blocks[len(n.Blockchain.Blocks)-1]

	var unmarshaledTime time.Time

	err := unmarshaledTime.UnmarshalBinary(currentBlock.Timestamp)
	if err != nil {
		log.Fatal(err.Error())
	}

	timeString := unmarshaledTime.Format(time.RFC3339)

	fmt.Println("===================")
	fmt.Printf("Created At: %s\n", timeString)
	fmt.Printf("Parent:     %x\n", currentBlock.ParentHash)
	fmt.Printf("Hash:       %x\n", currentBlock.Hash)
	fmt.Printf("Message:    %s\n", currentBlock.Data)
	fmt.Println("===================")
}

func (n *Node) Start() {
	blockchain, err := blocks.NewBlockchain()

	if err != nil {
		log.Fatal(err)
	}

	n.Blockchain = *blockchain

	go n.StartServer()
	go n.ConnectToNodes()

	for {
		n.AddChat()
	}
}
