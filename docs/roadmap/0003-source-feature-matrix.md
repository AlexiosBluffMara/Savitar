# Source Feature Matrix

Savitar will combine capabilities from several prior projects, but it will not merge their codebases. This matrix is the intake map for behavior, operator experience, and product scope.

## Source repos

| Source repo | Primary value | Savitar interpretation |
| --- | --- | --- |
| `ultraworkers/claw-code` | Rust CLI harness with sessions, slash commands, permissions, hooks, MCP inspection, JSON/text output, and subagent surfaces | Use as the model for a disciplined operator CLI and clear runtime surfaces |
| `AlexiosBluffMara/claw-code` | Clean-room rewrite framing, manifest output, workspace inventory, and audit mindset | Use as the model for provenance-friendly reporting and parity tracking |
| `AlexiosBluffMara/picoclaw` | Lightweight Go runtime, broad channel support, skills, routing, gateway patterns, launcher surfaces, and deploy-anywhere posture | Use as the model for Go-first efficiency, gateway breadth, and low-overhead packaging |
| `AlexiosBluffMara/Artemis` | Deep memory/skills, messaging gateway, cron, migration, remote execution backends, and production ambition | Use as the model for a long-lived agent platform instead of a single CLI loop |

## Feature extraction by category

### CLI and operator surfaces

- From `ultraworkers/claw-code`: direct CLI subcommands, doctor/status/model/config surfaces, session resume, slash-command parity, JSON output.
- From `AlexiosBluffMara/picoclaw`: setup-oriented commands such as onboarding, gateway start, model selection, and skills install.
- From `AlexiosBluffMara/Artemis`: setup, update, gateway, migration, and messaging-aware command design.

Savitar plan:
- Build a compact Go CLI with `doctor`, `plan`, `contracts`, `models`, then expand toward `status`, `gateway`, `memory`, `skills`, and `session`.

### Model routing

- From `AlexiosBluffMara/picoclaw`: rule-based routing and broad provider support.
- From `AlexiosBluffMara/Artemis`: multi-provider flexibility and operator switching.
- From `ultraworkers/claw-code`: explicit model inspection and switching surfaces.

Savitar plan:
- Keep a four-lane router: local Gemma lane, Copilot `0x`, Copilot `0.33x`, and Copilot `1x`.
- Start with policy-driven routing before any adaptive or learned routing.

### MCP, hooks, and skills

- From `ultraworkers/claw-code`: MCP config inspection, hooks, permission modes, skill and agent discovery.
- From `AlexiosBluffMara/picoclaw`: MCP server configuration, skill installation, steering, SubTurn, and hook-driven execution.
- From `AlexiosBluffMara/Artemis`: large-scale skills and plugin ecosystem expectations.

Savitar plan:
- Keep workspace MCP config in `.vscode/mcp.json`.
- Build a first-party doctor and contract layer around configured MCP servers.
- Prefer local skill folders and auditable hooks before any remote marketplace.

### Messaging and gateways

- From `AlexiosBluffMara/picoclaw`: broad chat-channel mindset and gateway process.
- From `AlexiosBluffMara/Artemis`: Discord, WhatsApp, Slack, and other gateway-backed messaging surfaces.

Savitar plan:
- First-class transports: Discord, WhatsApp, and iMessage.
- Add a public web UI as a first-class surface rather than treating it as a separate product.
- Defer other channels until the gateway abstraction is stable.
- Treat WhatsApp as a bridge problem with compliance and device ownership constraints.

### Automation and remote control

- From `AlexiosBluffMara/Artemis`: remote execution backends and cloud-oriented runtime thinking.
- From `AlexiosBluffMara/picoclaw`: Android and lightweight deployment posture.
- From user requirements: Peekaboo, Playwright, controlled web sessions, screenshots, screen recording, and remote computer usage.

Savitar plan:
- Browser automation via Playwright MCP.
- macOS desktop inspection and control via Peekaboo.
- Screen capture and recordings as explicit capabilities with local storage policy.
- Remote host execution only through explicit host allowlists and observable logs.
- Keep the eventual Pixel Fold native companion as a later surface, not a phase-0 dependency.

### Persona and public identity

- From `AlexiosBluffMara/Artemis`: a strong, persistent agent identity across channels.
- From user requirements: Savitar should feel like a real operator with consistent voice and judgment.

Savitar plan:
- Give Savitar a consistent human conversational style and a distinct point of view.
- Keep the system honest about being software whenever identity, consent, access, or trust are materially affected.

### Memory and subject-matter dumps

- From `AlexiosBluffMara/Artemis`: persistent memory, user modeling, searchable history, and skills learned from work.
- From `AlexiosBluffMara/claw-code`: workspace manifesting and audit summaries.
- From `AlexiosBluffMara/picoclaw`: JSONL memory and scheduled task reminders.

Savitar plan:
- Start with filesystem-backed memory packs and indexed summaries.
- Keep provenance attached to imported knowledge and repository analyses.
- Separate durable memory from transient session logs.

## What Savitar should not inherit

- Any source-code implementation copied from the source repos.
- Architecture sprawl without clear runtime boundaries.
- Marketplace-first complexity before local workflows are solid.
- Heavy dependencies that do not justify their cost on a Mac Mini host.

## Build order implied by this matrix

1. CLI surfaces and routing.
2. MCP manager, hooks, skills, and browser/desktop automation.
3. Messaging gateway and remote-host control.
4. Durable memory packs and repository intelligence.
5. Setup wizard, migration helpers, and production hardening.