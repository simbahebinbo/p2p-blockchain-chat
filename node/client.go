package node

import (
	"fmt"
	"net"
	"time"
)

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
