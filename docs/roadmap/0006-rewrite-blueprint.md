# Rewrite Blueprint

This document defines the new Savitar build direction after the scope reset in [ADR 0006](../adr/0006-rewrite-product-direction.md).

It is intentionally concrete. The goal is to give later implementation passes stable package boundaries, method-level intent, and a validation sequence that turns the current starter runtime into one coherent product.

## Product Thesis

Savitar should become the local-first operator for trusted community knowledge.

In practical terms, that means:

1. A student, faculty member, civic partner, or small-business owner can ask a question on Discord.
2. Savitar answers using curated knowledge packs and, when needed, approved live tools.
3. An operator can inspect the evidence, the route, and the reply in an authenticated dashboard.
4. The system records enough provenance that a reviewer can tell what Savitar knew, what it fetched, and why it answered the way it did.

## Rewrite Scope

### In scope for the rewrite

- Discord as the first end-user surface.
- An authenticated operator console.
- Knowledge packs, retrieval, citations, and source metadata.
- MCP-backed live evidence for GitHub, docs, and later approved web lookups.
- Run records, session continuity, and operator review points.

### Explicitly deferred

- WhatsApp.
- iMessage.
- Public anonymous web access.
- Pixel Fold-specific UX.
- LiteRT and on-device Android inference.
- Broad remote host automation as a public-facing feature.

Those items remain on the roadmap, but they do not define the first complete product.

## Stakeholders And Their Success Criteria

| Stakeholder | What they need from the rewrite |
| --- | --- |
| End users on Discord | Direct answers, visible trust, fast enough latency, no fake capabilities |
| Operators | One dashboard for conversations, knowledge packs, review queue, and run history |
| Faculty collaborators | Precise sourcing, clean collaboration language, and a serious civic/educational use case |
| Hackathon judges | One working system with Gemma 4 clearly central to the experience |
| Maintainers | Stable package boundaries, small interfaces, traceable behavior |

## Core Runtime Shape

The rewrite should converge on five major backend boundaries.

| Boundary | Responsibility | First output |
| --- | --- | --- |
| `gateway` | Normalize inbound and outbound surface traffic and surface trust policy | Discord adapter plus delivery plan types |
| `orchestrator` | Turn one inbound envelope into a routed, evidence-backed reply | Request handling program and run record |
| `memory` | Store packs, index sources, retrieve contextual evidence, attach citations | Retrieval contracts and source catalog |
| `operator` | Build review queues, dashboard snapshots, approvals, and run views | Operator service contracts |
| `webui` | Expose the operator service through authenticated HTTP routes | Backend route map and handler stubs |

## Turn Lifecycle

Every user-visible reply should eventually pass through the same sequence.

1. Normalize inbound surface data into `gateway.Envelope`.
2. Evaluate trust, operator scope, and any review requirements.
3. Resolve the user intent and the task complexity.
4. Load session state and subject hints.
5. Retrieve memory context from the knowledge layer.
6. Decide whether live tool calls are required.
7. Execute approved tool calls and normalize the results.
8. Route the reply request to the correct model lane.
9. Compose the final reply with references.
10. Record the run, session update, and delivery result.
11. Expose the turn in the operator console.

## Data Records That Must Exist

The rewrite should standardize the following records before adding more transports.

### `SourceRecord`

Tracks external material that informed a knowledge pack or tool result.

Required fields:

- stable ID
- subject
- source kind
- title
- origin URL or local origin description
- last-checked timestamp
- provenance note

### `KnowledgePack`

The durable document the system can load as reusable context.

Required fields:

- name
- subject
- source reference
- created timestamp
- updated timestamp
- body

### `EvidenceRecord`

The narrow excerpt actually used in a reply.

Required fields:

- subject
- pack or tool source
- excerpt text
- locator
- confidence or score

### `RunRecord`

The canonical record of one handled request.

Required fields:

