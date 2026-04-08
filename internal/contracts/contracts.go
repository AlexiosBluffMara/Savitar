package contracts

type Contract struct {
	Name         string
	Stage        string
	Goal         string
	Dependencies []string
	Guardrails   []string
}

func Default() []Contract {
	return []Contract{
		{
			Name:         "orchestrator-program",
			Stage:        "rewrite-pass-1",
			Goal:         "Turn one inbound envelope into an evidence-backed reply, route decision, run record, and delivery plan.",
			Dependencies: []string{"gateway", "memory-retrieval", "model-routing", "tool-broker", "run-recorder"},
			Guardrails:   []string{"clean-room", "honest-disclosure", "traceable-run-record"},
		},
		{
			Name:         "gateway-core",
			Stage:        "rewrite-pass-1",
			Goal:         "Normalize surface traffic, hold trust decisions, and keep transport delivery separate from reply composition.",
			Dependencies: []string{"discord-adapter", "surface-policy"},
			Guardrails:   []string{"surface-normalization", "explicit-trust-policy", "no-hidden-transport-logic"},
		},
		{
			Name:         "knowledge-retrieval",
			Stage:        "rewrite-pass-2",
			Goal:         "Store packs, index sources, retrieve excerpts, and return citations that can survive operator review.",
			Dependencies: []string{"pack-store", "source-catalog", "retrieval-index"},
			Guardrails:   []string{"traceable-source-material", "subject-bounded-context", "privacy-review"},
		},
		{
			Name:         "evidence-tools",
			Stage:        "rewrite-pass-2",
			Goal:         "Fetch live evidence through approved MCP servers when memory packs are insufficient or stale.",
			Dependencies: []string{"mcp-registry", "github-mcp", "context7", "playwright"},
			Guardrails:   []string{"smallest-scope", "observable-tool-use", "reference-required"},
		},
		{
			Name:         "discord-surface",
			Stage:        "rewrite-pass-3",
			Goal:         "Handle end-user messaging on Discord through the orchestrator program and a reviewable delivery path.",
			Dependencies: []string{"discord-app", "discordgo", "gateway-core", "orchestrator-program"},
			Guardrails:   []string{"mention-or-dm-default", "explicit-channel-policy", "no-fake-capabilities"},
		},
		{
			Name:         "session-and-runs",
			Stage:        "rewrite-pass-3",
			Goal:         "Persist session state, route choice, citations, and delivery outcomes for every handled turn.",
			Dependencies: []string{"session-store", "run-store", "structured-logging"},
			Guardrails:   []string{"durable-audit-trail", "no-secret-logging", "reviewable-history"},
		},
		{
			Name:         "operator-review",
			Stage:        "rewrite-pass-4",
			Goal:         "Expose a review queue, run history, and operator decisions over the same core runtime state.",
			Dependencies: []string{"run-store", "knowledge-retrieval", "approval-policy"},
			Guardrails:   []string{"operator-auth", "reason-required-for-overrides", "traceable-decisions"},
		},
		{
			Name:         "operator-webui",
			Stage:        "rewrite-pass-5",
			Goal:         "Expose the operator review surface through authenticated HTTP routes and a minimal console UI.",
			Dependencies: []string{"operator-review", "google-oauth", "session-validation", "http-server"},
			Guardrails:   []string{"authenticated-access-only", "public-internet-threat-model", "audit-log"},
		},
		{
			Name:         "model-routing",
			Stage:        "rewrite-pass-2",
			Goal:         "Keep local Gemma 4 routing primary while making route choice visible in runs, replies, and operator review.",
			Dependencies: []string{"ollama", "router-policy", "persona-policy"},
			Guardrails:   []string{"local-first", "explicit-cloud-fallback", "visible-route-reason"},
		},
		{
			Name:         "expansion-transports",
			Stage:        "later",
			Goal:         "Add WhatsApp and iMessage only after the Discord plus operator-console workflow is coherent and reviewable.",
			Dependencies: []string{"gateway-core", "orchestrator-program", "operator-webui"},
			Guardrails:   []string{"truthful-roadmap-language", "separate-adapter-boundaries", "policy-review"},
		},
		{
			Name:         "expansion-automation",
			Stage:        "later",
			Goal:         "Add browser automation, screen capture, and remote execution only after the evidence and operator story is stable.",
			Dependencies: []string{"evidence-tools", "operator-review", "shell-policy"},
			Guardrails:   []string{"visible-automation", "recording-consent", "explicit-host-allowlist"},
		},
	}
}
