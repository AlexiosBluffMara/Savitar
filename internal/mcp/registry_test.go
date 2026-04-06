package mcp

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alexiosbluffmara/savitar/internal/config"
)

func TestLoadVSCodeConfig(t *testing.T) {
	dir := t.TempDir()
	vscodeDir := filepath.Join(dir, ".vscode")
	if err := os.MkdirAll(vscodeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	content := `{
		"servers": {
			"github": {
				"type": "http",
				"url": "https://api.githubcopilot.com/mcp/"
			},
			"context7": {
				"command": "npx",
				"args": ["-y", "@upstash/context7-mcp"]
			}
		}
	}`
	if err := os.WriteFile(filepath.Join(vscodeDir, "mcp.json"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	servers, err := LoadVSCodeConfig(dir)
	if err != nil {
		t.Fatalf("LoadVSCodeConfig returned error: %v", err)
	}
	if len(servers) != 2 {
		t.Fatalf("expected 2 servers, got %d", len(servers))
	}

	byName := map[string]VSCodeServerConfig{}
	for _, s := range servers {
		byName[s.Name] = s
	}

	gh := byName["github"]
	if !gh.IsHTTP() {
		t.Error("github server should be HTTP")
	}
	if gh.URL != "https://api.githubcopilot.com/mcp/" {
		t.Errorf("unexpected URL: %q", gh.URL)
	}

	ctx7 := byName["context7"]
	if !ctx7.IsStdio() {
		t.Error("context7 server should be stdio")
	}
	if ctx7.Command != "npx" {
		t.Errorf("unexpected command: %q", ctx7.Command)
	}
}

func TestLoadVSCodeConfigMissingFile(t *testing.T) {
	servers, err := LoadVSCodeConfig(t.TempDir())
	if err != nil {
		t.Fatalf("expected nil error for missing file, got: %v", err)
	}
	if servers != nil {
		t.Errorf("expected nil servers for missing file, got: %v", servers)
	}
}

func TestRegistryResolvedMergesSavitarEnabled(t *testing.T) {
	dir := t.TempDir()
	vscodeDir := filepath.Join(dir, ".vscode")
	_ = os.MkdirAll(vscodeDir, 0o755)
	content := `{"servers": {"github": {"type": "http", "url": "https://example.com"}, "local": {"command": "echo", "args": []}}}`
	_ = os.WriteFile(filepath.Join(vscodeDir, "mcp.json"), []byte(content), 0o644)

	savitarCfg := []config.MCPServerConfig{
		{Name: "github", Enabled: false},
		{Name: "local", Enabled: true},
	}

	reg := NewRegistry(dir, savitarCfg)
	if err := reg.Load(); err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	resolved := reg.resolved()
	byName := map[string]mergedServer{}
	for _, s := range resolved {
		byName[s.vscode.Name] = s
	}

	if byName["github"].enabled {
		t.Error("github should be disabled per savitar config")
	}
	if !byName["local"].enabled {
		t.Error("local should be enabled per savitar config")
	}
}
