# Submission Strategy — Gemma 4 Good Hackathon

This document reflects the active hackathon strategy on the `gemma4-good-hackathon` branch.

## Core position

**Project name:** Savitar  
**Operator story:** Alexios Bluff Mara LLC runs a local-first Gemma 4 operator for Chicago small businesses, educators, and trusted community networks.  
**Primary surfaces:** Invite-only Discord for end users, operator web console for review and oversight.  
**Primary hardware:** Mac Mini M4 on bare metal with local Ollama, plus attached SSD for knowledge packs and run records.  
**Stretch device:** Pixel Fold as the operator companion and possible Live Mode surface.

Savitar should not be pitched as a generic multi-surface AI shell. For this competition, it should be pitched as a trustworthy local operator that helps people do real work without handing their conversations and local documents to a hosted model vendor.

## Official competition requirements

These points come directly from the official Kaggle overview and rules pages.

- One team can submit one writeup only.
- Maximum team size is five.
- The writeup must stay under 1,500 words.
- The video must be public on YouTube and 3 minutes or less.
- The code repository must be public.
- The live demo must be publicly accessible and not require a paywall.
- A media gallery and cover image are required.
- Winners must be prepared to provide reproducible code and documentation.
- External tools and data are allowed only when they are reasonably accessible at minimal cost.

## How judging really works

The public overview defines the scoring split:

- Impact & Vision: 40 points.
- Video Pitch & Storytelling: 30 points.
- Technical Depth & Execution: 30 points.

The host’s public welcome thread reinforces that the video is the judge’s primary lens. That means the project must be designed backward from the demo story, not forward from an architecture wish list.

## Award coverage strategy

### Main Track

Savitar is credible for the Main Track if the submission proves three things in one cohesive story:

1. A real problem exists for Chicago operators, educators, and community-serving small businesses.
2. Gemma 4 is central to the solution rather than bolted on.
3. The code, dashboard, and Discord flow are functional enough to survive scrutiny.

### Impact Track priorities

The highest-value impact lanes for the current project shape are:

| Track | Why Savitar fits |
|---|---|
| **Future of Education** | Educator Live Mode, role-gated classroom support, local course packs, and teacher-facing orchestration all map directly to the official education language around multi-tool agents that empower the educator. |
| **Digital Equity & Inclusivity** | Chicago small businesses and multilingual learners benefit from local deployment, low hardware cost, file-and-link workflows, and interface choices that do not assume enterprise tooling. |
| **Safety & Trust** | Invite-only Discord, operator review, provenance, explicit guardrails, evidence-backed answers, and local document control support a concrete trust story. |

Health & Sciences and Global Resilience can remain adjacent narratives, but they should not dilute the main pitch unless a real demo scenario is built.

### Special Technology Track priorities

| Prize | Status | Rationale |
|---|---|---|
| **Ollama** | Direct target | Savitar now runs Gemma 4 locally through Ollama on the Mac Mini. This is the cleanest special-tech fit in the repo today. |
| **Unsloth** | High-upside target | If a focused fine-tuned Gemma 4 adapter ships for educator or small-business workflows, Savitar becomes credible here too. |
| **Cactus** | Stretch | Possible only if the mobile companion truthfully routes tasks between surfaces or models in a live demo. |
| **LiteRT** | Stretch | Only pursue if an actual on-device Gemma flow exists. Do not pitch this as shipped before it exists. |
| **llama.cpp** | Not the main bet | The local Mac Mini plus Ollama story is stronger and more coherent for this repo. |

The rules explicitly state that a project may win both a Main Track prize and a Special Technology prize. Public discussion has not yet clarified whether one entry can win multiple prizes inside the same track, so Savitar should maximize category overlap without promising same-track stacking.

## Product story to tell

Savitar is best framed as a neighborhood-scale operator system, not a chatbot.

### Primary human story

Alexios Bluff Mara LLC operates a private Savitar instance for Chicago small-business owners and educators.

- A small-business owner uses Discord to ask about grants, forms, vendors, city programs, and operating checklists.
- An educator uses the same system in Live Mode to coordinate lesson materials, answer student questions, and surface curated resources during live instruction.
- The operator reviews outputs, curates local packs, and keeps sensitive material on local storage instead of hosted AI infrastructure.

### Why this story is stronger than the older one

- It avoids overcommitting to a specific university partnership before one exists.
- It stays closer to your real operator context and your LLC.
- It naturally covers education, digital equity, and trust in one product.
- It keeps the submission focused on what can actually be shown in a 3-minute video.

## System pillars judges should see

### 1. Local inference on hardware people can actually buy

- Gemma 4 on Ollama.
- Mac Mini M4 as the standard deployment story.
- Attached SSD for local knowledge packs, notes, course materials, and operating documents.

### 2. Invite-only Discord with real operational controls

- Role-based channels for small-business owners, educators, operators, and reviewers.
- File, link, and short video workflows instead of prompt-only interaction.
- Clear distinction between end-user channels and operator channels.

### 3. Evidence-backed answers instead of empty assistant theater

- Local markdown retrieval.
- Source metadata.
- Operator review queue.
- Honest fallback behavior when the model or evidence is insufficient.

### 4. Educator Live Mode

This should be the simplest educational story:

- one operator starts a live session,
- one subject pack is active,
- questions come in through a dedicated Discord channel,
- the dashboard shows what evidence was used,
- the educator can intervene, pin, or reject an answer.

That is much easier to demo than a broad “education platform,” and it directly serves the official Future of Education framing.

### 5. GitHub and Copilot as background leverage

Use GitHub and Copilot to strengthen the operator story rather than replace the local model story.

- GitHub stores the open code and issue flow.
- Copilot helps maintain documentation, research summaries, and implementation throughput.
- The local Gemma lane remains the actual submission centerpiece for end-user interaction.

## Demo narrative for the 3-minute video

### Segment 1: the problem

Chicago operators, teachers, and small-business owners need an AI system that can work with local documents, invite-only communities, and real oversight. They cannot rely on a black-box hosted assistant for every question.

### Segment 2: the system

Show Savitar on a Mac Mini running Gemma 4 locally through Ollama, with Discord as the working user surface and the web console as the operator surface.

### Segment 3: two concrete workflows

- Small-business workflow: grant or program discovery with local evidence.
- Educator workflow: Live Mode with local class materials and operator review.

### Segment 4: why this matters

- Lower hardware cost.
- Local control.
- Better trust posture.
- Useful in real community settings.

### Segment 5: technical credibility

Show the repo, the dashboard, the local model route, and the evidence trail.

## Deliverables checklist

| Deliverable | Requirement | Savitar guidance |
|---|---|---|
| Writeup | Public, under 1,500 words | Focus on one coherent story and one truthful architecture |
| Video | Public YouTube, 3 minutes or less | Prioritize real workflows and human impact over feature inventory |
| Code | Public repository | Keep the local-Ollama and Discord flow easy to inspect |
| Live demo | Publicly accessible | Expose a truthful demo surface, likely read-only for judges |
| Media gallery | Required | Screenshots of Discord, dashboard, evidence view, and local deployment |

## Immediate build priorities

1. Keep stripping cloud-only assumptions from the repo and docs.
2. Strengthen the local Ollama route and prove it in the dashboard and Discord flows.
3. Refine the operator dashboard around evidence, review, and Live Mode.
4. Prepare a single disciplined Unsloth fine-tune plan instead of vague “later” language.
5. Tighten the public demo story so it matches the official scoring weights.