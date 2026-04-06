# Savitar

Savitar is a clean-room autonomous agent starter focused on four things from day one: local-first model routing, MCP-driven tool access, messaging transports, and durable knowledge capture. This repository is the phase-0 foundation for that system: operating rules, environment setup, workspace MCP configuration, deployment notes, and a small Go CLI that expresses the intended boundaries without external runtime dependencies.

The current Discord surface can now answer normal messages through an Ollama-backed reply path while keeping operator diagnostics on explicit commands and failing closed unless local Ollama is available or cloud replies are explicitly allowed for that channel type.

## Current state

Savitar is a well-structured starter runtime with one real working transport. It is not yet a full autonomous operator platform.

What works right now:

- A Go CLI with `status`, `doctor`, `agents`, `skills`, `integrations`, `gateway`, `persona`, `session`, `models`, `contracts`, `plan`, and `discord` subcommands — all returning real workspace data.
- Local config loading, four-lane model profile definitions, and runtime state reporting.
- Workspace agent and skill discovery from `.github/agents` and `.github/skills`.
- Provider integration checks for Ollama, GitHub, Hugging Face, and Kaggle via keychain-backed credentials.
- A working Discord bot surface: mention-or-DM filtering, Ollama-backed conversational replies, per-user cooldown, concurrency caps, operator allowlist enforcement, and rolling in-memory chat history.

What does not work yet:

- No runtime MCP client. The workspace config lists GitHub, Context7, Playwright, and Peekaboo, but the runtime does not call any of them.
- No WhatsApp or iMessage transport beyond contract and config shape.
- No public web UI, Google login, or operator dashboard.
- No durable memory, indexed retrieval, or cross-session context.
- No structured logging, audit logs, retry policy, or long-running service mode.

See `docs/roadmap/0004-status-matrix.md` for the full shipped-versus-planned breakdown.

## Next

The immediate build priority is making the tool surfaces real:

1. Wire a Go MCP client to the runtime so the reply engine can call GitHub, Context7, Playwright, and Peekaboo.
2. Add controlled shell execution policy and repository analysis to the agent loop.
3. Persist session state to disk so conversation context survives process restarts.
4. Add structured logging and traceable run output so the system is observable without a debugger.

See `docs/roadmap/0005-build-order.md` for milestone-level detail.

## Later

Multi-surface operator platform:

- WhatsApp and iMessage transports behind the shared gateway model.
- Public web UI with Google login and operator workflows.
- Durable memory packs, indexed summaries, and long-lived cross-session context.
- Screen capture, browser automation, and macOS desktop control as first-class surfaces.
- Production hardening: supervised service mode, audit logs, launchd packaging, and release process.
- Eventual Pixel Fold native companion for Android-device actions.

See `docs/roadmap/0001-phase-plan.md` and `docs/roadmap/0002-parity-targets.md` for the full scope.

## Current focus

- Set up a repeatable Mac Mini work environment.
- Lock in clean-room and provenance discipline before feature work.
- Define the first runtime boundaries for models, MCP, transports, browser automation, and memory packs.
- Keep the core runtime readable, standard-library-first, and easy to audit.

## Quick start

1. Read `docs/setup/mac-mini.md` and `docs/setup/mcp-stack.md`.
2. Copy `config/savitar.example.json` to `config/savitar.local.json` and keep secrets local.
3. Run `./scripts/bootstrap-macos.sh`.
4. Run `./scripts/verify-prereqs.sh`.
5. Run `make doctor`, `make test`, and `make build` once Go is installed.
6. Review the current checkout, then run `export SAVITAR_REVIEWED_CHECKOUT=1` in the shell you trust.
7. Run `./scripts/store-local-integrations.sh`, then source `./scripts/use-local-integrations.sh`.
8. Run `./scripts/validate-local-integrations.sh`.
9. Build a reviewed Savitar binary, then run `./bin/savitar integrations`, `./bin/savitar discord status`, and `./bin/savitar discord preview "What can you help with?"` before bringing up live providers.

## Current CLI surface

- `savitar status` shows workspace, config, session, and surface counts.
- `savitar doctor` checks local prerequisites.
- `savitar agents` lists the custom Savitar agent team from `.github/agents`.
- `savitar skills` lists workspace skills from `.github/skills`.
- `savitar integrations` shows local provider integration readiness for Ollama Cloud, GitHub, Hugging Face, and Kaggle.
- `savitar gateway` summarizes Discord, WhatsApp, iMessage, and web UI surface status.
- `savitar persona` prints the current Savitar voice policy.
- `savitar session [show|init]` inspects or initializes local session state.
- `savitar discord [status|preview|run]` inspects, previews, or starts the Discord bot transport, including the model-backed reply path for normal messages.

## Project map

- `docs/adr/` contains the binding architectural decisions.
- `docs/roadmap/` contains the phase plan and parity targets.
- `docs/setup/` contains workstation and MCP setup guidance.
- `docs/setup/provider-integrations.md` covers local keychain-backed provider setup.
- `docs/operations/` contains the agent-team workflow and deployment runbooks.
- `docs/operations/discord-bot.md` is the operator runbook for the Discord surface.
- `docs/provenance/` contains clean-room evidence and templates.
- `cmd/savitar/` contains the CLI entry point.
- `internal/` contains the runtime packages.

## Phase order

1. Environment and policy baseline.
2. Core agent loop and model routing.
3. MCP orchestration and browser automation.
4. Messaging transports.
5. Memory packs and repository analysis.
6. Hardening, packaging, and deployment.

See `docs/roadmap/0001-phase-plan.md` for the execution order and `docs/roadmap/0002-parity-targets.md` for the capability targets.
See `docs/roadmap/0003-source-feature-matrix.md` for the clean-room feature intake from the source repositories.
See `docs/operations/agent-team.md` for the specialized agent team that should be used to build Savitar out.