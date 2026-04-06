#!/usr/bin/env bash

set -euo pipefail

required=(brew go node npm npx gh)
optional=(jq docker python3 uv adb osascript)
status=0

check_command() {
	local name="$1"
	local required_flag="$2"
	if command -v "$name" >/dev/null 2>&1; then
		printf "[ok] %s -> %s\n" "$name" "$(command -v "$name")"
		return
	fi

	if [[ "$required_flag" == "required" ]]; then
		printf "[missing] %s\n" "$name"
		status=1
	else
		printf "[optional-missing] %s\n" "$name"
	fi
}

for name in "${required[@]}"; do
	check_command "$name" required
done

for name in "${optional[@]}"; do
	check_command "$name" optional
done

exit "$status"