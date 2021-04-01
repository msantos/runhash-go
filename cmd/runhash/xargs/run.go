package xargs

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runhash/config"
	"runhash/hash"
	"strings"
	"syscall"
)

const (
	help = `xargs <key> <command> <...>

Description:

Failover a command between a sorted list of nodes. By default, the command
is executed sequentially on the sorted list of all nodes. If the command
exits non-0, the command is run again with the next node in the list.

'{}' in the command arguments is replaced with the selected node.

Example:

    RUNHASH_NODES="127.0.0.1 127.1.1.1" runhash xargs mykey nc {} 443

`
)

func Run(cfg *config.Config) {
	cfg.Help = help

	if cfg.Key == "" {
		cfg.Usage()
	}

	if len(cfg.Nodes) == 0 {
		cfg.Usage()
	}

	nodes := cfg.Nodes
	if !cfg.Sorted {
		nodes = hash.Sort(cfg.Key, cfg.Nodes)
	}

	nodes = cfg.Subset(nodes)

	var cmd string
	var oarg []string
	switch len(cfg.Args) {
	case 0:
		cfg.Usage()
	case 1:
		cmd = cfg.Args[0]
	default:
		cmd = cfg.Args[0]
		oarg = cfg.Args[1:]
	}

	exe, err := exec.LookPath(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s:%s\n", exe, err)
		os.Exit(127)
	}

	exitStatus := 0

	for _, node := range nodes {
		arg := replace(oarg, node, cfg.Replace)
		exitStatus = execv(exe, arg, syscall.Environ())
		if cfg.OKExit && exitStatus == 0 {
			break
		}
	}
	os.Exit(exitStatus)
}

func replace(arg []string, node, template string) []string {
	narg := make([]string, len(arg))
	for i, s := range arg {
		narg[i] = strings.ReplaceAll(s, template, node)
	}
	return narg
}

func execv(command string, args []string, env []string) int {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env

	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 111
	}
	waitCh := make(chan error, 1)
	go func() {
		waitCh <- cmd.Wait()
		close(waitCh)
	}()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan)

	for {
		select {
		case sig := <-sigChan:
			_ = cmd.Process.Signal(sig)
		case err := <-waitCh:
			var exitError *exec.ExitError
			if errors.As(err, &exitError) {
				waitStatus := exitError.Sys().(syscall.WaitStatus)
				if waitStatus.Signaled() {
					return 128 + int(waitStatus.Signal())
				}
				return waitStatus.ExitStatus()
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return 111
			}
			return 0
		}
	}
}
