# Provenance Entry

- Workstream: Local-only provider integrations for Ollama Cloud, GitHub, Hugging Face, Kaggle, and Discord token handling
- Date: 2026-04-06
- Author: GitHub Copilot
- External systems or docs consulted: Ollama Cloud authentication docs, GitHub CLI environment and auth docs, Hugging Face Hub quickstart auth docs, Kaggle public API docs, Kaggle CLI documentation
- Public interfaces or behaviors captured: `OLLAMA_API_KEY` bearer auth for `https://ollama.com/api`, `GH_TOKEN` and `gh auth` credential store behavior, `HF_TOKEN` env auth, `KAGGLE_API_TOKEN` support in the official Kaggle CLI, and local token-loading patterns through shell environment variables
- Explicit non-goals or material not copied: No provider SDK implementations, third-party auth handlers, or external scripts were copied into Savitar
- New files or modules created from scratch: provider integration config surfaces, runtime integration status output, keychain-backed shell helpers, provider setup docs, and local validation workflow
- Verification completed: runtime tests, CLI help/status coverage, provider auth documentation review, and local script validation planning
- Follow-up provenance needed: any future model-router requests against Ollama Cloud, GitHub API write actions, or Kaggle submission automation