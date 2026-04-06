# ADR 0002: Core Architecture

## Status

Accepted

## Decision

Savitar will be organized around explicit runtime boundaries so the system stays understandable as it gains tools and transports.

## Primary boundaries

- `app`: CLI entry point and human-facing commands.
- `config`: local configuration loading and defaults.
- `models`: routing policy for local and Copilot-backed model lanes.
- `runtime`: environment checks, capability summaries, and startup coordination.
- `contracts`: stable descriptions of transports, browser automation, memory, and repository analysis capabilities.

## External capability boundaries

- MCP servers provide tool access for GitHub, documentation lookup, browser automation, and macOS automation.
- Messaging transports are separate adapters for Discord, WhatsApp, and iMessage.
- Memory capture is treated as a product surface, not an incidental log folder.
- Repository analysis and controlled shell execution are isolated behind explicit contracts.

## Non-goals for phase 0

- No long-running daemon.
- No background job scheduler.
- No external database dependency.
- No coupled transport implementation hidden inside model logic.

## Rationale

This keeps the starter repo readable for non-specialists while leaving room for production hardening later.