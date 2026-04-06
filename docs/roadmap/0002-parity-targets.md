# Parity Targets

Savitar is aiming for behavioral parity across capability categories, not implementation parity with any other codebase.

## Target categories

- Multi-lane model routing with explicit local and Copilot-backed policies
- MCP-native tooling for GitHub, live documentation, browser automation, and macOS control
- Messaging surfaces for Discord, WhatsApp, and iMessage
- A public web UI with Google-based login and operator validation
- Controlled code execution and repository analysis
- Memory dumps and subject-area knowledge packs
- CLI-first and scriptable operator surfaces
- Screen capture, screenshots, and controlled remote host usage
- A later native companion path for Android-device actions
- A distinct, opinionated response style with explicit design choices

## Constraints

- Keep the core runtime lightweight and responsive.
- Prefer readable code over framework-heavy shortcuts.
- Preserve clean-room separation from external projects.
- Treat production readiness as a later hardening phase, not an excuse for overbuilding phase 0.

## Source repos informing the target

- `ultraworkers/claw-code`: CLI ergonomics, sessions, slash commands, permissions, MCP inspection, hooks, skills, and JSON output surfaces.
- `AlexiosBluffMara/claw-code`: clean-room porting discipline, inventory and audit views, and Python-friendly manifest/reporting patterns.
- `AlexiosBluffMara/picoclaw`: lightweight Go runtime expectations, messaging gateway breadth, model routing, skills, and deploy-anywhere thinking.
- `AlexiosBluffMara/Artemis`: gateway-first messaging, memory and skills depth, setup and migration surfaces, remote execution backends, and broader production ambition.

These sources are parity references only. Savitar implementations remain authored from scratch in this repository.