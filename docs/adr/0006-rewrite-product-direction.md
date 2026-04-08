# ADR 0006: Rewrite Product Direction

## Status

Accepted

## Context

Savitar has a strong clean-room foundation, a working Discord surface, and credible local-first Gemma 4 infrastructure. It also has a scope problem. The repository currently mixes one real product with several future surfaces, causing the public story to outrun the shipped system.

The rewrite should not discard the current working core. It should narrow the product definition so the next implementation passes build toward one coherent system instead of extending every planned surface at once.

## Decision

Savitar will be rewritten around one primary product:

**A local-first community knowledge operator for Illinois academic, civic, and small-business workflows.**

The rewrite is governed by the following rules.

1. The first complete system is **Discord plus an authenticated operator console**.
2. The system's primary job is to deliver **evidence-backed replies** using curated knowledge packs and approved live tool results.
3. Public messaging and operator workflows must share one core orchestration path so the web surface is not a separate product.
4. Knowledge, provenance, citations, and operator review are part of the core product, not later polish.
5. WhatsApp, iMessage, Pixel Fold mobile operator flows, LiteRT, remote host control, and broader automation remain valid expansion tracks, but they are not part of the first coherent product claim.
6. Public materials must describe only shipped surfaces in present tense. Planned surfaces must be labeled as planned.
7. Faculty collaboration language must stay precise. Savitar may describe conversations, feedback, domain review, or ongoing collaboration only to the extent those roles are true and approved for public mention.

## First complete product definition

The first complete product includes:

- Discord as the only end-user messaging surface.
- An authenticated operator web console for conversation review, knowledge inspection, run history, and approvals.
- Filesystem-backed knowledge packs with traceable source metadata.
- MCP-backed evidence retrieval where the runtime can actually reach and use approved tools.
- A run record that captures route choice, tool use, evidence references, and delivery outcome.

The first complete product does not require:

- WhatsApp.
- iMessage.
- A public anonymous demo mode.
- Pixel Fold-specific features.
- On-device Android inference.
- Broad remote execution.

## Consequences

- Rewrite work should prioritize orchestration, knowledge retrieval, citations, operator review, and one operator-facing web surface.
- Existing Discord functionality should be preserved as the initial production-grade transport.
- Roadmap language and contracts should favor the narrowed product over the previous everything-at-once framing.
- Hackathon and academic positioning should emphasize a truthful, faculty-informed, local-first system rather than a multi-surface platform that is still under construction.