#!/usr/bin/env bash

set -euo pipefail

if ! command -v brew >/dev/null 2>&1; then
	echo "Homebrew is required before running this bootstrap script." >&2
	exit 1
fi

formulae=(go jq uv android-platform-tools gh node python)

for formula in "${formulae[@]}"; do
	if brew list "$formula" >/dev/null 2>&1; then
		echo "already installed: $formula"
		continue
	fi

	brew install "$formula"
done

if brew list --cask docker >/dev/null 2>&1; then
	echo "already installed: docker"
elif command -v docker >/dev/null 2>&1; then
	echo "docker CLI already available: $(command -v docker)"
elif [[ -d "/Applications/Docker.app" ]]; then
	echo "Docker.app already present in /Applications"
else
	brew install --cask docker
fi

echo
echo "Bootstrap complete."
echo "Next steps:"
echo "  1. gh auth login"
echo "  2. Copy config/savitar.example.json to config/savitar.local.json"
echo "  3. Open the workspace in VS Code so .vscode/mcp.json is discovered"
echo "  4. Run ./scripts/verify-prereqs.sh"