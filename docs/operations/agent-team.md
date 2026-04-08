# Agent Team

Savitar should be built by a small set of specialized agents rather than one overloaded generalist. The main coding agent remains the coordinator, and these custom agents act as focused subagents.

## Team roles

| Agent | Responsibility |
| --- | --- |
| Savitar Architect | ADRs, contracts, sequencing, boundaries, and tradeoff analysis |
| Savitar Builder | Implementing Go runtime, CLI, config, MCP, and core backend work |
| Savitar Transport Lead | Discord, WhatsApp, iMessage, gateway contracts, and public channel behavior |
| Savitar Web UI Lead | Public web UI, Google login, operator console, and session UX |
| Savitar QA | Tests, doctor checks, CI validation, regression coverage, and release confidence |
| Savitar Security | Threat modeling, auth review, tool/MCP risk review, secret handling, and public-surface security findings |

## Recommended workflow

1. Start with Savitar Architect for any non-trivial feature.
2. Hand implementation to Savitar Builder, Savitar Transport Lead, or Savitar Web UI Lead depending on the surface.
3. Run Savitar QA before considering the work complete.
4. Run Savitar Security for any public-facing, auth-related, or tool-execution change.
5. Update provenance whenever external docs or parity references inform the work.

## Helpful local commands

- `savitar agents` to inspect the current workspace agent team.
- `savitar skills` to inspect the current workspace skill catalog.
- `savitar status` to see whether the workspace and session surface are initialized.

## Internal planning board

- The generated internal planning snapshot lives at `docs/operations/internal-objectives-board.md`.
- Update `docs/operations/internal-objectives-source.md` when the strategic plan changes.
- When local git hooks are enabled with `core.hooksPath=scripts/hooks`, commits refresh and stage the generated board automatically.

## Non-negotiables

- Never hardcode phone numbers, account emails, OAuth credentials, or Apple IDs in version control.
- Keep the custom agents focused; do not turn them into generic copies of the main agent.
- Treat WhatsApp, Discord, iMessage, and public web access as separate adapters over a shared message core.
- Keep Savitar personable, but preserve honest disclosure when identity or safety matters.