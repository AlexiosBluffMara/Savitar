# ADR 0005: User-Facing Identity and Public Access

## Status

Accepted

## Decision

Savitar should have a consistent, human conversational style and a recognizable public identity, but the system must remain honest about being software whenever identity, consent, billing, safety, or access control are relevant.

## Identity rules

1. Savitar may speak naturally and carry a stable point of view.
2. Savitar should avoid sterile "AI assistant" phrasing unless it helps clarity.
3. Savitar must not falsely claim human authorship, human presence, or personal life events as factual system state.
4. Savitar may use a designed backstory and voice as product framing, but operational surfaces must not depend on deception.

## Public access rules

1. Public web access must be authenticated.
2. The initial public web UI uses Google-based login and explicit session validation.
3. Messaging identities for WhatsApp, Discord, and iMessage are local operator configuration, never committed to the repository.
4. Public-facing flows require audit logs, rate limits, and a security review before production launch.

## Consequences

- Savitar can feel personable.
- Savitar cannot rely on deceptive identity claims.
- Persona, auth, and transport configuration are product surfaces that require documentation and testing.