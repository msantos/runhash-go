package sort_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	rsort "codeberg.org/msantos/runhash-go/cmd/runhash/internal/sort"
	"codeberg.org/msantos/runhash-go/internal/config"
)

type result struct {
	*config.Config
	exitCode int
	output   string
	stdin    string
}

var argv = map[string]result{
	"env": {
		exitCode: 0,
		output: `192.168.1.1
10.0.0.1
127.0.0.1`,
		Config: &config.Config{
			N:   0,
			Key: "mykey",
			Nodes: []string{
				"192.168.1.1",
				"10.0.0.1",
				"127.0.0.1",
			},
			Command: "sort",
		},
	},
	"argv": {
		exitCode: 0,
		output: `192.168.1.1
10.0.0.1
127.0.0.1`,
		Config: &config.Config{
			N:       0,
			Key:     "mykey",
			Nodes:   []string{},
			Command: "sort",
			Args: []string{
				"192.168.1.1",
				"10.0.0.1",
				"127.0.0.1",
			},
		},
	},
	"stdin": {
		exitCode: 0,
		output: `192.168.1.1
10.0.0.1
127.0.0.1`,
		stdin: `127.0.0.1
10.0.0.1
192.168.1.1
`,
		Config: &config.Config{
			N:       0,
			Key:     "mykey",
			Nodes:   []string{},
			Command: "sort",
			Args:    []string{"-"},
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

	var stdin io.Reader = os.Stdin
	if r.stdin != "" {
		stdin = bytes.NewBuffer([]byte(r.stdin))
	}

	cmd.Stdin = stdin

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
	if !strings.HasPrefix(output, r.output) {
		return fmt.Errorf("Expected: %s\nOutput: %s\nError: %w",
			r.output,
			output,
			errInvalidOutput,
		)
	}

	return nil
}

func TestRun_env(t *testing.T) {
	testrun(t, "env")
}

func TestRun_argv(t *testing.T) {
	testrun(t, "argv")
}

func TestRun_stdin(t *testing.T) {
	testrun(t, "stdin")
}

func testrun(t *testing.T, name string) {
	if os.Getenv("TESTING_RUNHASH_SORT_TESTRUN_"+strings.ToUpper(name)) == "1" {
		rsort.Run(argv[name].Config)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestRun_"+name)
	cmd.Env = append(os.Environ(), "TESTING_RUNHASH_SORT_TESTRUN_"+strings.ToUpper(name)+"=1")
	if err := run(cmd, argv[name]); err != nil {
		t.Errorf("%v", err)
		return
	}
}
