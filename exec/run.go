package exec

import (
	"fmt"
	"os"
	"os/exec"
	"runhash/config"
	"runhash/hash"
	"syscall"
)

const (
	help = `exec <key> <command> <...>

Description:

Execute a command on one node chosen from the list of nodes based on the
key. If RUNHASH_NODES is empty, the command will always run on this node.

Example:

    RUNHASH_NODES="$(uname -n) foo bar" runhash exec mykey ls -al
`
)

func Run(cfg *config.Config) {
	cfg.Help = help

	if cfg.Key == "" {
		cfg.Exit()
	}

	if len(cfg.Args) == 0 {
		cfg.Exit()
	}

	run := selectedNode(cfg)

	if !run {
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
	os.Exit(127)
}

func selectedNode(cfg *config.Config) bool {
	if len(cfg.Nodes) == 0 {
		return true
	}

	values := cfg.Nodes
	if !cfg.Sorted {
		values = hash.Sort(cfg.Key, cfg.Nodes)
	}
	if values[0] == cfg.Node {
		return true
	}
	return false
}
