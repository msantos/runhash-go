package exec_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	rexec "codeberg.org/msantos/runhash-go/internal/cmd/runhash/exec"
	"codeberg.org/msantos/runhash-go/internal/config"
)

type result struct {
	*config.Config
	exitCode int
	output   string
}

var argv = map[string]result{
	"always": {
		exitCode: 0,
		output:   "ok",
		Config: &config.Config{
			N:       0,
			Key:     "key1",
			Node:    "foo",
			Nodes:   []string{"abc", "foo", "bar"},
			Command: "exec",
			Args:    []string{"echo", "ok"},
		},
	},
	"notselected": {
		exitCode: 0,
		output:   "",
		Config: &config.Config{
			N:       1,
			Key:     "key1",
			Node:    "foo",
			Nodes:   []string{"abc", "foo", "bar"},
			Command: "exec",
			Args:    []string{"echo", "ok"},
		},
	},
	"selected": {
		exitCode: 0,
		output:   "ok",
		Config: &config.Config{
			N:       2,
			Key:     "key1",
			Node:    "foo",
			Nodes:   []string{"abc", "foo", "bar"},
			Command: "exec",
			Args:    []string{"echo", "ok"},
		},
	},
	"notfound": {
		exitCode: 127,
		output:   "",
		Config: &config.Config{
			N:       2,
			Key:     "key1",
			Node:    "foo",
			Nodes:   []string{"abc", "foo", "bar"},
			Command: "exec",
			Args:    []string{"abcdef123xx", "ok"},
		},
	},
	"notfoundandnotselected": {
		exitCode: 0,
		output:   "",
		Config: &config.Config{
			N:       1,
			Key:     "key1",
			Node:    "foo",
			Nodes:   []string{"abc", "foo", "bar"},
			Command: "exec",
			Args:    []string{"abcdef123xx", "ok"},
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

func TestRun_always(t *testing.T) {
	testrun(t, "always")
}

func TestRun_notselected(t *testing.T) {
	testrun(t, "notselected")
}

func TestRun_selected(t *testing.T) {
	testrun(t, "selected")
}

func TestRun_notfound(t *testing.T) {
	testrun(t, "notfound")
}

func TestRun_notfoundandnotselected(t *testing.T) {
	testrun(t, "notfoundandnotselected")
}

func testrun(t *testing.T, name string) {
	if os.Getenv("TESTING_RUNHASH_EXEC_TESTRUN_"+strings.ToUpper(name)) == "1" {
		rexec.Run(argv[name].Config)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestRun_"+name)
	cmd.Env = append(os.Environ(), "TESTING_RUNHASH_EXEC_TESTRUN_"+strings.ToUpper(name)+"=1")
	if err := run(cmd, argv[name]); err != nil {
		t.Errorf("%v", err)
		return
	}
}
