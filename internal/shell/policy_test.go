package shell

import (
	"context"
	"testing"
)

func TestPolicyAllowsAllowlistedCommand(t *testing.T) {
	p := DefaultPolicy()
	result := p.Run(context.Background(), "echo", "hello")
	if result.Err != nil {
		t.Fatalf("unexpected error: %v", result.Err)
	}
	if result.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d (stderr: %s)", result.ExitCode, result.Stderr)
	}
}

func TestPolicyBlocksDisallowedCommand(t *testing.T) {
	p := &Policy{Allowlist: []string{"echo"}}
	result := p.Run(context.Background(), "rm", "-rf", "/")
	if result.Err == nil {
		t.Fatal("expected error for disallowed command, got nil")
	}
}

func TestPolicyCustomAllowlist(t *testing.T) {
	p := &Policy{Allowlist: []string{"echo"}}
	allowed := p.Run(context.Background(), "echo", "ok")
	if allowed.Err != nil {
		t.Fatalf("allowed command returned error: %v", allowed.Err)
	}

	denied := p.Run(context.Background(), "git", "status")
	if denied.Err == nil {
		t.Error("expected git to be denied when not in allowlist")
	}
}

func TestPolicyCapturesOutput(t *testing.T) {
	p := DefaultPolicy()
	result := p.Run(context.Background(), "echo", "captured output")
	if result.Err != nil {
		t.Fatalf("unexpected error: %v", result.Err)
	}
	if result.Stdout == "" {
		t.Error("expected non-empty stdout")
	}
}
