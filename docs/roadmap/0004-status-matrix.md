# Product Status Matrix

This matrix reflects actual shipped capability versus roadmap intent as of the current binary and repo state. It is meant to stay honest and be updated as features ship.

Status key:
- **Shipped** — working in the current binary, validated locally
- **Partial** — code exists, boundaries defined, but not fully wired or production-ready
- **Not started** — in contracts or roadmap only

---

## CLI and operator surfaces

| Feature | Status | Notes |
| --- | --- | --- |
| `status`, `doctor`, `agents`, `skills` commands | Shipped | All return real workspace data |
| `integrations` command | Shipped | Reports Ollama, GitHub, HF, Kaggle readiness |
| `gateway`, `persona`, `session` commands | Shipped | Report live config state |
| `discord status`, `discord preview`, `discord run` | Shipped | Full Discord CLI surface |
| `models` command with four-lane profiles | Shipped | local-default, copilot-0x, copilot-0.33x, copilot-1x |
| `contracts` command | Shipped | Lists all named contracts with stage and goal |
| `plan` command | Shipped | Reports phase plan from roadmap |
| JSON output mode | Not started | No `--json` flag yet |
| Session resume across restarts | Not started | Session store exists but is not persisted to disk |
| Structured logging | Not started | No structured log output from the runtime |

---

## Model routing

| Feature | Status | Notes |
| --- | --- | --- |
| Four-lane profile definitions | Shipped | Declared in config and reported by CLI |
| Policy-driven routing by lane | Partial | Lane selection is defined; adaptive routing based on task type is not implemented |
| Fallback and error routing | Partial | Fail-closed on cloud if not opted in; no retry logic |
| Model switching from CLI | Not started | No `savitar models use <profile>` command |

---

## Reply engine and conversation

| Feature | Status | Notes |
| --- | --- | --- |
| Ollama-backed conversational replies | Shipped | Working through the Discord path |
| Per-user in-memory chat history | Shipped | Short rolling window, Discord only |
| Operator command bypass (deterministic responses) | Shipped | help, ping, status, etc. skip the model |
| Persona and voice policy | Shipped | Loaded from config, applied in replies |
| Reply engine for non-Discord surfaces | Not started | Engine exists but is only wired to Discord |
| Context-aware reply (MCP tools used inline) | Not started | No tool use in the reply path yet |

---

## Provider integrations

| Feature | Status | Notes |
| --- | --- | --- |
| Ollama local detection and credential check | Shipped | Validated through `integrations` command |
| GitHub credential check | Shipped | Validated through `integrations` command |
| Hugging Face credential check | Shipped | Validated through `integrations` command |
| Kaggle credential check | Shipped | Validated through `integrations` command |
| Keychain-backed secret storage | Shipped | Used by `store-local-integrations.sh` |
| Ollama API calls for inference | Shipped | Used by the Discord reply path |
| GitHub API calls (code, PRs, issues) | Not started | Credential present, no runtime client |
| Hugging Face API calls | Not started | Credential present, no runtime client |

---

## MCP orchestration

| Feature | Status | Notes |
| --- | --- | --- |
| MCP server config in `.vscode/mcp.json` | Shipped | Config exists and is inspectable |
| `integrations` reporting of MCP server state | Shipped | Four servers listed |
| Runtime MCP client that calls tools | Not started | No Go MCP client wired to the runtime |
| GitHub MCP tool use | Not started | |
| Context7 docs tool use | Not started | |
| Playwright browser automation | Not started | |
| Peekaboo macOS desktop inspection | Not started | |

---

## Skills and hooks

| Feature | Status | Notes |
| --- | --- | --- |
| Workspace skill discovery from `.github/skills` | Shipped | Reported by `skills` command |
| Workspace agent discovery from `.github/agents` | Shipped | Reported by `agents` command |
| Skill execution wired to the agent loop | Not started | Discovery works; execution is not wired |
| Hook runner for pre/post tool actions | Not started | Defined in contracts, not implemented |
| Auditable hook scripts | Not started | |

---

## Discord transport

| Feature | Status | Notes |
| --- | --- | --- |
| Discord bot token auth | Shipped | |
| Mention-or-DM inbound filter | Shipped | |
| Fail-closed DM cloud policy | Shipped | |
| Guild cloud opt-in | Shipped | Config flag, validated locally |
| Per-user cooldown (5s default) | Shipped | |
| Concurrency cap (2 concurrent replies) | Shipped | |
| Operator allowlist enforcement | Shipped | Works; currently zero operators configured locally |
| Conversational reply via Ollama | Shipped | |
| Rolling per-user-per-channel chat history | Shipped | In-memory only |
| Durable Discord conversation history | Not started | |
| Channel-specific reply policy | Partial | Opt-in flag exists; no per-channel granularity yet |

---

## Messaging transports (non-Discord)

| Feature | Status | Notes |
| --- | --- | --- |
| WhatsApp bridge | Not started | Contract defined (phase-3) |
| iMessage bridge | Not started | Contract defined (phase-3) |
| Shared gateway envelope across transports | Partial | Gateway struct exists; only Discord uses it |
| Inbound message normalization | Partial | Exists for Discord; not abstracted across transports |

---

## Web UI

| Feature | Status | Notes |
| --- | --- | --- |
| Public web UI backend | Not started | Contract defined (phase-3) |
| Google OAuth session validation | Not started | |
| Operator dashboard | Not started | |

---

## Memory packs

| Feature | Status | Notes |
| --- | --- | --- |
| Durable memory pack files | Not started | Contract defined (phase-4) |
| Subject-area indexed summaries | Not started | |
| Retrieval layer | Not started | |
| Cross-session context | Not started | |
| Session store (in-process) | Partial | Store struct exists; not persisted to disk |

---

## Automation and remote control

| Feature | Status | Notes |
| --- | --- | --- |
| Browser automation via Playwright | Not started | Contract defined (phase-2) |
| macOS desktop control via Peekaboo | Not started | Contract defined (phase-2) |
| Screen capture and recording | Not started | Contract defined (phase-2) |
| Controlled shell execution policy | Not started | Contract defined (phase-2) |
| Repository analysis pipeline | Not started | Contract defined (phase-2) |
| Remote host control | Not started | Contract defined (phase-3) |

---

## Hardening and production readiness

| Feature | Status | Notes |
| --- | --- | --- |
| Long-running service mode | Not started | |
| launchd packaging | Not started | |
| Audit logs | Not started | |
| Retry policy | Not started | |
| Structured logs | Not started | |
| Operator controls and access review | Partial | Allowlist logic exists; no admin surface to manage it |
| GitHub Actions artifact builds | Not started | |
| Release process | Not started | |
