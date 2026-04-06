#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"

trusted_bin() {
	local name="$1"
	for prefix in /opt/homebrew/bin /usr/local/bin "$HOME/.local/bin" /usr/bin /bin; do
		if [[ -x "$prefix/$name" ]]; then
			echo "$prefix/$name"
			return 0
		fi
	done
	command -v "$name"
}

PYTHON_BIN="/usr/bin/python3"
HF_BIN="$(trusted_bin hf)"
KAGGLE_BIN="$(trusted_bin kaggle)"
GH_BIN="$(trusted_bin gh)"
source "$SCRIPT_DIR/use-local-integrations.sh"

if [[ ! -x "$PYTHON_BIN" ]]; then
	echo "missing Python runtime at $PYTHON_BIN" >&2
	exit 1
fi

if [[ ! -x "$HF_BIN" ]]; then
	echo "missing Hugging Face CLI at $HF_BIN" >&2
	exit 1
fi

if [[ ! -x "$KAGGLE_BIN" ]]; then
	echo "missing Kaggle CLI at $KAGGLE_BIN" >&2
	exit 1
fi

if [[ -z "$GH_BIN" || ! -x "$GH_BIN" ]]; then
	echo "missing GitHub CLI" >&2
	exit 1
fi

if [[ -n "${OLLAMA_API_KEY:-}" ]]; then
	"$PYTHON_BIN" - <<'PY'
import os
import urllib.request

request = urllib.request.Request(
    "https://ollama.com/api/tags",
    headers={"Authorization": f"Bearer {os.environ['OLLAMA_API_KEY']}"},
)
with urllib.request.urlopen(request, timeout=30) as response:
    response.read(1)
PY
	echo "ollama: ok"
else
	echo "ollama: missing OLLAMA_API_KEY" >&2
	exit 1
fi

"$GH_BIN" auth status >/dev/null
echo "github: ok"

"$HF_BIN" auth whoami >/dev/null
echo "huggingface: ok"

"$KAGGLE_BIN" competitions list --page-size 1 -v >/dev/null
echo "kaggle: ok"