package respond

import (
	"context"
	"strings"
	"testing"

	"github.com/alexiosbluffmara/savitar/internal/chat"
	"github.com/alexiosbluffmara/savitar/internal/gateway"
)

type stubRepoSearcher struct {
	result repoEvidence
	err    error
}

func (s stubRepoSearcher) Search(_ context.Context, _ string) (repoEvidence, error) {
	return s.result, s.err
}

type stubLiveLookup struct {
	result liveEvidence
	err    error
}

func (s stubLiveLookup) Lookup(_ context.Context, _ string) (liveEvidence, error) {
	return s.result, s.err
}

func TestModelReplyInjectsRepoAndLiveEvidence(t *testing.T) {
	generator := &stubGenerator{replies: []string{"Here is the answer."}}
	engine := NewEngineWithGenerator(testRuntime(t.TempDir()), generator)
	engine.repoSearch = stubRepoSearcher{result: repoEvidence{
		Context: "Local repository markdown evidence:\n- docs/adr/0006 :: Decision",
		Sources: []string{"Local repo: docs/adr/0006 :: Decision"},
	}}
	engine.liveLookup = stubLiveLookup{result: liveEvidence{
		ToolContext: chat.ToolContext{ServerName: "web", ToolName: "duckduckgo-json", Result: "Live web lookup (https://example.com): Example title\nExample snippet"},
	}}

	reply, err := engine.Reply(context.Background(), gateway.Envelope{
		Surface:           gateway.SurfaceDiscord,
		ConversationID:    "conversation-live",
		SenderID:          "user-1",
		SenderDisplayName: "tester",
		Body:              "Who is Andrej Karpathy and how does this rewrite work?",
		Metadata:          map[string]string{"mode": "preview", "dm": "true"},
	})
	if err != nil {
		t.Fatalf("Reply returned error: %v", err)
	}
	if len(generator.requests) != 1 {
		t.Fatalf("expected one generator request, got %d", len(generator.requests))
	}
	request := generator.requests[0]
	if len(request.MemoryContext) == 0 || !strings.Contains(request.MemoryContext[len(request.MemoryContext)-1], "Local repository markdown evidence") {
		t.Fatalf("expected repo context in memory context, got %#v", request.MemoryContext)
	}
	if len(request.ToolContexts) != 1 {
		t.Fatalf("expected one live tool context, got %d", len(request.ToolContexts))
	}
	if !strings.Contains(reply, "Sources:") {
		t.Fatalf("expected reply to include sources footer, got %q", reply)
	}
	if !strings.Contains(reply, "Local repo:") {
		t.Fatalf("expected reply to include local repo source, got %q", reply)
	}
	if !strings.Contains(reply, "Live web lookup") {
		t.Fatalf("expected reply to include live source, got %q", reply)
	}
}

func TestParseDuckDuckGoResult(t *testing.T) {
	body := `<a rel="nofollow" class="result__a" href="//duckduckgo.com/l/?uddg=https%3A%2F%2Fgithub.com%2Fkarpathy%2Freader3">GitHub - karpathy/reader3</a><a class="result__snippet">Quick illustration of how one can easily read books together with LLMs.</a>`
	result, err := parseDuckDuckGoResult(body)
	if err != nil {
		t.Fatalf("parseDuckDuckGoResult returned error: %v", err)
	}
	if result.URL != "https://github.com/karpathy/reader3" {
		t.Fatalf("unexpected result URL: %q", result.URL)
	}
	if !strings.Contains(result.Snippet, "read books together with LLMs") {
		t.Fatalf("unexpected snippet: %q", result.Snippet)
	}
}

func TestNormalizeLiveLookupQuery(t *testing.T) {
	query := normalizeLiveLookupQuery("Who is Andrej Karpathy and how does the rewrite work?")
	if query != "Andrej Karpathy" {
		t.Fatalf("unexpected normalized query: %q", query)
	}
}

func TestNormalizeRepoSearchQuery(t *testing.T) {
	query := normalizeRepoSearchQuery("Who is Andrej Karpathy and how does the rewrite work?")
	if query != "the rewrite work" {
		t.Fatalf("unexpected normalized repo query: %q", query)
	}
}

