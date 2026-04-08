# Submission Strategy — Gemma 4 Good Hackathon

This document preserves the original broad submission framing.

The active rewrite narrows the first complete product to Discord, curated knowledge, evidence-backed replies, and an authenticated operator console. See [docs/adr/0006-rewrite-product-direction.md](../adr/0006-rewrite-product-direction.md) and [docs/roadmap/0006-rewrite-blueprint.md](../roadmap/0006-rewrite-blueprint.md) for the current build direction.

## Positioning

**Project name:** Savitar  
**Tagline:** A local-first, multi-surface AI agent for community knowledge, small business support, and academic collaboration — running Gemma 4 on Ollama.

**Duo submission:** AlexiosBluffMara with Illinois State University affiliation.

**Developer hardware:** Mac Mini M4 (16GB, server/inference), Google Pixel 9 Pro Fold (mobile companion, on Google Fi), Apple + Google ecosystem coverage.

---

## Target Tracks

### Primary: Main Track ($50K / $25K / $15K / $10K)

Savitar is a fully-functional, open-source Go runtime that routes conversations across Discord, WhatsApp, and iMessage through a unified gateway, backs them with Gemma 4 running locally on Ollama, and exposes everything through a clean web UI. It solves a real problem: community organizations, small businesses, and academic groups need private, affordable, and accessible AI — and they need it on *their* hardware, under *their* control.

### Secondary: Ollama Special Technology Prize ($10K)

Direct hit. The entire inference stack runs on Ollama. Local E4B and E2B for privacy-first work, cloud 31B for complex queries. The architecture specifically showcases how Ollama enables a production-grade agent with multi-model routing and cloud fallback.

**Can win both Main + Ollama simultaneously.**

### Tertiary: LiteRT Special Technology Prize ($10K)

Future-path angle. The Pixel 9 Pro Fold running a fine-tuned Gemma 4 E2B model via LiteRT for on-device mobile inference. Not the primary submission, but the architecture explicitly plans for a mobile companion surface (Phase 6 in the roadmap). Demonstrating that the same Savitar gateway can route to an on-device LiteRT model on Android makes the technical depth story stronger even if the mobile path isn't fully built yet.

**Maximum possible: Main ($50K) + Ollama ($10K) + LiteRT ($10K) + Impact ($10K) = $80K**

### Impact Tracks (pick up to 2):

| Track | Prize | Angle |
|---|---|---|
| **Digital Equity & Inclusivity** | $10K | Local-first AI for underserved Chicago communities and rural Illinois. Small business owners who can't afford API subscriptions get a private agent on a Mac Mini. Multilingual support through Gemma 4's native language capabilities. |
| **Future of Education** | $10K | ISU partnership. Academic collaboration network with verified identity. Students and faculty use Savitar to access course knowledge, research resources, and community data through a shared, privacy-preserving agent. |
| **Safety & Trust** | $10K | Clean-room architecture. Honest disclosure policy. Operator guardrails. Shell policy enforcement. No deceptive identity claims. Audit logs. |

**Recommended primary impact angle: Digital Equity & Inclusivity** — it's the strongest story given the Chicago PoC small business connection and multiple budget tiers.

---

## Evaluation Alignment

### Impact & Vision (40 pts)

| Criterion | Our answer |
|---|---|
| Real-world problem | Small businesses, community orgs, and academic groups need AI they control, on hardware they own, at costs they can afford — without sending private data to cloud APIs |
| Target users | Chicago PoC small business owners, ISU students and faculty, community college learners, rural Illinois nonprofits |
| Scale potential | Savitar is open-source, runs on hardware from a $499 Apple Education Mac Mini to a $2,999 Mac Studio, and connects to messaging platforms people already use. Schools already get Apple Education discounts. |
| Gemma 4 usage | Core model for all inference — E2B for budget hardware and mobile (LiteRT on Pixel), E4B for standard (Ollama on Mac Mini), 31B-cloud for peak demand |
| Ecosystem breadth | Apple (Mac Mini + iMessage) and Google (Pixel Fold + Fi + Android + LiteRT) — covers both ecosystems from a single agent |

### Video Pitch & Storytelling (30 pts)

**3-minute video outline:**
1. (0:00–0:30) — The problem: meet Maria, a Chicago small business owner. She needs help with LLC filings, grant applications, and customer questions. She can't afford $200/month API bills. She speaks Spanish at home.
2. (0:30–1:15) — The solution: Savitar running on a $599 Mac Mini in Maria's back office. She messages Artemis on Discord. It answers in English and Spanish, pulls up grant deadlines, and helps draft applications — all local, all private.
3. (1:15–1:45) — The architecture: show the web UI, the multi-surface gateway (Discord + WhatsApp + iMessage), the Gemma 4 model routing, the Ollama inference. Show it working live.
4. (1:45–2:15) — Scale it: ISU students use the same system for research collaboration. A community college runs their own instance. Schools get the Mac Mini for $499 through Apple Education. Show the budget tiers from $0 to $3000. Quick flash of the Pixel Fold as the mobile operator companion.
5. (2:15–2:45) — The "for good": why local-first matters for communities that can't trust cloud providers with their data. Why identity-verified channels build real networks. Why open-source means anyone can run this.
6. (2:45–3:00) — Call to action and credits.

### Technical Depth & Execution (30 pts)

| Criterion | Our answer |
|---|---|
| Architecture | Clean Go runtime, contract-driven, 17 test packages, modular boundaries |
| Model routing | Four-lane policy: local-default (E4B), fast (E2B), standard (cloud 31B), high-capability (cloud 31B) |
| Ollama integration | Local inference on :11434, cloud fallback to api.ollama.com, health checks, per-model routing |
| Multi-surface | Discord shipped, WhatsApp and iMessage contracted, unified gateway envelope |
| Mobile companion | Pixel 9 Pro Fold on Google Fi as mobile operator surface. Future: on-device E2B via LiteRT for offline/private mobile inference |
| MCP orchestration | 5 servers configured (GitHub, Context7, Tavily, Playwright, Peekaboo) |
| Memory system | Filesystem-backed packs with provenance, subject tagging, temporal retrieval |
| Web UI | Planned: operator dashboard with Google login, session view, conversation monitor |
| Clean-room | Full provenance log, ADR charter, no copied code |

---

## Submission Deliverables Checklist

| Deliverable | Status | Notes |
|---|---|---|
| Kaggle Writeup (≤1500 words) | Not started | Draft from this strategy doc |
| YouTube Video (≤3 min) | Not started | Script outlined above |
| Public code repo | Ready | github.com/AlexiosBluffMara/Savitar |
| Live demo URL | Not started | Need web UI endpoint |
| Cover image | Not started | |
| Media gallery | Not started | Screenshots of Discord, web UI, terminal |

---

## Timeline to May 18

| Week | Focus |
|---|---|
| **Apr 7–13** | MCP runtime client (Milestone 1), memory packs for ISU/grant knowledge |
| **Apr 14–20** | Web UI backend + minimal frontend, WhatsApp bridge contract |
| **Apr 21–27** | iMessage bridge contract, budget-tier deployment docs, demo polish |
| **Apr 28–May 4** | Record video, write Kaggle writeup, create media gallery |
| **May 5–11** | Live demo hosted, submission draft review |
| **May 12–18** | Final submission, buffer for fixes |
