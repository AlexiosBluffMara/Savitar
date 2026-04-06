package runtime

import (
	"testing"

	"github.com/alexiosbluffmara/savitar/internal/config"
)

func TestIntegrationStatusesReportEnvPresence(t *testing.T) {
	t.Setenv("OLLAMA_API_KEY", "test-ollama")

	loaded := config.Loaded{Config: config.Default(), Path: "config/savitar.local.json"}
	loaded.Config.Integrations.Ollama.Enabled = true

	rt := NewAtRoot(loaded, t.TempDir())
	statuses := rt.IntegrationStatuses()

	for _, status := range statuses {
		if status.Name != "ollama" {
			continue
		}

		if !status.CredentialPresent {
			t.Fatal("expected ollama credential to be reported as present")
		}
		if status.TokenEnv != "OLLAMA_API_KEY" {
			t.Fatalf("unexpected ollama token env: %q", status.TokenEnv)
		}
		return
	}

	t.Fatal("expected to find ollama integration status")
}
