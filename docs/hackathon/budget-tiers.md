# Deployment And Work Tiers

Savitar should now be described through local deployment tiers and separate background-work tiers. The local Gemma 4 story stays constant: inference runs on hardware you control through Ollama.

## Local deployment tiers

### Tier 0: Guest Access

- Budget: $0 for the user.
- Hardware: none beyond a phone or browser.
- Audience: students, small-business clients, community members.
- Experience: join an invite-only Discord server and interact with Savitar through approved channels.
- Tradeoff: this is access to someone else’s local deployment, not ownership.

### Tier 1: Solo Builder

- Budget: $100 to $300.
- Hardware: old Mac, Linux desktop, or other spare box that can host a small local Ollama lane.
- Model lane: `gemma4:e2b`.
- Best for: experimentation, personal knowledge packs, light classroom prep, and low-volume community pilots.

### Tier 2: Standard Operator Host

- Budget: $499 to $799.
- Hardware: Mac Mini M4 with 16 GB memory.
- Model lane: `gemma4:e4b-it-q8_0` primary, `gemma4:e2b` secondary.
- Best for: Alexios Bluff Mara LLC, Chicago small businesses, school pilots, and community organizations.
- Why it matters: this is the strongest Ollama prize story because it is local, affordable, and believable.

### Tier 3: Community Hub

- Budget: $1,200 to $2,500.
- Hardware: higher-memory Mac Mini, Mac Studio, or another always-on local host with attached SSD.
- Model lane: same Gemma 4 local story, more headroom for concurrent users, larger packs, and longer sessions.
- Best for: educator teams, community labs, neighborhood business associations, and operator-reviewed shared deployments.

### Tier 4: Educator Live Mode

- Budget: Tier 2 or Tier 3 plus the operator’s existing laptop or foldable phone.
- Hardware: same local host as above, plus a portable dashboard surface.
- Experience: one educator runs a live, moderated question stream against curated class materials and local knowledge packs.
- Best for: the Future of Education pitch.

## Background-work tiers

These are not inference tiers. They describe how much background leverage the operator uses for implementation, maintenance, and knowledge curation.

### Tier A: Manual Operator

- Everything is curated by hand.
- Best for the smallest private pilots.

### Tier B: Copilot Routine Support

- GitHub Copilot helps with documentation, issue triage, implementation speed, and operator workflow polish.
- Best for keeping the repo moving without changing the local inference story.

### Tier C: Copilot Deep Support

- Higher-complexity Copilot lanes help with research synthesis, code generation, and toolchain work.
- Best for moving faster on integrations, tests, and public demo polish.

### Tier D: Structured Open Development

- GitHub issues, PR review, provenance, and Copilot-assisted iteration form the collaboration backbone.
- Best for proving this is a maintainable open project, not a weekend prototype.

## What to emphasize in the writeup

- The end-user AI story stays local through Ollama.
- The Mac Mini is the believable standard host.
- The attached SSD makes local packs and retrieval credible.
- Discord is the trust and access surface.
- Copilot and GitHub amplify operator productivity but do not replace Gemma 4 in the product story.

## What to avoid

- Do not reintroduce hosted-model fallback language.
- Do not describe “all categories” as guaranteed wins.
- Do not over-index on mobile or edge prizes unless a real demo exists.
- Do not drift back into a generic institution-first story if the stronger operator story is Chicago plus educators.