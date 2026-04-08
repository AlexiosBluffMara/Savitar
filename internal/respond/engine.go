package respond

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/alexiosbluffmara/savitar/internal/chat"
	"github.com/alexiosbluffmara/savitar/internal/gateway"
	"github.com/alexiosbluffmara/savitar/internal/models"
	"github.com/alexiosbluffmara/savitar/internal/ollama"
	savitarruntime "github.com/alexiosbluffmara/savitar/internal/runtime"
)

type Engine struct {
	rt         *savitarruntime.Runtime
	generator  chat.Generator
	history    *conversationHistory
	repoSearch repoSearcher
	liveLookup liveLookup
}

func NewEngine(rt *savitarruntime.Runtime) Engine {
	var generator chat.Generator
	if rt != nil {
		candidate := ollama.NewReplyGenerator(rt.Config().Config)
		if candidate.Available() {
			generator = candidate
		}
	}
	return NewEngineWithGenerator(rt, generator)
}

func NewEngineWithGenerator(rt *savitarruntime.Runtime, generator chat.Generator) Engine {
	engine := Engine{
		rt:         rt,
		generator:  generator,
		history:    newConversationHistory(8),
		repoSearch: markdownRepoSearcher{rt: rt},
		liveLookup: newLiveLookup(rt),
	}
	return engine
}

func (e Engine) Reply(ctx context.Context, env gateway.Envelope) (string, error) {
	if e.rt == nil {
		return "Savitar is missing runtime context for this reply.", nil
	}

	input := cleanInput(env.Body, e.rt.Persona().Name)
	if input == "" {
		return e.helpText(), nil
	}

	command := firstWord(input)
	if requiresOperator(command) && !e.isOperator(env) {
		return e.operatorDeniedText(), nil
	}

	switch command {
	case "commands", "help":
		return e.helpText(), nil
	case "ping":
		return fmt.Sprintf("%s is online.", e.rt.Persona().Name), nil
	case "status", "health":
		return e.statusText()
	case "doctor":
		return e.doctorText(), nil
	case "plan":
		return e.planText(), nil
	case "models":
		return e.modelsText(), nil
	case "contracts":
		return e.contractsText(), nil
	case "gateway":
		return e.gatewayText(), nil
	case "persona":
		return e.personaText(), nil
	case "session":
		return e.sessionText()
	default:
		reply, err := e.modelReply(ctx, env, input)
		if err != nil {
			if env.Metadata["mode"] == "preview" {
				return "", err
			}
			return e.modelUnavailableText(input), nil
		}
		return reply, nil
	}
}

func (e Engine) helpText() string {
	name := e.rt.Persona().Name
	return strings.Join([]string{
		fmt.Sprintf("%s Discord conversational mode is live.", name),
		"Public commands: help and ping still work directly.",
		"Normal messages use the model-backed reply path when an Ollama target is configured.",
		"Operator-only commands for allowlisted Discord user IDs: status, doctor, plan, models, contracts, gateway, persona, session.",
		"DM the bot directly or mention it with one of those commands.",
		"The transport replies immediately, then edits the placeholder once the response is ready.",
	}, "\n")
}

func (e Engine) statusText() (string, error) {
	status, err := e.rt.Status()
	if err != nil {
		return "", err
	}

	mode := "defaults"
	if status.ConfigLoaded {
		mode = "loaded"
	}

	return strings.Join([]string{
		"Savitar status:",
		fmt.Sprintf("- agent: %s", status.AgentName),
		fmt.Sprintf("- config: %s (%s)", status.ConfigPath, mode),
		fmt.Sprintf("- model profiles: %d", status.ModelProfiles),
		fmt.Sprintf("- enabled MCP servers: %d", status.EnabledMCPServers),
		fmt.Sprintf("- workspace agents: %d", status.AgentCount),
		fmt.Sprintf("- workspace skills: %d", status.SkillCount),
		fmt.Sprintf("- public surfaces: %d/%d", status.EnabledSurfaces, status.TotalSurfaces),
		fmt.Sprintf("- session initialized: %t", status.SessionInitialized),
	}, "\n"), nil
}

func (e Engine) doctorText() string {
	lines := []string{"Savitar doctor:"}
	for _, check := range e.rt.Doctor() {
		status := "ok"
		if !check.Available && check.Required {
			status = "missing"
		} else if !check.Available {
			status = "optional-missing"
		}
		lines = append(lines, fmt.Sprintf("- %s: %s (%s)", check.Name, status, check.Details))
	}

	return strings.Join(lines, "\n")
}

func (e Engine) planText() string {
	steps := e.rt.Plan()
	lines := []string{"Savitar plan:"}
	for index, step := range steps {
		lines = append(lines, fmt.Sprintf("%d. %s", index+1, step))
	}

	return strings.Join(lines, "\n")
}

