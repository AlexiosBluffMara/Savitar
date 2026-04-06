# Mac Mini Setup

This document defines the baseline Savitar workstation setup for macOS on Apple Silicon.

## Required tools

- Homebrew
- Go
- Node.js and npm
- GitHub CLI
- Docker Desktop or Docker Engine
- `jq`
- `uv`
- Android platform tools for the Pixel Fold bridge

## Bootstrap path

1. Run `./scripts/bootstrap-macos.sh`.
2. Run `./scripts/verify-prereqs.sh`.
3. Authenticate GitHub CLI with `gh auth login`.
4. Open the workspace in VS Code so `.vscode/mcp.json` can be discovered.
5. Copy `config/savitar.example.json` to `config/savitar.local.json` and fill local secrets.
6. Keep Google OAuth, iMessage account details, and messaging identifiers in local config or the host keychain only.
7. Review the current checkout, then export `SAVITAR_REVIEWED_CHECKOUT=1` in the shell you trust.
8. Store provider tokens with `./scripts/store-local-integrations.sh` and source `./scripts/use-local-integrations.sh` in that reviewed shell.

## Local model lane

- Default target: an MLX-backed Gemma 4n E4B-class model.
- Install the local model tooling with `uv` once the exact approved model artifact is chosen.
- Keep the local model endpoint configurable through `config/savitar.local.json`.

## macOS permissions

- Grant Screen Recording for Peekaboo.
- Grant Accessibility for Peekaboo and future iMessage automation.
- Grant Automation permissions for the terminal or editor host that will run Savitar.

## Device bridge preparation

- Enable developer mode and USB debugging on the Pixel Fold.
- Confirm `adb devices` works before building the WhatsApp bridge.
- Keep any real phone number in local config or the host keychain, never in the repo.
- Treat the eventual Pixel Fold native companion as a later phase, not a current setup dependency.

## Discord bot setup notes

- Create a Discord application and bot in the Discord developer portal.
- Keep the bot token in the shell or host keychain only, exposed as `SAVITAR_DISCORD_BOT_TOKEN`.
- Put any trusted remote operator Discord user IDs only in local config under `transports.discord.operatorUserIDs`.
- The default Savitar Discord mode is mention-or-DM, which works without channel-wide message capture.
- If you want the bot to read all allowed guild messages, set `transports.discord.useMessageContentIntent` to `true`, keep a non-empty `allowedChannelIDs` list, and enable the privileged `MESSAGE_CONTENT` intent in the Discord developer portal.
- Validate the local surface with `savitar discord status` and `savitar discord preview status` before running `savitar discord run`.

## Provider integration notes

- Ollama Cloud uses `OLLAMA_API_KEY` against `https://ollama.com/api`.
- GitHub CLI prefers its own credential store, but Savitar also supports `GH_TOKEN` for explicit token-based use.
- Hugging Face uses `HF_TOKEN`; keep it scoped to the smallest permission set that still works.
- Kaggle's current official CLI supports `KAGGLE_API_TOKEN`, which can stay in the macOS Keychain and be exported only when needed.
- Install trusted user-level `hf` and `kaggle` CLIs outside the workspace before running the provider validator.
- Use `./scripts/validate-local-integrations.sh` after sourcing the integration loader to check the provider auth path without running workspace code under live tokens.
- Build and run `./bin/savitar` from a reviewed checkout before using repo-local commands with exported credentials.

## Public web setup notes

- The future public web UI will use Google-based login.
- Keep Google client secrets and web session secrets out of version control.
- Do not put any real iCloud or messaging account address into committed config examples.

## Validation

- `./scripts/verify-prereqs.sh`
- `make doctor`
- `make test`
- `make build`