---
name: savitar-transport-design
description: 'Design or extend Savitar messaging transports and the shared gateway. Use when working on Discord, WhatsApp, iMessage, message normalization, delivery flows, or channel identity policy.'
argument-hint: 'Describe the transport or gateway change'
---

# Savitar Transport Design Workflow

## When to Use

- Adding a new messaging platform.
- Changing inbound or outbound message handling.
- Designing the shared gateway process.

## Procedure

1. Define the normalized message envelope before adapter logic.
2. Split transport auth, delivery, identity, and message transformation concerns.
3. Keep account identifiers and secrets in local config only.
4. Update contracts, config, setup docs, and operator runbooks together.
5. Add transport-specific validation and a security review.