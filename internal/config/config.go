package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Nodes   []string
	Node    string
	Command string
	Key     string
	Args    []string
	N       int
	Replace string
	OKExit  bool
	Sorted  bool
	Help    string
}

var Version = "0.0.0"

func getenv(k, def string) string {
	if v, ok := os.LookupEnv(k); ok {
		return v
	}
	return def
}

func Nodename() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	return getenv("RUNHASH_NODE", hostname)
}

func Nodes() string {
	return getenv("RUNHASH_NODES", "")
}

func (cfg *Config) Usage() {
	if cfg.Help != "" {
		fmt.Fprintf(os.Stderr, "Usage: %s %s\n", os.Args[0], cfg.Help)
	}
	os.Exit(2)
}

func (cfg *Config) String() string {
	return fmt.Sprintf(
		"RUNHASH_NODE=\"%s\"\nRUNHASH_NODES=\"%s\"",
		cfg.Node,
		strings.Join(cfg.Nodes, " "),
	)
}

func PrintDefaults() {
	fmt.Fprintf(
		os.Stderr,
		"  RUNHASH_NODE=\"%s\"\n    Node identifier\n\n  RUNHASH_NODES=\"%s\"\n    Space separated list of nodes\n\n",
		Nodename(), Nodes(),
	)
}

func (cfg *Config) Subset(nodes []string) []string {
	if cfg.N <= 0 {
		return nodes
	}
	return nodes[:min(len(nodes), cfg.N)]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
