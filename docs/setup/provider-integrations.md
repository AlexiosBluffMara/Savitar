# Provider Integrations

Savitar now supports a local-first integration path for local Ollama, GitHub, Hugging Face, Kaggle, and Discord.

When the local Ollama lane is configured, the Discord preview and live bot surfaces use that same host for normal conversational replies. Guild channels and direct messages both fail closed unless a local Ollama target exists.

These helper scripts are repo code. Review the checkout you are running before entering or exporting any live credentials.

## Goal

- Keep secrets out of version control.
- Use the macOS Keychain for storage.
- Export credentials into the current shell only when the operator intentionally loads them.

## Local scripts

- `./scripts/store-local-integrations.sh` prompts for tokens and stores them in the macOS Keychain.
- `source ./scripts/use-local-integrations.sh` exports provider environment variables from the Keychain into the current shell.
- `./scripts/validate-local-integrations.sh` runs provider-only authentication checks without invoking the Savitar binary.
- `./bin/savitar discord preview "<message>"` exercises the model-backed reply path against the reviewed binary.

## Trusted CLI locations

- Install `hf` and `kaggle` outside the workspace, for example with `uv tool install huggingface_hub` and `uv tool install kaggle`.
- The validator prefers binaries in `/opt/homebrew/bin`, `/usr/local/bin`, and `~/.local/bin` before it falls back to the ambient shell path.

## Environment variables

- GitHub: `GH_TOKEN` and `GITHUB_TOKEN`
- Hugging Face: `HF_TOKEN`
- Kaggle: `KAGGLE_API_TOKEN`
- Discord: `SAVITAR_DISCORD_BOT_TOKEN`

## Keychain service names

- `savitar/github/token`
- `savitar/huggingface/token`
- `savitar/kaggle/api-token`
- `savitar/discord/bot-token`

## Local config

The ignored local override at `config/savitar.local.json` enables these integration surfaces without storing any credential values. It relies on the environment variables loaded from the Keychain-backed shell helper.

The local Ollama lane is configured through `models.localDefault` rather than `integrations`. Keep the provider set to `ollama`, the endpoint pointed at your local host, and the model tag aligned with what is actually pulled on the machine.

## Validation flow

1. Review the current checkout.
2. Run `export SAVITAR_REVIEWED_CHECKOUT=1` in the shell you trust.
3. Run `./scripts/store-local-integrations.sh`.
4. Source `./scripts/use-local-integrations.sh`.
5. Run `./scripts/validate-local-integrations.sh`.
6. Build `./bin/savitar` from the same reviewed checkout.
7. Run `./bin/savitar integrations` and `./bin/savitar discord status`.
8. Run `./bin/savitar discord preview "What can you help with?"` to confirm the conversational reply path.
9. Run `./bin/savitar discord run` when you are ready to bring the bot online.

## Security notes

- Do not commit provider tokens, PATs, or bot tokens.
- Do not store or expect any hosted Ollama token anywhere in this repo; Savitar now assumes local Ollama only.
- Prefer the smallest usable scope for every token.
- Keep token handling inside reviewed local shells, the macOS Keychain, and trusted provider CLIs.