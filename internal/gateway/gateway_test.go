package gateway

import (
	"testing"

	"github.com/alexiosbluffmara/savitar/internal/config"
)

func TestBuildPlanIncludesConfiguredSurfaces(t *testing.T) {
	cfg := config.Default()
	cfg.Transports.Discord.Enabled = true
	cfg.Transports.Discord.DisplayName = "Savitar Ops"
	cfg.Transports.WhatsApp.Enabled = true
	cfg.WebUI.Enabled = true

	plan := BuildPlan(cfg)
	if len(plan.Surfaces) != 4 {
		t.Fatalf("expected 4 surfaces, got %d", len(plan.Surfaces))
	}

	if plan.EnabledCount() != 3 {
		t.Fatalf("expected 3 enabled surfaces, got %d", plan.EnabledCount())
	}

	if plan.Surfaces[0].Identity != "Savitar Ops" {
		t.Fatalf("unexpected discord identity: %q", plan.Surfaces[0].Identity)
	}
}
