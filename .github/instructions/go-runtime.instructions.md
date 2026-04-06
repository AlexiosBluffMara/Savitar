---
name: "Savitar Go Runtime"
description: "Use when writing Savitar Go runtime code, CLI commands, model routing, contracts, config loading, or transport foundations."
applyTo: ["cmd/**/*.go", "internal/**/*.go"]
---
# Savitar Go Runtime

- Keep packages small and explicit.
- Prefer direct data flow over framework-heavy abstractions.
- When adding a capability, expose it through config, contracts, and a CLI or doctor surface when practical.
- Add focused tests for routing, defaults, parsing, or doctor behavior when runtime logic changes.
- Avoid hidden global state except for immutable defaults.