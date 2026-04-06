# ADR 0003: Model Routing Policy

## Status

Accepted

## Decision

Savitar will route work across four logical lanes.

1. `local-default`: the preferred local lane on Apple Silicon, backed by MLX and a Gemma 4n E4B-class model.
2. `copilot-0x`: the routine Copilot lane for low-risk, high-volume work.
3. `copilot-0.33x`: the medium-complexity Copilot lane.
4. `copilot-1x`: the high-complexity Copilot lane for expensive reasoning.

## Routing rules

- Prefer the local lane for latency-sensitive work, private context, and tasks that do not need a frontier cloud model.
- Prefer `copilot-0x` for routine drafting, classification, summarization, and low-risk execution.
- Escalate to `copilot-0.33x` for synthesis across multiple sources or medium-depth planning.
- Escalate to `copilot-1x` for architecture changes, debugging, multi-step code generation, or ambiguous tasks with higher failure cost.

## Local model preference

- On Apple Silicon, prefer MLX-backed Gemma routes when available.
- Prefer an E4B-class Gemma lane over a smaller E2B-class lane unless memory pressure forces a fallback.
- Keep the exact model artifact name configurable outside source control.

## Consequences

Routing policy is part of product behavior and should remain visible in both code and docs.