func (e Engine) modelsText() string {
	profiles := e.rt.ModelProfiles()
	lines := []string{"Savitar model profiles:"}
	for _, profile := range profiles {
		lines = append(lines, fmt.Sprintf("- %s: %s %s (%.2fx)", profile.Name, profile.Provider, profile.Model, profile.UsageMultiplier))
	}
	decision := e.rt.Router().Route(models.Task{Complexity: models.ComplexityComplex})
	lines = append(lines, fmt.Sprintf("Example complex-task route: %s (%s)", decision.Profile.Name, decision.Reason))
	return strings.Join(lines, "\n")
}

func (e Engine) contractsText() string {
	lines := []string{"Savitar contracts:"}
	for _, contract := range e.rt.Contracts() {
		lines = append(lines, fmt.Sprintf("- %s [%s]: %s", contract.Name, contract.Stage, contract.Goal))
	}
	return strings.Join(lines, "\n")
}

func (e Engine) gatewayText() string {
	lines := []string{"Savitar surfaces:"}
	for _, surface := range e.rt.GatewayPlan().Surfaces {
		lines = append(lines, fmt.Sprintf("- %s: enabled=%t kind=%s identity=%s details=%s", surface.Name, surface.Enabled, surface.Kind, surface.Identity, surface.Details))
	}
	return strings.Join(lines, "\n")
}

func (e Engine) personaText() string {
	profile := e.rt.Persona()
	return strings.Join([]string{
		fmt.Sprintf("%s persona:", profile.Name),
		fmt.Sprintf("- style: %s", profile.Style),
		fmt.Sprintf("- tone: %s", profile.Tone),
		fmt.Sprintf("- design bias: %s", profile.DesignBias),
		fmt.Sprintf("- commentary density: %s", profile.CommentaryDensity),
		fmt.Sprintf("- public bio: %s", profile.PublicBio),
		fmt.Sprintf("- disclosure policy: %s", profile.DisclosurePolicy),
	}, "\n")
}

func (e Engine) sessionText() (string, error) {
	report, err := e.rt.Session()
	if err != nil {
		return "", err
	}

	return strings.Join([]string{
		"Savitar session:",
		fmt.Sprintf("- path: %s", report.Path),
		fmt.Sprintf("- initialized: %t", report.Exists),
		fmt.Sprintf("- surface: %s", blankIfEmpty(report.State.CurrentSurface)),
		fmt.Sprintf("- model profile: %s", blankIfEmpty(report.State.CurrentModelProfile)),
		fmt.Sprintf("- conversation: %s", blankIfEmpty(report.State.ActiveConversationID)),
		fmt.Sprintf("- last command: %s", blankIfEmpty(report.State.LastCommand)),
	}, "\n"), nil
}

func (e Engine) modelUnavailableText(input string) string {
	name := e.rt.Persona().Name
	return strings.Join([]string{
		fmt.Sprintf("%s saw your message about: %q", name, trimSnippet(input, 96)),
		"The model-backed reply path is unavailable right now.",
		"Public commands are help and ping. Operator diagnostics require an allowlisted Discord user ID.",
	}, "\n")
}

func (e Engine) operatorDeniedText() string {
	return strings.Join([]string{
		"This Discord user is not allowlisted for operator diagnostics.",
		"Public commands remain available: help and ping.",
		"Add the sender ID to transports.discord.operatorUserIDs in local config if remote operator access is intended.",
	}, "\n")
}

func cleanInput(body string, agentName string) string {
	cleaned := strings.TrimSpace(body)
	if cleaned == "" {
		return ""
	}

	prefixes := []string{
		"@" + agentName,
		agentName,
	}
	for _, prefix := range prefixes {
		if prefix == "" {
			continue
		}
		if strings.HasPrefix(strings.ToLower(cleaned), strings.ToLower(prefix)) {
			cleaned = strings.TrimSpace(cleaned[len(prefix):])
			cleaned = strings.TrimLeft(cleaned, ":,- ")
		}
	}

	return strings.TrimSpace(cleaned)
}

func firstWord(input string) string {
	fields := strings.Fields(strings.ToLower(input))
	if len(fields) == 0 {
		return ""
	}
	return strings.Trim(fields[0], ".,:;!?")
}

func trimSnippet(input string, limit int) string {
	cleaned := strings.Join(strings.Fields(input), " ")
	runes := []rune(cleaned)
	if len(runes) <= limit {
		return cleaned
	}
	if limit <= 3 {
		return string(runes[:limit])
	}
	return string(runes[:limit-3]) + "..."
}

func blankIfEmpty(value string) string {
	if value == "" {
		return "<none>"
	}
	return value
}

