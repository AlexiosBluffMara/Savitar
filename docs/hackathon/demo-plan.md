# Live Demo Plan

This file distinguishes the current shipped state from the public-facing submission target.

## Current shipped state

- Discord reply path works locally.
- The operator dashboard prototype works locally in demo mode.
- Local markdown retrieval and evidence-style fallback exist.
- Public authenticated web access is not finished yet.

The current repo can already support a truthful local demo video. The remaining work is turning that into a public judge-facing live demo without pretending unfinished auth and moderation flows already exist.

## Submission target

The target public demo should show one coherent system:

```text
Invite-only Discord
  -> Savitar gateway
  -> local Gemma 4 on Ollama
  -> local markdown and document packs on SSD
  -> operator dashboard with evidence and review state
```

No part of the core inference story should depend on hosted-model fallback.

## Judge scenarios

### Scenario 1: Chicago small-business operator

- A business owner asks about a grant, filing step, or city resource.
- Savitar answers from local evidence packs.
- The dashboard shows the route, evidence, and operator visibility.

### Scenario 2: Educator Live Mode

- An educator runs a moderated class or workshop channel.
- Students ask questions against a local pack.
- The educator sees what sources Savitar used and can intervene.

### Scenario 3: Safety and review

- An answer triggers operator review or gets corrected.
- The dashboard makes the review trail visible.
- This demonstrates that Savitar is designed for oversight, not blind automation.

## Public demo posture

### What should be public

- The live demo URL.
- A read-only or sanitized judge view.
- The GitHub repository.
- The short video.

### What should stay private

- Real operator channels.
- Sensitive knowledge packs.
- Any personal or business data used during rehearsals.

## Demo recording priorities

1. Show the human problem first.
2. Show Discord as the working entry point.
3. Show the dashboard as the operator proof surface.
4. Show that Gemma 4 is local on the Mac Mini.
5. End on why this matters for education, equity, and trust.