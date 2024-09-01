package main

import (
	"flag"
	"fmt"
	"strings"
)

const NUM_INIT_NODES = 3

func main() {
	sPort := flag.Int("s", 8000, "the node's server port")
	name := flag.String("n", "unnamed node", "the node's name")
	knownPorts := flag.String("k", "", "comma separated list of ports")

	flag.Parse()

	n := Node{KnownNodes: make(map[string]int),
		Name:       *name,
		ServerPort: *sPort,
	}

	ports := strings.Split(*knownPorts, ",")

	for i := range ports {
		port := fmt.Sprintf("localhost:%s", ports[i])
		n.KnownNodes[port] = 0
	}

	n.Start()

	select {}
}
