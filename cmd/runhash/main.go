package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"runhash/cmd/runhash/exec"
	"runhash/cmd/runhash/sort"
	"runhash/cmd/runhash/xargs"
	"runhash/config"
	"strings"
)

func usage(cfg *config.Config) {
	fmt.Fprintf(os.Stderr, `%s %s
Usage: %s [-n <number>] <command> [<key> [<...>]]

    Commands:

      sort - sort nodes 
      exec - run command on matching node
      xargs - failover command between nodes
      version - display version

Environment Variables:

`, path.Base(os.Args[0]), config.Version, os.Args[0])

	cfg.PrintDefaults()
	fmt.Fprintf(os.Stderr, "Options:\n\n")
	flag.PrintDefaults()
}

func args() *config.Config {
	cfg, err := config.Env()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	flag.Usage = func() {
		usage(cfg)
	}

	n := flag.Int("n", 0, "number of nodes to return, use 0 for all nodes")
	node := flag.String("node", "", "set identifier for this node")
	nodes := flag.String("nodes", "", "set node list")
	replace := flag.String("replace", "{}",
		"xargs: replace occurrences of string with selected node")
	okExit := flag.Bool("okexit", true,
		"xargs: exit if command returns status 0")
	sorted := flag.Bool("sorted", false,
		"Do not sort: use the existing sort order for nodes")

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	args := flag.Args()

	cfg.N = *n
	cfg.Replace = *replace
	cfg.OKExit = *okExit
	cfg.Sorted = *sorted

	if *node != "" {
		cfg.Node = *node
	}

	if *nodes != "" {
		cfg.Nodes = strings.Fields(*nodes)
	}

	cfg.Command = args[0]

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
	}
}
