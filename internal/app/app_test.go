package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunHelpIncludesNewCommands(t *testing.T) {
	tempDir := t.TempDir()
	previousDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned error: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Chdir returned error: %v", err)
	}
	defer func() {
		_ = os.Chdir(previousDir)
	}()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exitCode := Run(stdout, stderr, "dev", []string{"help"})
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}

	output := stdout.String()
	for _, command := range []string{"savitar status", "savitar agents", "savitar skills", "savitar integrations", "savitar gateway", "savitar persona", "savitar session [show|init|list]", "savitar discord [status|preview|run]", "savitar mcp [status]", "savitar repo analyze", "savitar memory [list"} {
		if !strings.Contains(output, command) {
			t.Fatalf("expected help output to contain %q", command)
		}
	}
}

func TestRunSessionInitCreatesLocalState(t *testing.T) {
	tempDir := t.TempDir()
	previousDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned error: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Chdir returned error: %v", err)
	}
	defer func() {
		_ = os.Chdir(previousDir)
	}()

	if err := os.MkdirAll(filepath.Join(tempDir, "config"), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exitCode := Run(stdout, stderr, "dev", []string{"session", "init"})
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d with stderr %q", exitCode, stderr.String())
	}

	sessionPath := filepath.Join(tempDir, ".savitar", "session-index.json")
	if _, err := os.Stat(sessionPath); err != nil {
		t.Fatalf("expected session file at %s: %v", sessionPath, err)
	}

	if !strings.Contains(stdout.String(), "initialized") {
		t.Fatalf("expected session output to mention initialization, got %q", stdout.String())
	}
}

func TestRunDiscordPreviewReturnsReply(t *testing.T) {
	tempDir := t.TempDir()
	previousDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned error: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Chdir returned error: %v", err)
	}
	defer func() {
		_ = os.Chdir(previousDir)
	}()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exitCode := Run(stdout, stderr, "dev", []string{"discord", "preview", "status"})
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d with stderr %q", exitCode, stderr.String())
	}

	if !strings.Contains(stdout.String(), "Savitar status") {
		t.Fatalf("expected preview output to contain status text, got %q", stdout.String())
	}
}

func TestRunIntegrationsReturnsProviderTable(t *testing.T) {
	tempDir := t.TempDir()
	previousDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd returned error: %v", err)
	}

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Chdir returned error: %v", err)
	}
	defer func() {
		_ = os.Chdir(previousDir)
	}()

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exitCode := Run(stdout, stderr, "dev", []string{"integrations"})
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d with stderr %q", exitCode, stderr.String())
	}

	output := stdout.String()
	for _, provider := range []string{"ollama", "github", "huggingface", "kaggle"} {
		if !strings.Contains(output, provider) {
			t.Fatalf("expected integrations output to contain %q, got %q", provider, output)
		}
	}
}
