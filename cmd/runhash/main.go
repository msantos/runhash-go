package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"go.iscode.ca/runhash/cmd/runhash/internal/exec"
	"go.iscode.ca/runhash/cmd/runhash/internal/sort"
	"go.iscode.ca/runhash/cmd/runhash/internal/xargs"
	"go.iscode.ca/runhash/internal/config"
)

func usage() {
	fmt.Fprintf(os.Stderr, `%s %s
Usage: %s [-n <number>] <command> [<key> [<...>]]

Commands:

  sort - sort nodes
  exec - run command on matching node
  xargs - failover command between nodes
  version - display version

Environment Variables:

`, path.Base(os.Args[0]), config.Version, os.Args[0])

	config.PrintDefaults()
	fmt.Fprintf(os.Stderr, "Options:\n\n")
	flag.PrintDefaults()
}

func args() *config.Config {
	n := flag.Int("n", 0, "number of nodes to return, use 0 for all nodes")
	node := flag.String("node", config.Nodename(), "set identifier for this node")
	nodes := flag.String("nodes", config.Nodes(), "set node list")
	replace := flag.String("replace", "{}",
		"xargs: replace occurrences of string with selected node")
	okExit := flag.Bool("okexit", true,
		"xargs: exit if command returns status 0")
	sorted := flag.Bool("sorted", false,
		"Do not sort: use the existing sort order for nodes")

	flag.Usage = func() {
		usage()
	}

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}

	args := flag.Args()

	cfg := &config.Config{
		N:       *n,
		Replace: *replace,
		OKExit:  *okExit,
		Sorted:  *sorted,
		Node:    *node,
		Nodes:   strings.Fields(*nodes),
		Command: args[0],
	}

	if len(args) > 1 {
		cfg.Key = args[1]
	}

	if len(args) > 2 {
		cfg.Args = args[2:]
	}

	return cfg
}

func main() {
	cfg := args()

	switch cfg.Command {
	case "exec":
		exec.Run(cfg)
	case "sort":
		sort.Run(cfg)
	case "xargs":
		xargs.Run(cfg)
	case "version":
		fmt.Fprintf(os.Stderr, "%s\n", config.Version)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cfg.Command)
		flag.Usage()
		os.Exit(127)
	}
}
