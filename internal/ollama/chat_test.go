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

func TestGenerateSendsLocalChatRequest(t *testing.T) {
	var seen chatRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/chat" {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "" {
			t.Fatalf("unexpected authorization header %q", got)
		}
		if err := json.NewDecoder(r.Body).Decode(&seen); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		_ = json.NewEncoder(w).Encode(chatResponse{Message: chatMessage{Role: "assistant", Content: "hello back"}})
	}))
	defer server.Close()

	cfg := config.Default()
	cfg.Models.LocalDefault.Provider = "ollama"
	cfg.Models.LocalDefault.Endpoint = server.URL
	cfg.Models.LocalDefault.Model = "local-model"

	generator := NewReplyGenerator(cfg)
	reply, err := generator.Generate(context.Background(), chat.Request{
		Surface:    "discord",
		UserInput:  "Hello",
		ReplyLimit: 500,
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
	if seen.Model != cfg.Models.LocalDefault.Model {
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

func TestGenerateReturnsLocalError(t *testing.T) {
	cfg := config.Default()
	cfg.Models.LocalDefault.Provider = "ollama"
	cfg.Models.LocalDefault.Endpoint = ""
	cfg.Models.LocalDefault.Model = "local-model"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(errorResponse{Error: "local model unavailable"})
	}))
	defer server.Close()

	cfg.Models.LocalDefault.Endpoint = server.URL
	cfg.Models.LocalDefault.Model = "local-model"

	generator := NewReplyGenerator(cfg)
	_, err := generator.Generate(context.Background(), chat.Request{
		Surface:   "discord",
		UserInput: "Hi",
		Route: models.Decision{Profile: models.Profile{
			Name:     "local-default",
			Provider: "ollama",
		}},
	})
	if err == nil {
		t.Fatal("expected local ollama error")
	}
	if !strings.Contains(err.Error(), "local model unavailable") {
		t.Fatalf("expected local error in %q", err)
	}
}

func TestGenerateFailsClosedWithoutLocalTarget(t *testing.T) {
	cfg := config.Default()
	cfg.Models.LocalDefault.Provider = "mlx"

	generator := NewReplyGenerator(cfg)
	_, err := generator.Generate(context.Background(), chat.Request{
		Surface:   "discord",
		UserInput: "hello",
		Route: models.Decision{Profile: models.Profile{
			Name:     "local-default",
			Provider: "mlx",
		}},
	})
	if err == nil {
		t.Fatal("expected request to fail closed without a local Ollama target")
	}
}
