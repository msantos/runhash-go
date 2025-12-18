package exec

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"go.iscode.ca/runhash/internal/config"
	"go.iscode.ca/runhash/internal/hash"
)

const (
	help = `exec <key> <command> <...>

Description:

Execute a command on a subset of nodes chosen from the list of nodes
based on the key. The number of nodes chosen is set by the -n option.

The command will always run on this node if either:
* RUNHASH_NODES is empty
* the node is included in RUNHASH_NODES and -n is 0 (all nodes)

Example:

    RUNHASH_NODES="$(uname -n) foo bar" runhash -n 1 exec mykey ls -al

`
)

func Run(cfg *config.Config) {
	cfg.Help = help

	if cfg.Key == "" {
		cfg.Usage()
	}

	if len(cfg.Args) == 0 {
		cfg.Usage()
	}

	if !selectedNode(cfg) {
		os.Exit(0)
	}

	cmd := cfg.Args[0]
	arg := cfg.Args[0:]

	exe, err := exec.LookPath(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s:%s\n", exe, err)
		os.Exit(127)
	}

	if err := syscall.Exec(exe, arg, syscall.Environ()); err != nil {
		fmt.Fprintf(os.Stderr, "%s:%s\n", exe, err)
	}
	os.Exit(126)
}

func selectedNode(cfg *config.Config) bool {
	if len(cfg.Nodes) == 0 {
		return true
	}

	values := cfg.Nodes
	if !cfg.Sorted {
		values = hash.Sort(cfg.Key, cfg.Nodes)
	}
	for _, v := range cfg.Subset(values) {
		if v == cfg.Node {
			return true
		}
	}
	return false
}
