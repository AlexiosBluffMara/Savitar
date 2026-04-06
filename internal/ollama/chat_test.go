package ollama

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alexiosbluffmara/savitar/internal/chat"
	"github.com/alexiosbluffmara/savitar/internal/config"
	"github.com/alexiosbluffmara/savitar/internal/models"
)

func TestGenerateSendsChatRequestWithBearerToken(t *testing.T) {
	t.Setenv("OLLAMA_API_KEY", "test-token")

	var seen chatRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("unexpected authorization header %q", got)
		}
		if err := json.NewDecoder(r.Body).Decode(&seen); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		_ = json.NewEncoder(w).Encode(chatResponse{Message: chatMessage{Role: "assistant", Content: "hello back"}})
	}))
	defer server.Close()

	cfg := config.Default()
	cfg.Integrations.Ollama.Enabled = true
	cfg.Integrations.Ollama.BaseURL = server.URL + "/api"

	generator := NewReplyGenerator(cfg)
	reply, err := generator.Generate(context.Background(), chat.Request{
		Surface:            "discord",
		UserInput:          "Hello",
		ReplyLimit:         500,
		AllowCloudFallback: true,
		Route: models.Decision{Profile: models.Profile{
			Name:     "copilot-0x",
			Provider: "copilot",
		}},
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if reply != "hello back" {
		t.Fatalf("unexpected reply %q", reply)
	}
	if seen.Model != cfg.Integrations.Ollama.CloudModel {
		t.Fatalf("unexpected model %q", seen.Model)
	}
	if seen.Stream {
		t.Fatal("expected non-streaming chat request")
	}
	if len(seen.Messages) < 2 {
		t.Fatalf("expected system and user messages, got %d", len(seen.Messages))
	}
	if seen.Messages[0].Role != "system" {
		t.Fatalf("expected first message to be system, got %q", seen.Messages[0].Role)
	}
	if !strings.Contains(seen.Messages[0].Content, "software operator assistant") {
		t.Fatalf("expected guardrail system prompt, got %q", seen.Messages[0].Content)
	}
}

func TestGenerateFallsBackFromLocalToCloud(t *testing.T) {
	t.Setenv("OLLAMA_API_KEY", "test-token")

	localCalls := 0
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		localCalls++
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(errorResponse{Error: "local model unavailable"})
	}))
	defer localServer.Close()

	cloudCalls := 0
	cloudServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cloudCalls++
		_ = json.NewEncoder(w).Encode(chatResponse{Message: chatMessage{Role: "assistant", Content: "cloud fallback"}})
	}))
	defer cloudServer.Close()

	cfg := config.Default()
	cfg.Models.LocalDefault.Provider = "ollama"
	cfg.Models.LocalDefault.Endpoint = localServer.URL
	cfg.Models.LocalDefault.Model = "local-model"
	cfg.Integrations.Ollama.Enabled = true
	cfg.Integrations.Ollama.BaseURL = cloudServer.URL + "/api"

	generator := NewReplyGenerator(cfg)
	reply, err := generator.Generate(context.Background(), chat.Request{
		Surface:            "discord",
		UserInput:          "Hi",
		AllowCloudFallback: true,
		Route: models.Decision{Profile: models.Profile{
			Name:     "local-default",
			Provider: "ollama",
		}},
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if reply != "cloud fallback" {
		t.Fatalf("unexpected reply %q", reply)
	}
	if localCalls != 1 {
		t.Fatalf("expected one local attempt, got %d", localCalls)
	}
	if cloudCalls != 1 {
		t.Fatalf("expected one cloud fallback attempt, got %d", cloudCalls)
	}
}

func TestGenerateFailsClosedForPrivateContextWithoutLocalTarget(t *testing.T) {
	t.Setenv("OLLAMA_API_KEY", "test-token")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected no cloud request for private context without explicit fallback")
	}))
	defer server.Close()

	cfg := config.Default()
	cfg.Integrations.Ollama.Enabled = true
	cfg.Integrations.Ollama.BaseURL = server.URL + "/api"

	generator := NewReplyGenerator(cfg)
	_, err := generator.Generate(context.Background(), chat.Request{
		Surface: "discord",
		Task: models.Task{
			PrivateContext: true,
		},
		Route: models.Decision{Profile: models.Profile{
			Name:     "local-default",
			Provider: "mlx",
		}},
	})
	if err == nil {
		t.Fatal("expected private-context request to fail closed without a local target")
	}
}

func TestGenerateFailsClosedForGuildWithoutCloudOptIn(t *testing.T) {
	t.Setenv("OLLAMA_API_KEY", "test-token")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected no cloud request when guild cloud replies are not allowed")
	}))
	defer server.Close()

	cfg := config.Default()
	cfg.Integrations.Ollama.Enabled = true
	cfg.Integrations.Ollama.BaseURL = server.URL + "/api"

	generator := NewReplyGenerator(cfg)
	_, err := generator.Generate(context.Background(), chat.Request{
		Surface:            "discord",
		UserInput:          "hello",
		AllowCloudFallback: false,
		Route: models.Decision{Profile: models.Profile{
			Name:     "copilot-0x",
			Provider: "copilot",
		}},
	})
	if err == nil {
		t.Fatal("expected guild request to fail closed without explicit cloud opt-in")
	}
}

func TestGeneratePrefersLocalBeforeCloudForGuildReplies(t *testing.T) {
	t.Setenv("OLLAMA_API_KEY", "test-token")

	localCalls := 0
	localServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		localCalls++
		_ = json.NewEncoder(w).Encode(chatResponse{Message: chatMessage{Role: "assistant", Content: "local reply"}})
	}))
	defer localServer.Close()

	cloudCalls := 0
	cloudServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cloudCalls++
		_ = json.NewEncoder(w).Encode(chatResponse{Message: chatMessage{Role: "assistant", Content: "cloud reply"}})
	}))
	defer cloudServer.Close()

	cfg := config.Default()
	cfg.Models.LocalDefault.Provider = "ollama"
	cfg.Models.LocalDefault.Endpoint = localServer.URL
	cfg.Models.LocalDefault.Model = "local-model"
	cfg.Integrations.Ollama.Enabled = true
	cfg.Integrations.Ollama.BaseURL = cloudServer.URL + "/api"

	generator := NewReplyGenerator(cfg)
	reply, err := generator.Generate(context.Background(), chat.Request{
		Surface:            "discord",
		UserInput:          "hello",
		AllowCloudFallback: true,
		Route: models.Decision{Profile: models.Profile{
			Name:     "copilot-0x",
			Provider: "copilot",
		}},
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if reply != "local reply" {
		t.Fatalf("unexpected reply %q", reply)
	}
	if localCalls != 1 {
		t.Fatalf("expected one local call, got %d", localCalls)
	}
	if cloudCalls != 0 {
		t.Fatalf("expected no cloud call, got %d", cloudCalls)
	}
}
