# Build Order: Next 4 Milestones

This document defines the concrete implementation order for the next phase of Savitar development. Each milestone is a complete, shippable increment. Milestones are ordered by dependency: each one unlocks the next.

---

## Milestone 1: MCP runtime client

**What it is:** Wire a Go MCP client to the runtime so the reply engine can call external tools during a conversation.

**Why it comes first:** Every subsequent milestone (browser automation, repo analysis, memory) depends on the runtime being able to invoke tools. Right now the workspace config lists four MCP servers but the runtime never calls any of them. Nothing else meaningful can be added to the agent loop until this gap is closed.

**What to build:**

1. Add a `internal/mcp/` package with a minimal Go MCP client that can connect to a local MCP server over stdio or HTTP.
2. Wire it to `internal/runtime/runtime.go` so the runtime can enumerate available tools from each configured server at startup.
3. Add a `savitar mcp status` command that reports which servers are reachable and which tools they expose.
4. Extend the reply engine in `internal/respond/engine.go` to accept a tool-call result from a running MCP server and incorporate it into the response.
5. Validate with Context7 first (documentation lookup), then GitHub (repo inspection).

**Acceptance:** `savitar mcp status` lists tools from at least two servers. A Discord reply that asks about a codebase triggers a GitHub MCP lookup and incorporates the result.

**Files touched:** `internal/mcp/` (new), `internal/runtime/runtime.go`, `internal/respond/engine.go`, `cmd/savitar/main.go`.

---

## Milestone 2: Controlled shell and repository analysis

**What it is:** Add a policy-gated shell execution layer and a repository analysis pipeline so Savitar can clone, inspect, and summarize code under explicit operator control.

**Why it comes second:** This is the primary Phase 2 gap between "knows about tools" and "can safely use tools." It also enables the memory pack work in Milestone 4 because you need a way to ingest external repositories before you can summarize them.

**What to build:**

1. Add a `internal/shell/` package with a policy struct that enforces: explicit command allowlist, working-directory confinement, max-runtime timeout, and stdout/stderr capture to a local log.
2. Add a `internal/repo/` package with a pipeline that clones a target repository to a local workspace directory, runs `git log`, `git diff`, and basic file enumeration, and returns a structured summary.
3. Wire `internal/shell/` policy to the MCP tool-call path from Milestone 1 so that any shell-backed tool call passes through the policy check.
4. Add a `savitar repo analyze <url>` command that runs the pipeline and writes a structured summary to `workspace/repo-analysis/`.
5. Add provenance metadata (source URL, commit hash, timestamp, operator who triggered it) to every analysis output.

**Acceptance:** `savitar repo analyze <url>` produces a structured summary with provenance attached. A shell command that is not on the allowlist is rejected and logged without execution.

**Files touched:** `internal/shell/` (new), `internal/repo/` (new), `internal/mcp/` (policy integration), `cmd/savitar/main.go`.

---

## Milestone 3: Durable session and structured logging

**What it is:** Persist session state to disk and add structured log output so conversation context survives restarts and the system is observable without a debugger.

**Why it comes third:** The Discord reply surface works, but every restart loses all context. Structured logs are needed before adding more tool surfaces because the failure modes multiply quickly and you need traceable run output to diagnose them. This milestone unblocks memory packs and production hardening.

**What to build:**

1. Extend `internal/session/store.go` to persist session state to `~/.savitar/sessions/<id>.json` (or the configured data directory) on every write. Load from disk on startup if a session file exists.
2. Add a session GC policy: sessions older than a configurable max age (default 7 days) are pruned at startup.
3. Add a `internal/log/` package that wraps the standard library logger with structured JSON output. Fields: timestamp, level, component, surface, session ID, and message.
4. Replace all `fmt.Println` and `log.Printf` calls in `internal/runtime/`, `internal/discord/`, and `internal/respond/` with structured log calls.
5. Add a `savitar session list` command that shows all on-disk sessions with their last-active timestamp and surface.

**Acceptance:** Restart the Discord bot, send a message, restart again. The bot correctly references the prior conversation. `savitar session list` shows the persisted sessions. Log output is valid JSON.

**Files touched:** `internal/session/store.go`, `internal/log/` (new), `internal/runtime/runtime.go`, `internal/discord/bot.go`, `internal/respond/engine.go`, `cmd/savitar/main.go`.

---

## Milestone 4: Memory pack foundation

**What it is:** Add a filesystem-backed memory pack layer that captures subject-area summaries, attaches provenance, and makes them retrievable by the reply engine.

**Why it comes fourth:** Milestones 1–3 together provide: tool calls that produce structured output (M1), a repository ingestion pipeline that produces summaries with provenance (M2), and a durable session layer that can reference prior context (M3). Memory packs are the synthesis layer on top of all three.

**What to build:**

1. Add a `internal/memory/` package with a `Pack` type: a named, versioned, filesystem-backed document with provenance metadata, a subject tag, and free-form content.
2. Add a writer that serializes a Pack to `~/.savitar/memory/<subject>/<name>.md` using the same frontmatter format used elsewhere in the project.
3. Add a reader that loads all Packs matching a subject tag and returns them as context for the reply engine.
4. Wire the repository analysis pipeline from Milestone 2 to automatically write a memory pack when an analysis run completes.
5. Extend `internal/respond/engine.go` to load relevant memory packs by subject before generating a reply, so prior analyses inform current responses.
6. Add `savitar memory list` and `savitar memory show <name>` commands.

**Acceptance:** Run `savitar repo analyze <url>`. A memory pack appears under `~/.savitar/memory/`. Ask a Discord question about that repository. The reply cites the prior analysis without re-fetching.

**Files touched:** `internal/memory/` (new), `internal/repo/` (memory write call), `internal/respond/engine.go`, `cmd/savitar/main.go`.

---

## After milestone 4

The next natural boundary after these four milestones is the gateway expansion: adding a second real transport (WhatsApp or iMessage) that shares the same gateway envelope, session layer, memory packs, and MCP tool path that the Discord surface already uses. That work is well-defined in `docs/roadmap/0001-phase-plan.md` Phase 3 and in `internal/contracts/contracts.go`.
