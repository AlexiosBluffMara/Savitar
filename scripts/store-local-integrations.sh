#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"

if [[ "${SAVITAR_REVIEWED_CHECKOUT:-}" != "1" ]]; then
	echo "set SAVITAR_REVIEWED_CHECKOUT=1 only after reviewing this checkout before entering live credentials" >&2
	exit 1
fi

store_secret() {
	local service="$1"
	local prompt="$2"
	echo "storing $service in the macOS Keychain" >&2
	/usr/bin/security add-generic-password -U -a "$USER" -s "$service" -w
}

store_secret savitar/github/token "GitHub token or PAT"
store_secret savitar/huggingface/token "Hugging Face token"
store_secret savitar/kaggle/api-token "Kaggle API token"
store_secret savitar/discord/bot-token "Discord bot token"