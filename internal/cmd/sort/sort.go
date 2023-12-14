package sort

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"codeberg.org/msantos/runhash-go/internal/config"
	"codeberg.org/msantos/runhash-go/internal/hash"
)

const (
	help = `sort <key> [<node> <...>]

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

	var err error
	if nodes[0] == "-" {
		nodes, err = readFromStdin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "stdin: %s\n", err)
			os.Exit(1)
		}
		if len(nodes) == 0 {
			os.Exit(0)
		}
	}

	if !cfg.Sorted {
		nodes = hash.Sort(cfg.Key, nodes)
	}

	for _, node := range cfg.Subset(nodes) {
		fmt.Println(node)
	}
}

func readFromStdin() (nodes []string, err error) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		node := strings.TrimSpace(scanner.Text())
		if node == "" {
			continue
		}
		nodes = append(nodes, node)
	}
	return nodes, scanner.Err()
}
