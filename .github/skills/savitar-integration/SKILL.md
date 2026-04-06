---
name: savitar-integration
description: 'Add or change a Savitar model route, MCP server, transport, browser automation feature, or memory workflow. Use when implementing new integrations, extending contracts, or wiring new capabilities into the runtime.'
argument-hint: 'Describe the integration or capability you are adding'
---

# Savitar Integration Workflow

## When to Use

- Adding a new MCP server.
- Adding or changing Discord, WhatsApp, or iMessage transport behavior.
- Extending browser automation, shell execution, repository analysis, or memory capture.
- Changing model routing or agent persona policy.

## Procedure

1. Read `docs/adr/0001-clean-room-charter.md` through `docs/adr/0004-integrations-and-guardrails.md`.
2. Update `config/savitar.example.json` so the new capability has an obvious local configuration path.
3. Update `.vscode/mcp.json` if the capability depends on a workspace MCP server.
4. Extend the relevant contract or runtime doctor checks so the capability has a visible boundary.
5. Add or update setup and roadmap docs when the environment or sequencing changes.
6. Add a provenance entry under `docs/provenance/` when the work depends on external products, public documentation, or parity targets.
7. Add or update tests before calling the capability complete.