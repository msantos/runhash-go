package xargs_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	rxargs "codeberg.org/msantos/runhash-go/cmd/runhash/internal/xargs"
	"codeberg.org/msantos/runhash-go/internal/config"
)

type result struct {
	*config.Config
	exitCode int
	output   string
}

var argv = map[string]result{
	"okexit": {
		exitCode: 0,
		output:   "192.168.1.1",
		Config: &config.Config{
			Key: "key1",
			Nodes: []string{
				"192.168.1.1",
				"10.0.0.1",
				"127.0.0.1",
			},
			Command: "xargs",
			Args:    []string{"echo", "{}"},
			Replace: "{}",
			OKExit:  true,
		},
	},
	"exitcode": {
		exitCode: 111,
		output: `10.0.0.1
127.0.0.1
192.168.1.1`,
		Config: &config.Config{
			Key: "key2",
			Nodes: []string{
				"192.168.1.1",
				"10.0.0.1",
				"127.0.0.1",
			},
			Command: "xargs",
			Args:    []string{"/bin/sh", "-c", "echo '#@#'; exit 111"},
			Replace: "#@#",
			OKExit:  true,
		},
	},
	"all": {
		exitCode: 0,
		output: `+10.0.0.1+
+127.0.0.1+
+192.168.1.1+`,
		Config: &config.Config{
			Key: "key2",
			Nodes: []string{
				"192.168.1.1",
				"10.0.0.1",
				"127.0.0.1",
			},
			Command: "xargs",
			Args:    []string{"echo", "+{}+"},
			Replace: "{}",
			OKExit:  false,
		},
	},
}

var (
	errInvalidOutput = errors.New("unexpected output")
	errExitCode      = errors.New("unexpected exit code")
)

func run(cmd *exec.Cmd, r result) error {
	var buf bytes.Buffer

	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	var ee *exec.ExitError

	err := cmd.Run()
	if err != nil {
		if !errors.As(err, &ee) {
			return err
		}
	}

	if cmd.ProcessState.ExitCode() != r.exitCode {
		return fmt.Errorf("Expected: %d\nExitCode: %d\nError: %w",
			r.exitCode,
			cmd.ProcessState.ExitCode(),
			errExitCode,
		)
	}

	output := strings.TrimSpace(buf.String())
	if output != r.output {
		return fmt.Errorf("Expected: %s\nOutput: %s\nError: %w",
			r.output,
			output,
			errInvalidOutput,
		)
	}

	return nil
}

func TestRun_okexit(t *testing.T) {
	testrun(t, "okexit")
}

func TestRun_exitcode(t *testing.T) {
	testrun(t, "exitcode")
}

func TestRun_all(t *testing.T) {
	testrun(t, "all")
}

func testrun(t *testing.T, name string) {
	if os.Getenv("TESTING_RUNHASH_XARGS_TESTRUN_"+strings.ToUpper(name)) == "1" {
		rxargs.Run(argv[name].Config)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestRun_"+name)
	cmd.Env = append(os.Environ(), "TESTING_RUNHASH_XARGS_TESTRUN_"+strings.ToUpper(name)+"=1")
	if err := run(cmd, argv[name]); err != nil {
		t.Errorf("%v", err)
		return
	}
}
