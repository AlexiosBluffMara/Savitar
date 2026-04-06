---
name: savitar-provenance
description: 'Create or update a Savitar provenance log entry for clean-room work. Use when a feature depends on external documentation, parity targets, third-party products, or public APIs and the repo needs an auditable record of what was consulted.'
argument-hint: 'Describe the feature or workstream that needs a provenance entry'
---

# Savitar Provenance Workflow

## When to Use

- Starting a feature inspired by another system.
- Studying public API or setup documentation before implementation.
- Recording the source material behind an integration or deployment change.

## Procedure

1. Copy the template in `./assets/entry-template.md` into a new file under `docs/provenance/`.
2. Record the public sources, product docs, or requirements that informed the work.
3. State clearly what was not copied and which implementation details were authored from scratch.
4. Link the provenance entry from the relevant ADR, roadmap item, or pull request.
5. Update the entry if the scope changes or new external references are consulted.