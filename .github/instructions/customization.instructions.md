---
name: "Savitar Customization Files"
description: "Use when creating or editing Savitar custom agents, instructions, skills, hooks, prompts, or workspace MCP configuration."
applyTo: [".github/**/*.md", ".github/**/*.json", ".vscode/**/*.json"]
---
# Savitar Customization Files

- Keep each custom agent single-purpose with the smallest useful toolset.
- Make skill and agent descriptions keyword-rich so delegation works reliably.
- Keep hooks short, auditable, and free of secrets.
- Prefer workspace-scoped customization for Savitar-specific behavior.
- Avoid circular handoffs between agents.