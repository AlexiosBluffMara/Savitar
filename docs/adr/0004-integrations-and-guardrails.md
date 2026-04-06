# ADR 0004: Integrations and Guardrails

## Status

Accepted

## Decision

Savitar will support integrations through well-scoped bridges and explicit guardrails.

## Transport boundaries

- Discord uses a bot-token-backed integration.
- WhatsApp is treated as a bridge problem, not a core runtime assumption. The production path should be either a compliant WhatsApp Business API setup or an explicitly owned Android companion bridge.
- iMessage remains macOS-local and should flow through Apple-supported automation surfaces such as AppleScript or Shortcuts.

## Tooling boundaries

- GitHub MCP is the default GitHub surface.
- Context7 is the default live documentation surface.
- Playwright handles controlled web sessions.
- Peekaboo handles macOS-native UI inspection and control where browser automation is insufficient.

## Safety rules

- Personal phone numbers and account identifiers stay out of version control.
- Tool access should default to the smallest useful scope.
- Browser automation and shell execution must remain observable and reviewable.
- New integrations require setup docs, local config entries, and provenance records.

## Consequences

This keeps Savitar lightweight while still leaving a path to broad capability coverage.