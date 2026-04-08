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
OLLAMA_BIN="$(trusted_bin ollama)"
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

if [[ -z "$OLLAMA_BIN" || ! -x "$OLLAMA_BIN" ]]; then
	echo "missing Ollama CLI" >&2
	exit 1
fi

"$OLLAMA_BIN" list >/dev/null
echo "ollama: ok"

"$GH_BIN" auth status >/dev/null
echo "github: ok"

"$HF_BIN" auth whoami >/dev/null
echo "huggingface: ok"

"$KAGGLE_BIN" competitions list --page-size 1 -v >/dev/null
echo "kaggle: ok"