- surface
- conversation ID
- request time
- route selected
- tool calls attempted
- references attached
- review outcome
- delivery outcome

### `OperatorDecision`

Tracks a human approval, rejection, or escalation.

Required fields:

- operator identity
- decision type
- reason
- timestamp
- linked run ID

## Package-Level Build Notes

### `internal/gateway`

Keep the current envelope model, but add:

- inbound adapter contracts
- outbound delivery contracts
- approval and trust policy contracts
- reference metadata on outbound delivery

### `internal/orchestrator`

This is the new center of gravity.

It should own:

- intent resolution
- knowledge loading
- tool planning
- tool execution
- route selection
- reply composition
- run recording

It should not own surface-specific transport code.

### `internal/memory`

The current pack store is useful, but the rewrite needs more than file IO.

Add:

- retrieval request and result types
- excerpt records
- source catalog contracts
- indexer contracts
- clear subject-hint behavior

### `internal/operator`

This package should stay transport-neutral.

It should answer questions like:

- what conversations need review
- what runs happened recently
- what knowledge packs are available
- what operator approvals are pending

### `internal/webui`

Do not start with a frontend framework fight.

Start with backend route contracts and one authenticated HTTP server boundary. The frontend can be a later layer over the same operator service.

## Rewrite Pass Order

### Pass 1: Direction, contracts, and scaffolding

- land ADR 0006
- add the orchestrator, operator, gateway, memory, and web UI contracts
- align the `plan` and `contracts` surfaces with the narrowed direction

### Pass 2: Knowledge and evidence layer

- implement source catalog records
- add retrieval and excerpt selection
- attach references to replies
- write tests around subject selection and citation formatting

### Pass 3: Discord orchestration rewrite

- replace the direct reply-engine path with the orchestrator program
- persist run records and session updates
- expose review-required turns cleanly

### Pass 4: Operator backend

- implement operator service stores
- expose authenticated HTTP endpoints
- add dashboard snapshot, review queue, and run detail routes

### Pass 5: Operator frontend

- build the minimal console UI
- show conversation feed, references, route choice, and review queue
- keep authentication mandatory

### Pass 6: Expansion after coherence

- WhatsApp
- iMessage
- deeper browser automation
- mobile companion flows

## Validation Rules

Every pass should answer these questions clearly.

1. What is the new boundary?
2. What config or contract exposes it?
3. How does an operator inspect it?
4. What tests prove it works?
5. What claims can now be made in public tense?

## Public Claims Policy For The Rewrite

### Safe present-tense claims

- Discord is the primary messaging surface.
- Savitar runs Gemma 4 locally through Ollama.
- The system is being rewritten around curated knowledge, live evidence, and an operator console.
- Future transport expansion is planned after the core workflow is complete.

### Claims that must stay future-tense until shipped

- WhatsApp support.
- iMessage support.
- Multi-surface unified dashboard.
- Public demo mode without authentication.
- Pixel Fold operator workflows.
- LiteRT on-device inference.

## Faculty Collaboration Language

The rewrite should preserve room for real academic collaboration without overstating it.

Recommended language pattern:

> Savitar has been informed by conversations with Illinois State University faculty including Sally Xie, Mangolika Bhattacharya, Rudro Baksi, and Somnath Lahiri. Those conversations have ranged from informal discussion to more extended information sharing. Public-facing materials should describe those roles precisely and should not imply institutional endorsement without explicit approval.

Operational rule:

- keep exact collaboration notes in local working memory or approved project documentation
- distinguish conversation, feedback, collaboration, and endorsement
- never upgrade one of those categories in public copy without permission

## Hackathon Narrative After The Rewrite

The strongest honest story is no longer "one agent on every surface."

It is:

**One local-first knowledge operator that helps Illinois communities, faculty, students, and small businesses access trusted guidance through a working Discord surface and an auditable operator console.**

That story is narrower, but it is more defensible and more likely to survive technical scrutiny.