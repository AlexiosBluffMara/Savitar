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
			Name:         "core-agent-loop",
			Stage:        "phase-1",
			Goal:         "Route tasks, generate replies, persist session state, and coordinate the runtime.",
			Dependencies: []string{"go", "config", "model-router", "reply-engine"},
			Guardrails:   []string{"clean-room", "structured-logs", "honest-disclosure"},
		},
		{
			Name:         "mcp-orchestration",
			Stage:        "phase-2",
			Goal:         "Broker access to GitHub, docs, browser, and macOS tooling through MCP.",
			Dependencies: []string{"github-mcp", "context7", "playwright", "peekaboo"},
			Guardrails:   []string{"smallest-scope", "observable-tool-use"},
		},
		{
			Name:         "provider-integrations",
			Stage:        "phase-2",
			Goal:         "Use local-only credentials to reach Ollama Cloud, GitHub, Hugging Face, Kaggle, and Discord without committing secrets.",
			Dependencies: []string{"macos-keychain", "ollama", "gh", "huggingface-hub", "kaggle-cli"},
			Guardrails:   []string{"keychain-first", "scoped-tokens", "local-config-only"},
		},
		{
			Name:         "conversation-replies",
			Stage:        "phase-2",
			Goal:         "Generate model-backed conversational replies through Ollama while preserving deterministic operator commands.",
			Dependencies: []string{"ollama-api", "reply-engine", "persona-config"},
			Guardrails:   []string{"honest-disclosure", "operator-command-bypass", "no-secret-logging"},
		},
		{
			Name:         "discord-transport",
			Stage:        "phase-3",
			Goal:         "Handle inbound and outbound Discord messages through a bot-token-backed adapter over the shared gateway envelope.",
			Dependencies: []string{"discord-app", "discord-bot-token", "discordgo", "discord-message-content-intent"},
			Guardrails:   []string{"mention-or-dm-default", "explicit-channel-policy", "no-secret-logging"},
		},
		{
			Name:         "whatsapp-bridge",
			Stage:        "phase-3",
			Goal:         "Bridge WhatsApp through a compliant API or an explicitly owned Android companion service.",
			Dependencies: []string{"pixel-fold", "adb", "bridge-service"},
			Guardrails:   []string{"no-phone-number-in-repo", "policy-review"},
		},
		{
			Name:         "imessage-bridge",
			Stage:        "phase-3",
			Goal:         "Bridge iMessage through macOS-local automation surfaces.",
			Dependencies: []string{"osascript", "messages", "accessibility-permission"},
			Guardrails:   []string{"macos-only", "local-host-only"},
		},
		{
			Name:         "agent-identity",
			Stage:        "phase-1",
			Goal:         "Maintain a consistent Savitar persona across CLI, messaging, and web surfaces.",
			Dependencies: []string{"persona-config", "reply-policy"},
			Guardrails:   []string{"honest-disclosure", "no-personal-identifiers-in-repo"},
		},
		{
			Name:         "browser-automation",
			Stage:        "phase-2",
			Goal:         "Use Playwright and Peekaboo for controlled web and macOS automation.",
			Dependencies: []string{"playwright", "peekaboo"},
			Guardrails:   []string{"persistent-session-policy", "visible-automation"},
		},
		{
			Name:         "screen-capture",
			Stage:        "phase-2",
			Goal:         "Capture screenshots and screen recordings for browser and desktop workflows.",
			Dependencies: []string{"peekaboo", "macos-screen-recording-permission"},
			Guardrails:   []string{"recording-consent", "local-storage-policy"},
		},
		{
			Name:         "memory-packs",
			Stage:        "phase-4",
			Goal:         "Create durable subject-area memory dumps, summaries, and retrieval indices.",
			Dependencies: []string{"snapshot-store", "indexer"},
			Guardrails:   []string{"traceable-source-material", "privacy-review"},
		},
		{
			Name:         "skills-and-hooks",
			Stage:        "phase-2",
			Goal:         "Load workspace skills and enforce deterministic safety hooks around agent actions.",
			Dependencies: []string{"skill-loader", "hook-runner"},
			Guardrails:   []string{"auditable-scripts", "bounded-hook-runtime"},
		},
		{
			Name:         "repository-analysis",
			Stage:        "phase-2",
			Goal:         "Clone, inspect, summarize, and compare external repositories under controlled execution.",
			Dependencies: []string{"git", "shell-policy"},
			Guardrails:   []string{"read-only-by-default", "provenance-required"},
		},
		{
			Name:         "remote-host-control",
			Stage:        "phase-3",
			Goal:         "Run approved commands and workflows on explicitly configured remote or local hosts.",
			Dependencies: []string{"shell-policy", "host-registry"},
			Guardrails:   []string{"explicit-host-allowlist", "observable-command-log"},
		},
		{
			Name:         "public-webui",
			Stage:        "phase-3",
			Goal:         "Expose a public Savitar web UI with Google login, session controls, and operator guardrails.",
			Dependencies: []string{"web-backend", "google-oauth", "session-store"},
			Guardrails:   []string{"public-internet-threat-model", "secret-store", "audit-log"},
		},
	}
}
