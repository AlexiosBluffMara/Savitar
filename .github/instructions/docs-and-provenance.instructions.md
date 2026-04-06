---
name: "Savitar Docs and Provenance"
description: "Use when updating Savitar ADRs, roadmap docs, setup guides, public-surface docs, or provenance records."
applyTo: ["docs/**/*.md", "README.md"]
---
# Savitar Docs and Provenance

- ADRs should record decisions and consequences, not vague intentions.
- Roadmap files should separate current work from later ambition.
- Setup docs should distinguish required now versus planned later.
- Never commit phone numbers, account emails, bot tokens, OAuth secrets, or Apple IDs into docs or examples.
- If a feature was informed by public docs or parity targets, add or update a provenance entry.