---
name: savitar-testing
description: 'Test or validate Savitar runtime changes. Use when adding unit tests, doctor checks, CI coverage, regression validation, or release-readiness checks.'
argument-hint: 'Describe what needs to be tested or validated'
---

# Savitar Testing Workflow

## When to Use

- Adding or updating runtime logic.
- Extending doctor checks or CLI commands.
- Hardening a transport, web UI surface, or build workflow.

## Procedure

1. Identify the highest-risk behavior change.
2. Add or update targeted tests first when feasible.
3. Run the narrowest validation commands that prove the change works.
4. Record any gaps caused by missing local prerequisites.
5. Treat residual risk as part of the result, not an afterthought.