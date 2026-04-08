# Provenance Entry

- Workstream: Ollama-backed conversational replies for the Savitar Discord surface
- Date: 2026-04-06
- Author: GitHub Copilot
- External systems or docs consulted: Ollama API introduction docs, Ollama `/api/chat` request and response examples, Ollama error-format docs
- Public interfaces or behaviors captured: `POST /api/chat` with `model`, `messages`, and `stream=false`; local Ollama API base URLs; JSON error responses with an `error` field
- Explicit non-goals or material not copied: No Ollama SDK code, prompts, tests, or third-party Discord bot logic were copied into Savitar
- New files or modules created from scratch: `internal/chat`, `internal/ollama/chat.go`, Ollama chat tests, conversational responder history logic, and Discord backpressure controls
- Verification completed: provider validation script passed for Ollama, GitHub, Hugging Face, and Kaggle; `go test ./...` passed; `go build -o ./bin/savitar ./cmd/savitar` passed; `./bin/savitar discord preview "What can you help with right now?"` returned a live Ollama-backed reply; `./bin/savitar discord run` connected successfully as `Artemis#1566` in one guild after the conversational and throttling changes
- Follow-up provenance needed: any future long-lived conversation storage, moderation pipeline, or multi-provider reply routing beyond Ollama