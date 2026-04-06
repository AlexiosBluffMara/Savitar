# Provenance Entry

- Workstream: Discord transport implementation for Savitar
- Date: 2026-04-06
- Author: GitHub Copilot
- External systems or docs consulted: Discord Gateway documentation at `docs.discord.com/developers/events/gateway#gateway-intents` and the `github.com/bwmarrin/discordgo` package documentation at `pkg.go.dev`
- Public interfaces or behaviors captured: bot-token-backed session startup, gateway intents, privileged `MESSAGE_CONTENT` behavior, Discord message send and edit flow, immediate placeholder reply pattern, mention-or-DM-safe default behavior, and operator allowlisting for diagnostic commands
- Explicit non-goals or material not copied: No Discord client source files, sample bot implementations, or third-party prompt logic were copied into Savitar
- New files or modules created from scratch: `internal/discord`, `internal/respond`, `docs/operations/discord-bot.md`, CLI Discord surfaces, config example expansion
- Verification completed: `go test ./...`, `go build ./...`, and local preview-path validation through the new Discord reply engine
- Follow-up provenance needed: slash commands, approval-gated tool execution from Discord, and any future Discord moderation or multi-guild operator workflows