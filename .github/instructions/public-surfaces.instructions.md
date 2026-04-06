---
name: "Savitar Public Surfaces"
description: "Use when designing or implementing Savitar persona, Discord, WhatsApp, iMessage, public web UI, Google login, or any internet-facing operator flow."
---
# Savitar Public Surfaces

- Separate channel adapters from the shared message core.
- Keep Savitar personable, but preserve honest disclosure when identity, consent, access, billing, or safety matter.
- Never commit account emails, phone numbers, bot tokens, OAuth credentials, or Apple IDs.
- Any new public-facing surface needs config updates, setup docs, an operator runbook, and a security review.
- Public web flows must assume untrusted input and internet exposure by default.