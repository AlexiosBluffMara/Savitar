#!/usr/bin/env bash

if ! (return 0 2>/dev/null); then
	echo "source this script instead of executing it: . ./scripts/use-local-integrations.sh" >&2
	exit 1
fi

if [[ "${SAVITAR_REVIEWED_CHECKOUT:-}" != "1" ]]; then
	echo "refusing to export live credentials until SAVITAR_REVIEWED_CHECKOUT=1 is set for a reviewed checkout" >&2
	return 1
fi

load_secret() {
	local service="$1"
	/usr/bin/security find-generic-password -a "$USER" -s "$service" -w 2>/dev/null || true
}

export_from_keychain() {
	local env_name="$1"
	local service="$2"
	local value
	value="$(load_secret "$service")"
	if [[ -n "$value" ]]; then
		export "${env_name}=${value}"
	fi
}

export_from_keychain GH_TOKEN savitar/github/token
export_from_keychain HF_TOKEN savitar/huggingface/token
export_from_keychain KAGGLE_API_TOKEN savitar/kaggle/api-token
export_from_keychain SAVITAR_DISCORD_BOT_TOKEN savitar/discord/bot-token

if [[ -n "${GH_TOKEN:-}" && -z "${GITHUB_TOKEN:-}" ]]; then
	export GITHUB_TOKEN="$GH_TOKEN"
fi

export KAGGLE_CONFIG_DIR="${KAGGLE_CONFIG_DIR:-$HOME/.kaggle}"
export HF_HOME="${HF_HOME:-$HOME/.cache/huggingface}"