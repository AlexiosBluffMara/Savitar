# MCP Stack

Savitar starts with a small MCP stack that covers repository context, live documentation, browser automation, and macOS automation.

## Workspace default

The committed workspace configuration lives in `.vscode/mcp.json`.

## Core servers

| Server | Purpose | Default mode | Configuration |
| --- | --- | --- | --- |
| GitHub MCP | Repository, issues, PRs, Actions, and GitHub-native context | Remote | `https://api.githubcopilot.com/mcp/` |
| Context7 | Live, version-specific documentation lookup | Local `npx` | `@upstash/context7-mcp` |
| Playwright MCP | Controlled browser sessions and automation | Local `npx` | `@playwright/mcp@latest` |
| Peekaboo | macOS-native UI inspection and control | Local `npx` | `@steipete/peekaboo` |

## Recommended expansion after phase 1

- Slack MCP for team notifications and operator workflows.
- Postgres MCP for durable operational state if Savitar outgrows flat-file memory.
- A read-only filesystem MCP when Savitar needs broader visibility than the active workspace allows.

## GitHub MCP note

- The remote GitHub MCP server is the default because it avoids local binary management and works well with VS Code.
- If a local fallback is needed later, use the official Docker image from `ghcr.io/github/github-mcp-server` with a scoped PAT and explicit toolsets.

## Guardrails

- Keep the server set small enough that tool choice remains clear.
- Prefer read-only or smallest-scope configurations first.
- Add new servers only when the capability is part of the roadmap, not because the registry exists.