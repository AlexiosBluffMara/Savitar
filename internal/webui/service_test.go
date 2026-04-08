package webui

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alexiosbluffmara/savitar/internal/config"
	"github.com/alexiosbluffmara/savitar/internal/operator"
	savitarruntime "github.com/alexiosbluffmara/savitar/internal/runtime"
)

type deniedAuthenticator struct{}

func (deniedAuthenticator) SessionFromRequest(*http.Request) (Session, error) {
	return Session{}, os.ErrPermission
}

func TestRoutesServeDashboardInDemoMode(t *testing.T) {
	tempDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tempDir, "docs", "hackathon"), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "docs", "hackathon", "demo-plan.md"), []byte("# Demo\n\nSavitar web UI demo source.\n"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	loaded := config.Loaded{Config: config.Default(), Path: filepath.Join(tempDir, "config", "savitar.local.json")}
	loaded.Config.Knowledge.RepoMarkdownDirs = []string{"docs/hackathon"}
	rt := savitarruntime.NewAtRoot(loaded, tempDir)
	if _, err := rt.InitSession(); err != nil {
		t.Fatalf("InitSession returned error: %v", err)
	}

	service := Service{
		Runtime:       rt,
		Operator:      &operator.Service{Store: NewRuntimeStore(rt)},
		Authenticator: DemoAuthenticator{},
		DemoMode:      true,
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()
	service.Routes().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d with body %q", recorder.Code, recorder.Body.String())
	}
	body := recorder.Body.String()
	for _, want := range []string{"Savitar Operator Console", "Gemma 4 Good prototype", "Tier 2: Standard Operator Host", "local-demo"} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected dashboard HTML to contain %q", want)
		}
	}
}

func TestDashboardAPIFallsBackToRepoSources(t *testing.T) {
	tempDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tempDir, "docs", "hackathon"), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "docs", "hackathon", "submission-strategy.md"), []byte("# Strategy\n\nMac Mini M4 and Pixel Fold.\n"), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	loaded := config.Loaded{Config: config.Default(), Path: filepath.Join(tempDir, "config", "savitar.local.json")}
	loaded.Config.Knowledge.RepoMarkdownDirs = []string{"docs/hackathon"}
	rt := savitarruntime.NewAtRoot(loaded, tempDir)

	service := Service{
		Runtime:       rt,
		Operator:      &operator.Service{Store: NewRuntimeStore(rt)},
		Authenticator: DemoAuthenticator{},
		DemoMode:      true,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	recorder := httptest.NewRecorder()
	service.Routes().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d with body %q", recorder.Code, recorder.Body.String())
	}

	var payload DashboardData
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	if payload.KnowledgeMode != "repo-markdown-fallback" {
		t.Fatalf("expected repo markdown fallback mode, got %q", payload.KnowledgeMode)
	}
	if len(payload.RepoSources) == 0 {
		t.Fatalf("expected repo sources in dashboard payload")
	}
	if payload.AuthMode != "local-demo" {
		t.Fatalf("expected local demo auth mode, got %q", payload.AuthMode)
	}
}

func TestRoutesRejectUnauthorizedRequests(t *testing.T) {
	loaded := config.Loaded{Config: config.Default(), Path: "config/savitar.local.json"}
	rt := savitarruntime.NewAtRoot(loaded, t.TempDir())
	service := Service{
		Runtime:       rt,
		Operator:      &operator.Service{Store: NewRuntimeStore(rt)},
		Authenticator: deniedAuthenticator{},
	}

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	recorder := httptest.NewRecorder()
	service.Routes().ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}
