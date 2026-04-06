# Deployment Notes

Savitar is not production-ready yet, but the repo now has a deployment path that scales from local Mac Mini usage to build automation.

## Local runtime target

- Primary host: macOS on Apple Silicon.
- Primary process model: a local Savitar binary plus workspace MCP servers.
- Primary operations mode: launch manually during development, then promote to `launchd` once the agent loop is stable.

## GitHub Actions target

- `ci.yml` will run formatting and tests on every push and pull request.
- `build-artifacts.yml` will build versioned binaries for supported targets and upload them as workflow artifacts.
- Release signing, notarization, and installer packaging are intentionally deferred until the daemon and transport bridges stabilize.
- The future public web UI will need its own deployment lane with Google OAuth secrets, session secrets, rate limiting, and reverse-proxy hardening.

## Production hardening gates

1. Split secrets from local config cleanly.
2. Add structured logging and traceable run IDs.
3. Add transport-specific integration tests.
4. Add a supervised service mode for the Mac Mini host.
5. Add backup and restore procedures for memory snapshots.