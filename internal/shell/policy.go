package shell

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// DefaultAllowlist is the set of commands permitted by default.
// Only base command names are listed; arguments are passed through.
var DefaultAllowlist = []string{
	"git",
	"gh",
	"go",
	"make",
	"cat",
	"head",
	"tail",
	"grep",
	"find",
	"ls",
	"wc",
	"sort",
	"uniq",
	"diff",
	"echo",
	"true",
	"false",
}

// Policy controls what shell commands are permitted to run.
type Policy struct {
	// Allowlist is the set of permitted base command names.
	// If empty, DefaultAllowlist is used.
	Allowlist []string
	// WorkDir is the working directory for all commands.
	// If empty, the process's current directory is used.
	WorkDir string
	// MaxRuntime is the per-command timeout.
	// If zero, 60 seconds is used.
	MaxRuntime time.Duration
}

// Result captures the outcome of a command run.
type Result struct {
	Command  string
	Args     []string
	Stdout   string
	Stderr   string
	ExitCode int
	Err      error
}

// Run executes command with the given args under the policy.
// It returns the captured output, error, and a non-zero exit code if the
// command failed. Policy violations are returned as errors without running
// the command.
func (p *Policy) Run(ctx context.Context, command string, args ...string) Result {
	base := baseCommand(command)
	if !p.permitted(base) {
		return Result{
			Command: command,
			Args:    args,
			Err:     fmt.Errorf("command %q is not in the shell allowlist", base),
		}
	}

	timeout := p.MaxRuntime
	if timeout <= 0 {
		timeout = 60 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, command, args...)
	if p.WorkDir != "" {
		cmd.Dir = p.WorkDir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
			err = nil // non-zero exit is not a Go error; caller checks ExitCode
		}
	}

	return Result{
		Command:  command,
		Args:     args,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
		Err:      err,
	}
}

func (p *Policy) permitted(base string) bool {
	list := p.Allowlist
	if len(list) == 0 {
		list = DefaultAllowlist
	}
	for _, name := range list {
		if strings.EqualFold(name, base) {
			return true
		}
	}
	return false
}

func baseCommand(command string) string {
	trimmed := strings.TrimSpace(command)
	// Strip any directory prefix so "git" and "/usr/bin/git" both match "git".
	parts := strings.Split(trimmed, "/")
	return parts[len(parts)-1]
}

// DefaultPolicy returns a Policy with the default allowlist, no working-dir
// restriction, and a 60-second timeout.
func DefaultPolicy() *Policy {
	return &Policy{}
}
