package customization

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverAgentsAndSkills(t *testing.T) {
	root := t.TempDir()
	agentsDir := filepath.Join(root, ".github", "agents")
	skillsDir := filepath.Join(root, ".github", "skills", "demo-skill")

	if err := os.MkdirAll(agentsDir, 0o755); err != nil {
		t.Fatalf("MkdirAll agents returned error: %v", err)
	}

	if err := os.MkdirAll(skillsDir, 0o755); err != nil {
		t.Fatalf("MkdirAll skills returned error: %v", err)
	}

	agentFile := filepath.Join(agentsDir, "demo.agent.md")
	skillFile := filepath.Join(skillsDir, "SKILL.md")

	if err := os.WriteFile(agentFile, []byte(`---
name: "Demo Agent"
description: "Use when testing workspace discovery"
tools: [read, search]
user-invocable: false
---
body
`), 0o600); err != nil {
		t.Fatalf("WriteFile agent returned error: %v", err)
	}

	if err := os.WriteFile(skillFile, []byte(`---
name: demo-skill
description: "Use when testing workspace discovery"
argument-hint: "describe the task"
---
body
`), 0o600); err != nil {
		t.Fatalf("WriteFile skill returned error: %v", err)
	}

	agents, err := DiscoverAgents(root)
	if err != nil {
		t.Fatalf("DiscoverAgents returned error: %v", err)
	}

	skills, err := DiscoverSkills(root)
	if err != nil {
		t.Fatalf("DiscoverSkills returned error: %v", err)
	}

	if len(agents) != 1 {
		t.Fatalf("expected 1 agent, got %d", len(agents))
	}

	if agents[0].Name != "Demo Agent" {
		t.Fatalf("unexpected agent name: %q", agents[0].Name)
	}

	if agents[0].UserInvocable {
		t.Fatal("expected agent to be non-invocable")
	}

	if len(agents[0].Tools) != 2 {
		t.Fatalf("expected 2 tools, got %d", len(agents[0].Tools))
	}

	if len(skills) != 1 {
		t.Fatalf("expected 1 skill, got %d", len(skills))
	}

	if skills[0].ArgumentHint != "describe the task" {
		t.Fatalf("unexpected skill argument hint: %q", skills[0].ArgumentHint)
	}
}
