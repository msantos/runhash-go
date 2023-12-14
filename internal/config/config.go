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

func Env() (*Config, error) {
	nodes := getenv("RUNHASH_NODES", "")

	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("gethostname: %w", err)
	}
	node := getenv("RUNHASH_NODE", hostname)

	return &Config{
		Node:  node,
		Nodes: strings.Fields(nodes),
	}, nil
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

func (cfg *Config) PrintDefaults() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	fmt.Fprintf(os.Stderr,
		"  RUNHASH_NODE=\"%s\"\n    %s (default %s)\n\n",
		cfg.Node,
		"Node identifier",
		hostname,
	)

	fmt.Fprintf(os.Stderr,
		"  RUNHASH_NODES=\"%s\"\n    %s (default %s)\n\n",
		strings.Join(cfg.Nodes, " "),
		"Space separated list of nodes",
		"\"\"",
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
