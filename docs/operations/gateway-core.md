# Gateway Core

Savitar now has a shared gateway core in `internal/gateway`.

This package is the boundary that future Discord, WhatsApp, iMessage, and web UI transports should share.

## Current responsibilities

- Define the normalized surface list for local status and future transport startup.
- Hold a clean-room `Envelope` model for inbound and outbound messages.
- Keep identity, delivery mode, and transport details separate from app-level persona logic.
- Feed the live Discord adapter through the same `Envelope` model used by future transports.

## Immediate follow-on work

1. Add message normalization from WhatsApp, iMessage, and web UI into `gateway.Envelope`.
2. Add delivery planning and approval hooks before outbound sends.
3. Add per-surface auth and trust policy checks before public or remote actions.