package sort

import (
	"bufio"
	"fmt"
	"os"
	"runhash/config"
	"runhash/hash"
	"strings"
)

const (
	help = `sort [-n <number>] <key> [<node> <...>]

Description:

Select one or more nodes from a list of nodes. Nodes are read:

* as command line arguments
* space delimited nodes from the RUNHASH_NODES environment variable
* if both of the above are not defined, from stdin

Example:

      runhash sort mykey 127.0.0.1 192.168.1.1 10.0.0.1

      RUNHASH_NODES="127.0.0.1 192.168.1.1 10.0.0.1" runhash sort mykey

      echo -e "127.0.0.1\n192.168.1.1\n10.0.0.1" | runhash sort mykey -
`
)

func Run(cfg *config.Config) {
	cfg.Help = help

	nodes := cfg.Nodes

	if cfg.Key == "" {
		cfg.Usage()
	}

	if len(cfg.Args) > 0 {
		nodes = cfg.Args
	}

	if len(nodes) == 0 {
		cfg.Usage()
	}

	if nodes[0] == "-" {
		nodes = readFromStdin()
		if len(nodes) == 0 {
			os.Exit(0)
		}
	}

	if !cfg.Sorted {
		nodes = hash.Sort(cfg.Key, nodes)
	}

	nodes = cfg.Subset(nodes)
	for _, node := range nodes {
		fmt.Println(node)
	}
}

func readFromStdin() (nodes []string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		node := strings.TrimSpace(scanner.Text())
		if node == "" {
			continue
		}
		nodes = append(nodes, node)
	}
	return nodes
}
