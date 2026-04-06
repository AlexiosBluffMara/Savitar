package respond

import (
	"context"
	"strings"
	"testing"

	"github.com/alexiosbluffmara/savitar/internal/chat"
	"github.com/alexiosbluffmara/savitar/internal/config"
	"github.com/alexiosbluffmara/savitar/internal/gateway"
	savitarruntime "github.com/alexiosbluffmara/savitar/internal/runtime"
)

type stubGenerator struct {
	replies  []string
	requests []chat.Request
	err      error
}

func (s *stubGenerator) Generate(_ context.Context, request chat.Request) (string, error) {
	s.requests = append(s.requests, request)
	if s.err != nil {
		return "", s.err
	}
	if len(s.replies) == 0 {
		return "", nil
	}
	reply := s.replies[0]
	s.replies = s.replies[1:]
	return reply, nil
}

func testRuntime(root string) *savitarruntime.Runtime {
	return savitarruntime.NewAtRoot(config.Loaded{
		Config: config.Default(),
		Path:   "config/savitar.local.json",
	}, root)
}

func TestReplyReturnsHelpForEmptyInput(t *testing.T) {
	engine := NewEngine(testRuntime(t.TempDir()))
	reply, err := engine.Reply(context.Background(), gateway.Envelope{Surface: gateway.SurfaceDiscord})
	if err != nil {
		t.Fatalf("Reply returned error: %v", err)
	}

	if !strings.Contains(reply, "Public commands") {
		t.Fatalf("expected help text, got %q", reply)
	}
}

func TestReplyStatusIncludesAgentName(t *testing.T) {
	generator := &stubGenerator{}
	engine := NewEngineWithGenerator(testRuntime(t.TempDir()), generator)
	reply, err := engine.Reply(context.Background(), gateway.Envelope{
		Surface: gateway.SurfaceDiscord,
		Body:    "status",
		Metadata: map[string]string{
			"mode": "preview",
		},
	})
	if err != nil {
		t.Fatalf("Reply returned error: %v", err)
	}

	if !strings.Contains(reply, "agent: Savitar") {
		t.Fatalf("expected status reply to include agent name, got %q", reply)
	}
	if len(generator.requests) != 0 {
		t.Fatal("expected operator command to bypass model generator")
	}
}

func TestDefaultReplyPointsBackToCommandsWhenModelUnavailable(t *testing.T) {
	engine := NewEngine(testRuntime(t.TempDir()))
	reply, err := engine.Reply(context.Background(), gateway.Envelope{Surface: gateway.SurfaceDiscord, Body: "please summarize the repo"})
	if err != nil {
		t.Fatalf("Reply returned error: %v", err)
	}

	if !strings.Contains(reply, "model-backed reply path is unavailable") {
		t.Fatalf("expected unavailable-model reply, got %q", reply)
	}
	if !strings.Contains(reply, "help") {
		t.Fatalf("expected default reply to mention commands, got %q", reply)
	}
}

func TestDefaultReplyUsesModelGenerator(t *testing.T) {
	generator := &stubGenerator{replies: []string{"Model-backed answer."}}
	engine := NewEngineWithGenerator(testRuntime(t.TempDir()), generator)
	reply, err := engine.Reply(context.Background(), gateway.Envelope{
		Surface:           gateway.SurfaceDiscord,
		ConversationID:    "conversation-1",
		SenderID:          "user-1",
		SenderDisplayName: "tester",
		Body:              "please summarize the repo",
	})
	if err != nil {
		t.Fatalf("Reply returned error: %v", err)
	}
	if reply != "Model-backed answer." {
		t.Fatalf("unexpected model-backed reply: %q", reply)
	}
	if len(generator.requests) != 1 {
		t.Fatalf("expected one model request, got %d", len(generator.requests))
	}
}

func TestReplyCarriesConversationHistoryAcrossTurns(t *testing.T) {
	generator := &stubGenerator{replies: []string{"First reply.", "Second reply."}}
	engine := NewEngineWithGenerator(testRuntime(t.TempDir()), generator)

	firstEnv := gateway.Envelope{
		Surface:        gateway.SurfaceDiscord,
		ConversationID: "conversation-2",
		SenderID:       "user-1",
		Body:           "hello there",
	}
	if _, err := engine.Reply(context.Background(), firstEnv); err != nil {
		t.Fatalf("first Reply returned error: %v", err)
	}

	secondEnv := gateway.Envelope{
		Surface:        gateway.SurfaceDiscord,
		ConversationID: "conversation-2",
		SenderID:       "user-1",
		Body:           "what did I just ask you?",
	}
	if _, err := engine.Reply(context.Background(), secondEnv); err != nil {
		t.Fatalf("second Reply returned error: %v", err)
	}

	if len(generator.requests) != 2 {
		t.Fatalf("expected two model requests, got %d", len(generator.requests))
	}
	history := generator.requests[1].History
	if len(history) != 2 {
		t.Fatalf("expected prior user and assistant turns, got %d entries", len(history))
	}
	if history[0].Role != "user" || history[0].Content != "hello there" {
		t.Fatalf("unexpected first history turn: %+v", history[0])
	}
	if history[1].Role != "assistant" || history[1].Content != "First reply." {
		t.Fatalf("unexpected second history turn: %+v", history[1])
	}
}

func TestReplyDoesNotShareHistoryAcrossGuildUsers(t *testing.T) {
	generator := &stubGenerator{replies: []string{"User one.", "User two."}}
	engine := NewEngineWithGenerator(testRuntime(t.TempDir()), generator)

	firstEnv := gateway.Envelope{
		Surface:        gateway.SurfaceDiscord,
		ConversationID: "shared-channel",
		SenderID:       "user-1",
		Body:           "hello there",
	}
	if _, err := engine.Reply(context.Background(), firstEnv); err != nil {
		t.Fatalf("first Reply returned error: %v", err)
	}

	secondEnv := gateway.Envelope{
		Surface:        gateway.SurfaceDiscord,
		ConversationID: "shared-channel",
		SenderID:       "user-2",
		Body:           "what did they say?",
	}
	if _, err := engine.Reply(context.Background(), secondEnv); err != nil {
		t.Fatalf("second Reply returned error: %v", err)
	}

	if len(generator.requests) != 2 {
		t.Fatalf("expected two model requests, got %d", len(generator.requests))
	}
	if len(generator.requests[1].History) != 0 {
		t.Fatalf("expected shared guild channel history to be isolated per sender, got %d entries", len(generator.requests[1].History))
	}
}

func TestOperatorCommandDeniedForNonAllowlistedUser(t *testing.T) {
	engine := NewEngine(testRuntime(t.TempDir()))
	reply, err := engine.Reply(context.Background(), gateway.Envelope{Surface: gateway.SurfaceDiscord, SenderID: "user-1", Body: "status"})
	if err != nil {
		t.Fatalf("Reply returned error: %v", err)
	}

	if !strings.Contains(reply, "not allowlisted") {
		t.Fatalf("expected operator denial text, got %q", reply)
	}
}