func TestShouldUseLiveLookupRequiresQuestionOrLookupVerb(t *testing.T) {
	if shouldUseLiveLookup("please summarize the repo for me") {
		t.Fatal("expected plain conversational text not to trigger live lookup")
	}
	if !shouldUseLiveLookup("Who is Andrej Karpathy?") {
		t.Fatal("expected question to trigger live lookup")
	}
	if !shouldUseLiveLookup("look up Andrej Karpathy") {
		t.Fatal("expected explicit lookup command to trigger live lookup")
	}
}

func TestModelReplySkipsLiveLookupInDiscordDMByDefault(t *testing.T) {
	generator := &stubGenerator{replies: []string{"Here is the answer."}}
	engine := NewEngineWithGenerator(testRuntime(t.TempDir()), generator)
	engine.liveLookup = stubLiveLookup{result: liveEvidence{
		ToolContext: chat.ToolContext{ServerName: "web", ToolName: "duckduckgo-json", Result: "Live web lookup (https://example.com): Example title\nExample snippet"},
	}}

	reply, err := engine.Reply(context.Background(), gateway.Envelope{
		Surface:           gateway.SurfaceDiscord,
		ConversationID:    "conversation-live-private",
		SenderID:          "user-1",
		SenderDisplayName: "tester",
		Body:              "Who is Andrej Karpathy and how does this rewrite work?",
		Metadata:          map[string]string{"dm": "true"},
	})
	if err != nil {
		t.Fatalf("Reply returned error: %v", err)
	}
	if len(generator.requests) != 1 {
		t.Fatalf("expected one generator request, got %d", len(generator.requests))
	}
	if len(generator.requests[0].ToolContexts) != 0 {
		t.Fatalf("expected DM live lookup to be disabled by default, got %#v", generator.requests[0].ToolContexts)
	}
	if strings.Contains(reply, "Live web lookup") {
		t.Fatalf("expected reply not to include live web source, got %q", reply)
	}
	if strings.Contains(reply, "Sources:") && strings.Contains(reply, "Live web:") {
		t.Fatalf("expected no live web footer in DM by default, got %q", reply)
	}
	if strings.Contains(reply, "Local repo:") {
		t.Fatalf("expected no repo source footer without repo evidence, got %q", reply)
	}
	if reply != "Here is the answer." {
		t.Fatalf("unexpected reply: %q", reply)
	}
}

func TestModelReplyFallsBackToGroundedEvidenceWhenGeneratorFails(t *testing.T) {
	generator := &stubGenerator{err: context.DeadlineExceeded}
	engine := NewEngineWithGenerator(testRuntime(t.TempDir()), generator)
	engine.repoSearch = stubRepoSearcher{result: repoEvidence{
		Context: "Local repository markdown evidence:\n- docs/roadmap/0006 :: Discord plus authenticated operator console.",
		Sources: []string{"Local repo: docs/roadmap/0006"},
	}}
	engine.liveLookup = stubLiveLookup{result: liveEvidence{
		ToolContext: chat.ToolContext{ServerName: "web", ToolName: "duckduckgo-json", Result: "Live web lookup (https://github.com/karpathy/reader3): GitHub - karpathy/reader3\nQuick illustration of how one can easily read books together with LLMs."},
	}}

	reply, err := engine.Reply(context.Background(), gateway.Envelope{
		Surface:        gateway.SurfaceDiscord,
		ConversationID: "fallback-conversation",
		SenderID:       "user-1",
		Body:           "Who is Andrej Karpathy and how does the rewrite work?",
		Metadata:       map[string]string{"mode": "preview", "dm": "true"},
	})
	if err != nil {
		t.Fatalf("expected grounded fallback reply, got error: %v", err)
	}
	if !strings.Contains(reply, "grounded context") {
		t.Fatalf("expected fallback explanation, got %q", reply)
	}
	if !strings.Contains(reply, "Local repository:") {
		t.Fatalf("expected local repo evidence in fallback, got %q", reply)
	}
	if !strings.Contains(reply, "Live web:") {
		t.Fatalf("expected live web evidence in fallback, got %q", reply)
	}
	if !strings.Contains(reply, "Sources:") {
		t.Fatalf("expected sources footer in fallback, got %q", reply)
	}
}