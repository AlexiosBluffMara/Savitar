package session

import (
	"path/filepath"
	"testing"
)

func TestLoadReturnsDefaultStateWhenMissing(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), "session.json"))
	report, err := store.Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if report.Exists {
		t.Fatal("expected report.Exists to be false")
	}

	if report.State.CurrentSurface != "cli" {
		t.Fatalf("unexpected current surface: %q", report.State.CurrentSurface)
	}
}

func TestInitWritesState(t *testing.T) {
	store := NewStore(filepath.Join(t.TempDir(), ".savitar", "session.json"))
	report, err := store.Init(State{
		Version:             1,
		CurrentSurface:      "discord",
		CurrentModelProfile: "copilot-0x",
		LastCommand:         "status",
	})
	if err != nil {
		t.Fatalf("Init returned error: %v", err)
	}

	if !report.Exists {
		t.Fatal("expected report.Exists to be true")
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if loaded.State.CurrentModelProfile != "copilot-0x" {
		t.Fatalf("unexpected current model profile: %q", loaded.State.CurrentModelProfile)
	}
}
