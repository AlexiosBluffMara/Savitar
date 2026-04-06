---
name: savitar-architecture
description: 'Plan or revise Savitar architecture, ADRs, contracts, system boundaries, phase sequencing, or multi-surface integration. Use when scoping features, defining interfaces, or breaking work into clean-room milestones.'
argument-hint: 'Describe the subsystem or feature to design'
---

# Savitar Architecture Workflow

## When to Use

- Designing a new subsystem.
- Revising runtime boundaries.
- Planning transport, web UI, auth, memory, or deployment work.

## Procedure

1. Read the ADRs and `docs/roadmap/0003-source-feature-matrix.md`.
2. Define the smallest boundary that delivers the feature.
3. Update contracts, config examples, and roadmap notes before implementation details.
4. List dependencies, guardrails, and validation steps.
5. Add provenance if public docs or parity references informed the design.