func (e Engine) isOperator(env gateway.Envelope) bool {
	if env.Metadata["mode"] == "preview" {
		return true
	}
	if env.Surface != gateway.SurfaceDiscord {
		return true
	}
	for _, userID := range e.rt.Config().Config.Transports.Discord.OperatorUserIDs {
		if userID == env.SenderID {
			return true
		}
	}
	return false
}

func requiresOperator(command string) bool {
	switch command {
	case "status", "health", "doctor", "plan", "models", "contracts", "gateway", "persona", "session":
		return true
	default:
		return false
	}
}

func (e Engine) modelReply(ctx context.Context, env gateway.Envelope, input string) (string, error) {
	if e.generator == nil {
		return "", fmt.Errorf("no model-backed reply generator is configured")
	}

	key := conversationKey(env)
	task := inferTask(input, env)
	memoryContext := e.loadMemoryContext(input)
	repoEvidence, _ := e.loadRepoEvidence(ctx, input)
	if repoEvidence.Context != "" {
		memoryContext = append(memoryContext, repoEvidence.Context)
	}
	liveEvidence, _ := e.loadLiveEvidence(ctx, env, input)
	request := chat.Request{
		Surface:           string(env.Surface),
		ConversationID:    key,
		SenderDisplayName: env.SenderDisplayName,
		UserInput:         input,
		History:           e.history.load(key),
		Task:              task,
		Route:             e.rt.Router().Route(task),
		ReplyLimit:        e.replyLimit(env),
		MemoryContext:     memoryContext,
		ToolContexts:      liveEvidence,
	}

	reply, err := e.generator.Generate(ctx, request)
	if err != nil {
		fallback := e.evidenceFallbackReply(input, repoEvidence, liveEvidence)
		if fallback != "" {
			return e.attachSources(fallback, append(repoEvidence.Sources, collectToolSources(liveEvidence)...), e.replyLimit(env)), nil
		}
		return "", err
	}
	reply = strings.TrimSpace(reply)
	if reply == "" {
		return "", fmt.Errorf("empty model-backed reply")
	}

	e.history.append(key,
		chat.Turn{Role: "user", Content: input},
		chat.Turn{Role: "assistant", Content: reply},
	)
	return e.attachSources(reply, append(repoEvidence.Sources, collectToolSources(liveEvidence)...), e.replyLimit(env)), nil
}

func (e Engine) evidenceFallbackReply(input string, repo repoEvidence, live []chat.ToolContext) string {
	parts := []string{"The model-backed reply path is unavailable right now, but I found grounded context."}
	if strings.TrimSpace(repo.Context) != "" {
		parts = append(parts, "Local repository:")
		parts = append(parts, trimSnippet(repo.Context, 320))
	}
	for _, toolContext := range live {
		if strings.TrimSpace(toolContext.Result) == "" {
			continue
		}
		parts = append(parts, "Live web:")
		parts = append(parts, trimSnippet(toolContext.Result, 280))
		break
	}
	if len(parts) == 1 {
		return ""
	}
	parts = append(parts, fmt.Sprintf("Original request: %q", trimSnippet(input, 120)))
	return strings.Join(parts, "\n")
}

func (e Engine) loadRepoEvidence(ctx context.Context, input string) (repoEvidence, error) {
	if e.repoSearch == nil {
		return repoEvidence{}, nil
	}
	return e.repoSearch.Search(ctx, input)
}

func (e Engine) loadLiveEvidence(ctx context.Context, env gateway.Envelope, input string) ([]chat.ToolContext, error) {
	if e.liveLookup == nil || !e.allowLiveWebLookup(env) || !shouldUseLiveLookup(input) {
		return nil, nil
	}
	result, err := e.liveLookup.Lookup(ctx, input)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(result.ToolContext.Result) == "" {
		return nil, nil
	}
	return []chat.ToolContext{result.ToolContext}, nil
}

func (e Engine) allowLiveWebLookup(env gateway.Envelope) bool {
	if e.rt == nil {
		return false
	}
	if env.Metadata["mode"] == "preview" {
		return true
	}
	cfg := e.rt.Config().Config.Transports.Discord
	if env.Surface == gateway.SurfaceDiscord {
		if env.Metadata["dm"] == "true" {
			return cfg.AllowLiveWebLookupInDMs
		}
		return cfg.AllowLiveWebLookupInGuilds
	}
	return false
}

func collectToolSources(toolContexts []chat.ToolContext) []string {
	sources := make([]string, 0, len(toolContexts))
	for _, toolContext := range toolContexts {
		result := strings.TrimSpace(toolContext.Result)
		if result == "" {
			continue
		}
		firstLine := result
		if idx := strings.Index(firstLine, "\n"); idx >= 0 {
			firstLine = firstLine[:idx]
		}
		sources = append(sources, firstLine)
	}
	return sources
}

