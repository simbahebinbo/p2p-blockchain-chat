package node

import (
	"fmt"
	"net"
)

type Client struct {
	port              int
	connectionAddress string
}

func (c *Client) Connect(serverPort int) int {
	for {
		conn, err := net.Dial("tcp", c.connectionAddress)
		if err != nil {
			return 0
		} else {

			var buffer [1024]byte
			messageType := []byte("CONN:")
			host := []byte(fmt.Sprintf("localhost:%d", serverPort))

			copy(buffer[:32], messageType)
			copy(buffer[32:], host)

			conn.Write(buffer[:])
			conn.Close()
			return 1
		}
	}
}
