package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadReturnsDefaultsWhenFileMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if loaded.Exists {
		t.Fatal("expected missing config to report Exists=false")
	}

	if loaded.Config.Agent.Name != "Savitar" {
		t.Fatalf("unexpected default agent name: %q", loaded.Config.Agent.Name)
	}

	if loaded.Config.WebUI.AuthProvider != "google-oauth" {
		t.Fatalf("unexpected default web auth provider: %q", loaded.Config.WebUI.AuthProvider)
	}

	if !loaded.Config.Transports.Discord.RequireMention {
		t.Fatal("expected Discord transport to require mentions by default")
	}

	if loaded.Config.Transports.Discord.MaxResponseChars != 1800 {
		t.Fatalf("unexpected default Discord max response chars: %d", loaded.Config.Transports.Discord.MaxResponseChars)
	}

	if loaded.Config.Transports.Discord.PerUserCooldownSeconds != 5 {
		t.Fatalf("unexpected default Discord per-user cooldown: %d", loaded.Config.Transports.Discord.PerUserCooldownSeconds)
	}

	if loaded.Config.Transports.Discord.MaxConcurrentReplies != 2 {
		t.Fatalf("unexpected default Discord concurrent reply limit: %d", loaded.Config.Transports.Discord.MaxConcurrentReplies)
	}

	if loaded.Config.Transports.Discord.AllowCloudRepliesInDMs {
		t.Fatal("expected Discord cloud DMs to be disabled by default")
	}

	if loaded.Config.Transports.Discord.AllowCloudRepliesInGuilds {
		t.Fatal("expected Discord cloud guild replies to be disabled by default")
	}

	if loaded.Config.Transports.Discord.AllowLiveWebLookupInDMs {
		t.Fatal("expected Discord live web lookup in DMs to be disabled by default")
	}

	if loaded.Config.Transports.Discord.AllowLiveWebLookupInGuilds {
		t.Fatal("expected Discord live web lookup in guilds to be disabled by default")
	}

	if len(loaded.Config.Transports.Discord.OperatorUserIDs) != 0 {
		t.Fatalf("expected default Discord operator user allowlist to be empty, got %d entries", len(loaded.Config.Transports.Discord.OperatorUserIDs))
	}

	if loaded.Config.Integrations.Ollama.APIKeyEnv != "OLLAMA_API_KEY" {
		t.Fatalf("unexpected default Ollama API key env: %q", loaded.Config.Integrations.Ollama.APIKeyEnv)
	}

	if loaded.Config.Integrations.Kaggle.TokenEnv != "KAGGLE_API_TOKEN" {
		t.Fatalf("unexpected default Kaggle token env: %q", loaded.Config.Integrations.Kaggle.TokenEnv)
	}

	if !loaded.Config.Knowledge.RequireSourceMetadata {
		t.Fatal("expected knowledge retrieval to require source metadata by default")
	}

	if loaded.Config.Operator.RunLogDir != ".savitar/runs" {
		t.Fatalf("unexpected default operator run log dir: %q", loaded.Config.Operator.RunLogDir)
	}

	if len(loaded.Config.Knowledge.RepoMarkdownDirs) != 4 {
		t.Fatalf("expected default repo markdown dirs, got %d entries", len(loaded.Config.Knowledge.RepoMarkdownDirs))
	}

	if !loaded.Config.Knowledge.EnableLiveWebLookup {
		t.Fatal("expected live web lookup to be enabled by default")
	}

	if loaded.Config.Knowledge.LiveWebProvider != "duckduckgo-json" {
		t.Fatalf("unexpected live web provider: %q", loaded.Config.Knowledge.LiveWebProvider)
	}
}

func TestLoadOverridesDefaultsFromFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "savitar.local.json")
	data := []byte(`{"agent":{"name":"Testitar"},"automation":{"allowShellExecution":false}}`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if !loaded.Exists {
		t.Fatal("expected config to report Exists=true")
	}

	if loaded.Config.Agent.Name != "Testitar" {
		t.Fatalf("unexpected loaded agent name: %q", loaded.Config.Agent.Name)
	}

	if loaded.Config.Automation.AllowShellExecution {
		t.Fatal("expected AllowShellExecution to be overridden to false")
	}

	if loaded.Config.Models.LocalDefault.Model != "gemma4:e4b" {
		t.Fatalf("expected default model to remain populated, got %q", loaded.Config.Models.LocalDefault.Model)
	}

	if loaded.Config.Transports.IMessage.AccountEnv != "SAVITAR_IMESSAGE_ACCOUNT" {
		t.Fatalf("expected default iMessage account env to remain populated, got %q", loaded.Config.Transports.IMessage.AccountEnv)
	}
}
