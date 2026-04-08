package runtime

import (
	"strings"
	"testing"

	"github.com/alexiosbluffmara/savitar/internal/config"
)

func TestIntegrationStatusesReportLocalOllamaLane(t *testing.T) {
	loaded := config.Loaded{Config: config.Default(), Path: "config/savitar.local.json"}

	rt := NewAtRoot(loaded, t.TempDir())
	statuses := rt.IntegrationStatuses()

	for _, status := range statuses {
		if status.Name != "ollama" {
			continue
		}

		if !status.Enabled {
			t.Fatal("expected local ollama lane to be reported as enabled")
		}
		if status.AuthSource != "local-http" {
			t.Fatalf("unexpected ollama auth source: %q", status.AuthSource)
		}
		if status.TokenEnv != "" {
			t.Fatalf("expected no ollama token env, got %q", status.TokenEnv)
		}
		if !strings.Contains(status.Details, loaded.Config.Models.LocalDefault.Model) {
			t.Fatalf("expected ollama details to mention the local model, got %q", status.Details)
		}
		return
	}

	t.Fatal("expected to find ollama integration status")
}
