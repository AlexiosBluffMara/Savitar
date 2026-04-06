---
name: savitar-security
description: 'Review Savitar for security risks. Use when working on public web UI, Google auth, messaging transports, remote execution, MCP/tool safety, secret handling, or prompt-injection-resistant workflows.'
argument-hint: 'Describe the surface or change to review'
---

# Savitar Security Workflow

## When to Use

- Exposing a new public endpoint or login flow.
- Adding a messaging adapter.
- Expanding tool execution, MCP access, or remote host control.

## Procedure

1. Identify the assets, trust boundaries, and attacker entry points.
2. Review auth, session, transport, and secret-handling assumptions.
3. Check for tool abuse, privilege escalation, prompt injection, or data exfiltration paths.
4. Rank findings by severity and exploitability.
5. Recommend mitigations that fit Savitar's lightweight architecture.