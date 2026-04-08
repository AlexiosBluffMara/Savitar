# Web UI Prototype Runbook

This runbook covers the current Savitar operator web UI prototype.

It is a local demo surface for hackathon validation, not a public deployment.

## Required now

- A working Go toolchain on the local host.
- A Savitar checkout with the normal local config flow.
- A reviewed shell if you plan to demonstrate live provider integrations.
- Optional but useful: initialize the local session once so the dashboard has a current session record.

## Current prototype status

- The backend and HTML dashboard exist as a launchable prototype surface inside `internal/webui` and the `savitar webui serve --demo` command.
- The dashboard shows real runtime state, model routing, integration readiness, and repo-backed knowledge fallback.
- Public authentication is still missing. Google OAuth, signed sessions, rate limits, and audit logging remain pending.

## Start the prototype

- `./bin/savitar webui serve --demo --addr :8080`
- Or without a prior build: `go run ./cmd/savitar webui serve --demo --addr :8080`
- Open `http://127.0.0.1:8080` and keep it local-only.

For package-level validation without running a long-lived process, use `go test ./internal/webui`.

## What the prototype shows

- Runtime summary: current local model, enabled surfaces, MCP count, agents, and skills.
- Model routing: the local Gemma lane plus the managed Copilot lanes.
- Knowledge view: on-disk knowledge packs when present, or repo markdown fallback from the configured evidence directories when the pack store is empty.
- Recent runtime state: the local session record and any queued review JSON files.
- Hackathon budget selector: the Mac Mini recommendation, the lower-cost local path, the cloud tradeoff, and the Pixel Fold companion story.

## Security posture right now

- Demo mode is intended for trusted local use only.
- Do not expose the current prototype through Cloudflare Tunnel, router port-forwarding, or a public reverse proxy.
- Public judge access needs a separate authenticated mode that complies with `docs/adr/0005-user-facing-identity.md`.

## Troubleshooting

- Empty knowledge panel: no files were found in the configured knowledge directory, so the UI fell back to repo markdown. Add packs under `.savitar/knowledge` if you want the pack view.
- No recent runtime state: run `savitar session init` once, or wait until the orchestrator starts writing durable run records.
- Empty review queue: no review-item JSON files exist under `.savitar/reviews` yet.
- `webui.enabled` still shows as planned in the surface table: set it locally only when you want the gateway plan to advertise the web surface by default.