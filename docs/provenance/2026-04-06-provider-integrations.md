# Provenance Entry

- Workstream: Local-first provider integrations for local Ollama, GitHub, Hugging Face, Kaggle, and Discord token handling
- Date: 2026-04-06
- Author: GitHub Copilot
- External systems or docs consulted: Ollama local API docs, GitHub CLI environment and auth docs, Hugging Face Hub quickstart auth docs, Kaggle public API docs, Kaggle CLI documentation
- Public interfaces or behaviors captured: local Ollama CLI and `POST /api/chat` usage on a reviewed host, `GH_TOKEN` and `gh auth` credential store behavior, `HF_TOKEN` env auth, `KAGGLE_API_TOKEN` support in the official Kaggle CLI, and local token-loading patterns through shell environment variables
- Explicit non-goals or material not copied: No provider SDK implementations, third-party auth handlers, or external scripts were copied into Savitar
- New files or modules created from scratch: provider integration config surfaces, runtime integration status output, keychain-backed shell helpers, provider setup docs, and local validation workflow
- Verification completed: runtime tests, CLI help/status coverage, provider auth documentation review, and local script validation planning
- Follow-up provenance needed: any future GitHub API write actions, Kaggle submission automation, or changes to local Ollama deployment assumptions