func (e Engine) attachSources(reply string, sources []string, limit int) string {
	sources = sourceFooters(sources)
	if len(sources) == 0 {
		return reply
	}
	footerLines := []string{"", "Sources:"}
	for index, source := range sources {
		if index >= 4 {
			break
		}
		footerLines = append(footerLines, "- "+source)
	}
	footer := strings.Join(footerLines, "\n")
	combined := strings.TrimSpace(reply) + footer
	if limit > 0 && len([]rune(combined)) > limit {
		return strings.TrimSpace(reply)
	}
	return combined
}

// loadMemoryContext infers a subject from the input and returns relevant
// memory pack bodies. It does a best-effort subject match and silently ignores
// errors (no memory is preferable to breaking the reply path).
func (e Engine) loadMemoryContext(input string) []string {
	if e.rt == nil {
		return nil
	}
	subject := inferSubject(input)
	if subject == "" {
		return nil
	}
	packs, err := e.rt.MemoryBySubject(subject)
	if err != nil || len(packs) == 0 {
		return nil
	}
	const maxPacks = 3
	result := make([]string, 0, maxPacks)
	for _, p := range packs {
		if len(result) >= maxPacks {
			break
		}
		if strings.TrimSpace(p.Body) != "" {
			result = append(result, p.Body)
		}
	}
	return result
}

// inferSubject does a simple keyword scan of the input to pick a memory subject.
func inferSubject(input string) string {
	lower := strings.ToLower(input)
	subjects := []struct {
		keywords []string
		subject  string
	}{
		{[]string{"repo", "repository", "codebase", "github", "git"}, "repo"},
		{[]string{"discord", "bot", "transport"}, "discord"},
		{[]string{"model", "ollama", "llm", "ai"}, "models"},
		{[]string{"memory", "knowledge", "pack"}, "memory"},
		{[]string{"mcp", "tool", "playwright", "peekaboo"}, "mcp"},
	}
	for _, s := range subjects {
		for _, kw := range s.keywords {
			if strings.Contains(lower, kw) {
				return s.subject
			}
		}
	}
	return ""
}

func (e Engine) replyLimit(env gateway.Envelope) int {
	if e.rt == nil {
		return 0
	}
	if env.Surface == gateway.SurfaceDiscord {
		return e.rt.Config().Config.Transports.Discord.MaxResponseChars
	}
	return 900
}

func inferTask(input string, env gateway.Envelope) models.Task {
	cleaned := strings.ToLower(strings.TrimSpace(input))
	words := strings.Fields(cleaned)
	task := models.Task{Complexity: models.ComplexityRoutine}

	if env.Metadata["dm"] == "true" {
		task.PrivateContext = true
	}
	if len(env.Attachments) > 0 {
		task.Complexity = models.ComplexityStandard
	}
	if len(words) > 18 || strings.Contains(cleaned, "?") {
		task.Complexity = models.ComplexityStandard
	}
	for _, hint := range []string{"architecture", "design", "plan", "strategy", "review", "debug", "investigate", "analyze", "root cause", "compare", "refactor"} {
		if strings.Contains(cleaned, hint) {
			task.Complexity = models.ComplexityComplex
			break
		}
	}

	return task
}

func conversationKey(env gateway.Envelope) string {
	if env.Surface == gateway.SurfaceDiscord && env.Metadata["dm"] != "true" && env.ConversationID != "" && env.SenderID != "" {
		return string(env.Surface) + ":" + env.ConversationID + ":" + env.SenderID
	}
	if env.ConversationID != "" {
		return string(env.Surface) + ":" + env.ConversationID
	}
	if env.SenderID != "" {
		return string(env.Surface) + ":" + env.SenderID
	}
	return string(env.Surface) + ":default"
}

type conversationHistory struct {
	mu          sync.RWMutex
	maxMessages int
	turns       map[string][]chat.Turn
}

func newConversationHistory(maxMessages int) *conversationHistory {
	return &conversationHistory{
		maxMessages: maxMessages,
		turns:       map[string][]chat.Turn{},
	}
}

func (h *conversationHistory) load(key string) []chat.Turn {
	if h == nil || key == "" {
		return nil
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	stored := h.turns[key]
	cloned := make([]chat.Turn, len(stored))
	copy(cloned, stored)
	return cloned
}

func (h *conversationHistory) append(key string, turns ...chat.Turn) {
	if h == nil || key == "" || len(turns) == 0 {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.turns[key] = append(h.turns[key], turns...)
	if h.maxMessages > 0 && len(h.turns[key]) > h.maxMessages {
		h.turns[key] = append([]chat.Turn(nil), h.turns[key][len(h.turns[key])-h.maxMessages:]...)
	}
}
