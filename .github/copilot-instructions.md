# Project Guidelines

## Code Style

- Keep the core Go runtime standard-library-first unless a dependency clearly pays for itself.
- Prefer small packages, explicit names, and direct data flow over abstraction-heavy designs.
- Add brief comments only where behavior would otherwise be non-obvious.

## Architecture

- Treat `docs/adr/0001-clean-room-charter.md` as binding.
- Keep model routing, MCP integration, transports, memory, and automation as separate boundaries.
- Any new external integration must update the config example, the relevant setup docs, the agent-team docs when relevant, and the provenance log.

## Build and Test

- Bootstrap a Mac with `./scripts/bootstrap-macos.sh`.
- Verify local prerequisites with `./scripts/verify-prereqs.sh`.
- Use `make doctor`, `make test`, and `make build` once Go is available.

## Conventions

- Do not hardcode secrets, tokens, phone numbers, or Apple IDs.
- Prefer workspace MCP config in `.vscode/mcp.json` over undocumented local setup.
- Describe parity features in clean-room behavior terms instead of copying implementations from external projects.
- For public-facing persona work, prefer natural conversational style but preserve honest disclosure when identity, consent, access, billing, or safety are relevant.