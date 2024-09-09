package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/elitracy/chat-blockchain/node"
)

const NUM_INIT_NODES = 3

func main() {
	sPort := flag.Int("s", 8000, "the node's server port")
	name := flag.String("n", "unnamed node", "the node's name")
	knownPorts := flag.String("k", "-1", "comma separated list of ports")

	flag.Parse()

	n := node.Node{
		KnownNodes: make(map[string]int),
		Name:       *name,
		ServerPort: *sPort,
	}

	var ports []string

	if *knownPorts != "-1" {
		ports = strings.Split(*knownPorts, ",")
	}

	for _, port := range ports {
		host := fmt.Sprintf("localhost:%s", port)
		n.KnownNodes[host] = 0
	}

	n.Start()

	select {}
}
