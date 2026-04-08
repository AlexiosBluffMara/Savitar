# Phase Plan

This file is the legacy broad-scope phase map.

The active rewrite direction is defined in [docs/adr/0006-rewrite-product-direction.md](../adr/0006-rewrite-product-direction.md) and [docs/roadmap/0006-rewrite-blueprint.md](0006-rewrite-blueprint.md). Use those files when planning the narrowed product around Discord, knowledge retrieval, and the operator console.

## Phase 0: Foundation

- Clean-room charter and provenance process
- Mac Mini setup guide and bootstrap scripts
- Workspace MCP configuration
- Minimal Go CLI and runtime doctor

## Phase 1: Core loop

- Command router
- model routing implementation
- task envelope and session state
- structured logging
- direct CLI status, doctor, contracts, and plan surfaces
- persona and reply-style policy

## Phase 2: Tool surfaces

- MCP orchestration layer
- controlled shell execution policy
- browser session management
- repository intake and analysis pipeline
- screen capture and recording control
- hooks, guardrails, and tool permission policy
- skills inventory and local skill loading

## Phase 3: Messaging transports

- Discord bot runtime
- WhatsApp bridge contract
- iMessage bridge contract
- inbound and outbound message normalization
- messaging gateway process
- remote host execution contract
- initial public web UI backend and operator session model

## Phase 4: Memory packs

- subject-area memory dumps
- indexed summaries and retrieval
- replayable task context
- source-repository feature intake summaries
- project and user memory policies
- persona consistency across CLI, messaging, and web surfaces

## Phase 5: Hardening

- long-running service mode
- retry policy
- operator controls
- traceable audit logs
- setup wizard and migration helpers
- packaging and release discipline
- Google OAuth validation and public internet threat hardening

## Phase 6: Packaging and deployment

- GitHub Actions artifact builds
- release process
- launchd service packaging
- Mac Mini rollout checklist
- Pixel 9 Pro Fold mobile operator dashboard (responsive web UI)
- Google Fi as WhatsApp Business carrier number
- Gemma 4 E2B on-device inference via LiteRT on Tensor G4 (exploratory)
- Apple Education procurement documentation for school deployments