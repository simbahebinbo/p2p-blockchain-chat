package node

import (
	"fmt"
	"net"
	"time"
)

const MAX_CONNECTION_TRIES = 10

type Client struct {
	port              int
	connectionAddress string
}

// returns 1 once connected or 0 if max connection retries has been reached
func (c *Client) Connect(serverPort int) int {
	connectionStatus := 0
	tries := 0

	for connectionStatus == 0 || tries == MAX_CONNECTION_TRIES {
		Console.Write("[localhost:%d] Trying %s\n", c.port, c.connectionAddress)

		conn, err := net.Dial("tcp", c.connectionAddress)
		if err != nil {
			time.Sleep(2 * time.Second)
			tries++
		} else {

			var buffer [1024]byte
			messageType := []byte("CONN:")
			host := []byte(fmt.Sprintf("localhost:%d", serverPort))

			Console.Write("[localhost:%d] Connecting to %s\n", serverPort, c.connectionAddress)

			copy(buffer[:32], messageType)
			copy(buffer[32:], host)

			conn.Write(buffer[:])
			conn.Close()
			return 1
		}
	}
	return 0
}
