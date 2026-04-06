# ADR 0001: Clean-Room Charter

## Status

Accepted

## Decision

Savitar is a clean-room project. We can study public behavior, interfaces, and documentation for other systems, but we do not copy source code, internal data structures, prompts, test suites, or unpublished implementation details from Hermes, Picoclaw, OpenClaw, or any other parity target.

## Required practices

1. Record external inputs in `docs/provenance/` before or during implementation.
2. Describe parity goals in behavior terms, not file-for-file imitation.
3. Keep new code authored from scratch in this repository.
4. Treat secrets, phone numbers, API tokens, Apple IDs, and device identifiers as local configuration only.
5. Prefer public docs, official APIs, and vendor-supported bridges over reverse-engineered shortcuts.

## Prohibited practices

1. Copying code or prompts from external repositories.
2. Porting tests or fixtures from parity targets.
3. Committing private credentials or personal contact details.
4. Calling the project complete without updating provenance when external material informed the work.

## Consequences

- Feature parity is allowed.
- Code reuse is not.
- Provenance is an operating process, not a one-time note.