package persona

import (
	"testing"

	"github.com/alexiosbluffmara/savitar/internal/config"
)

func TestFromConfigProvidesFallbackName(t *testing.T) {
	profile := FromConfig(config.AgentConfig{})
	if profile.Name != "Savitar" {
		t.Fatalf("expected fallback name Savitar, got %q", profile.Name)
	}
}

func TestRequiresExplicitDisclosure(t *testing.T) {
	profile := FromConfig(config.AgentConfig{
		Name: "Savitar",
		Persona: config.AgentPersonaConfig{
			DisclosurePolicy: "Always disclose when identity matters",
		},
	})

	if !profile.RequiresExplicitDisclosure() {
		t.Fatal("expected disclosure policy to require explicit disclosure")
	}
}
