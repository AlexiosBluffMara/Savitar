# Deployment Tiers — Hardware, Privacy, Budget, and Audience

Savitar is designed to run at every level of investment and technical ability. This document outlines the concrete deployment tiers, what each costs, who it serves, and what tradeoffs it makes.

---

## The Unifying Project: One Agent, Every Surface, Every Budget

The single idea that ties everything together: **one AI agent that a person, a school, a small business, or a community org can deploy on the hardware they already own — and reach everyone they serve through the messaging platforms those people already use.**

This isn't about building separate apps for each platform. It's about building one system that spans:

- **Apple ecosystem** — Mac Mini (server), iMessage (transport), Apple Education pricing for schools
- **Google ecosystem** — Pixel Fold (mobile operator), Google Fi (connectivity), Android companion app (future LiteRT on-device inference), Google login for web UI
- **Open ecosystem** — Discord (community), WhatsApp (international), Ollama (local inference), any Linux box

The Mac Mini is the server. The Pixel Fold is the operator's mobile dashboard. Discord is where the community lives. iMessage is how Americans text. WhatsApp is how the rest of the world texts. The web UI is how judges see it all come together.

Every tier below uses the same Savitar runtime, the same Gemma 4 models, the same gateway, and the same memory packs. The only variable is hardware cost and which transports you enable.

---

## Tier 0: Free / Zero Hardware — Cloud-Only Student Setup

**Budget:** $0  
**Hardware:** Any laptop or Chromebook with a browser  
**Audience:** Students, learners, people trying it out  
**Tech-savviness:** Low — just needs Discord and a browser

### What You Get
- Join a community Discord server that already runs Savitar
- Chat with the agent through `@mention` or DMs
- No installation, no setup, no hardware purchase

### Model Routing
- All inference handled by the community operator's hardware
- No local model — you're a user, not an operator

### Privacy
- Messages pass through Discord's servers
- Agent replies generated on the operator's local machine
- Your conversation history stays with the operator's instance (not in the cloud)

### Limitations
- Dependent on someone else running the server
- No customization, no private knowledge

---

## Tier 1: Budget Local — Raspberry Pi 5 or Old Mac

**Budget:** $100–$300  
**Hardware:** Raspberry Pi 5 (8GB), old MacBook Air, or any Linux box with 8GB RAM  
**Audience:** Hobbyists, students with old hardware, small nonprofits  
**Tech-savviness:** Medium — comfortable with terminal

### What You Get
- Savitar CLI running locally
- Gemma 4 E2B (2.3B effective, ~4GB model file) via Ollama
- Discord bot in a single server
- Basic memory packs loaded from disk

### Model Routing
- Local-only: `gemma4:e2b` on Ollama
- No cloud fallback (to stay within hardware limits)
- 128K context window, but keepalive and generation will be slow

### Privacy
- Full local inference — nothing leaves your network
- All data on your hardware
- No cloud API keys needed

### Limitations
- Slow generation (~5–10 tokens/sec on Pi 5)
- E2B is capable but less nuanced than E4B
- Single transport only (Discord or CLI)

---

## Tier 2: Standard Local — Mac Mini M4 (Recommended)

**Budget:** $499–$799 ($499 with Apple Education pricing)  
**Hardware:** Mac Mini M4 with 16GB unified memory  
**Audience:** Small businesses, community orgs, academic departments, LLC operators  
**Tech-savviness:** Medium — follow setup guide, run shell scripts

### Apple Education Pricing

Schools, colleges, and universities get the Mac Mini M4 for **$499** through Apple Education (≈$41.58/mo for 12 months). This is available year-round to K-12 institutions, community colleges, and higher education. Students and faculty also qualify for personal education pricing. The seasonal Back to School promotion adds free accessories (AirPods, peripherals) on top of the education discount.

This means a community college IT department can deploy a Savitar instance for their ESL program for under $500 — less than 3 months of most SaaS AI subscriptions.

### What You Get
- Full Savitar runtime with all CLI commands
- Gemma 4 E4B (4.5B effective, 11GB Q8 quantization) as primary model
- Gemma 4 E2B as fast fallback for high-load periods
- Discord + iMessage transports (iMessage native on macOS)
- Web UI for operator dashboard
- Memory packs for subject-area knowledge
- MCP tool integration (GitHub, Tavily, Context7)

### Model Routing
- Primary: `gemma4:e4b-it-q8_0` (local)
- Fast lane: `gemma4:e2b` (local)
- Cloud fallback: `gemma4:31b-cloud` via Ollama Cloud (opt-in)
- Policy: local-first, cloud only for complex queries with explicit operator consent

### Privacy
- All default inference is local — zero cloud calls unless opted in
- iMessage integration is native macOS — no bridge server needed
- Discord messages processed locally, replies generated on your hardware
- Cloud opt-in is per-conversation, not global

### Why This Is the Sweet Spot
- $499–$599 gets you Apple Silicon with 16GB unified memory ($499 edu, $599 retail)
- E4B Q8 fits comfortably with room for the OS and tools
- iMessage is free and native — no WhatsApp Business API cost
- Runs 24/7 as a headless server (launchd service)
- This is what the live demo runs on
- Apple Education pricing means schools can deploy for under $500

---

## Tier 3: Enhanced Local — Mac Mini M4 Pro or Mac Studio

**Budget:** $1,399–$2,999  
**Hardware:** Mac Mini M4 Pro (24GB+) or Mac Studio M4 Max (32GB+)  
**Audience:** Larger organizations, multi-department academic setups, high-traffic community servers  
**Tech-savviness:** Medium-High — manages multiple services

### What You Get
- Everything in Tier 2, plus:
- Gemma 4 E4B at higher quantization (Q8 or full precision)
- Support for longer context and larger memory packs
- Multiple simultaneous conversations with fast generation
- All three transports: Discord + iMessage + WhatsApp
- Browser automation via Playwright MCP
- Repository analysis pipeline

### Model Routing
- Primary: `gemma4:e4b` at full quality
- Cloud fallback: `gemma4:31b-cloud` for complex multi-step tasks
- Concurrent request handling for busy servers

### Privacy
- Same local-first model as Tier 2
- More headroom for additional tools and integrations without sacrificing performance

---

## Tier 4: Cloud-Assisted — VPS or Cloud VM

**Budget:** $20–$100/month  
**Hardware:** Cloud VM with GPU (Lambda, RunPod, Vast.ai) or CPU-only VPS  
**Audience:** Organizations that want always-on access from anywhere, teams with no local hardware, international deployment  
**Tech-savviness:** High — manages cloud infrastructure

### What You Get
- Savitar deployed to a cloud VM
- Ollama running in the cloud (GPU VM) or Ollama Cloud API
- All three transports
- Web UI accessible from anywhere
- Automated deployment via scripts

### Model Routing
- GPU VM: local Ollama with `gemma4:e4b` or even larger models
- CPU-only: `gemma4:e2b` with cloud fallback to `31b-cloud`
- Ollama Cloud only: `gemma4:31b-cloud` for all inference ($, but no hardware to manage)

### Privacy
- Data is on the cloud provider's infrastructure
- Encrypted in transit, but the VM operator has access
- Suitable for less sensitive use cases or where availability matters more than data residency

### When This Makes Sense
- Team members in different cities or countries
- WhatsApp Business API requires a public webhook endpoint
- Need 24/7 uptime without managing physical hardware

---

## Tier 5: Mobile Companion — Google Pixel Fold on Google Fi

**Budget:** $1,799 (Pixel 9 Pro Fold) + $20–$65/mo (Google Fi)  
**Hardware:** Google Pixel 9 Pro Fold, Tensor G4 processor, 16GB RAM  
**Audience:** Operators who need mobile access, demo presenters, field workers  
**Tech-savviness:** Low for using, High for on-device fine-tuning

### What You Get
- **Mobile operator dashboard** — manage Savitar from anywhere via the web UI in the Fold's 8" inner display
- **Google Fi connectivity** — reliable data for remote management, works internationally
- **Gemini built-in** — Pixel 9 Pro Fold ships with Gemini Advanced (1 year free), complementary to Savitar's Gemma 4 stack
- **Future: on-device Gemma 4 E2B** — fine-tuned E2B model running via LiteRT directly on the Tensor G4 for offline/private mobile inference

### What This Tier Is (and Isn't)

This tier is **not** about running the full Savitar server on a phone. The Mac Mini stays home as the inference server. The Pixel Fold is:

1. **The operator's mobile window** — check conversations, manage memory packs, monitor model routing from your pocket
2. **A demo device** — walk into a hackathon demo or a community meeting and show the web UI on the Fold's big screen
3. **A future inference endpoint** — Gemma 4 E2B via LiteRT on Tensor G4 for completely offline, completely private mobile AI
4. **A WhatsApp endpoint** — Google Fi number for the WhatsApp Business bridge

### Google Fi as Infrastructure

Google Fi is relevant because:
- It's a Google-native carrier — tight integration with Pixel hardware
- International coverage for WhatsApp-heavy use cases
- Low-cost plans ($20/mo basic) for always-on connectivity
- The Fi number can serve as the WhatsApp Business number for the agent

### On-Device Fine-Tuning (Future Path)

Fine-tuning Gemma 4 E2B for mobile tool usage via LiteRT is technically possible but complex. This is explicitly a Phase 6+ goal. For the hackathon, the architecture *documents* this path and shows the gateway contract that would route to an on-device model, even though the primary demo runs on the Mac Mini. Judges see that the system is designed for it, even if it's not built yet.

### Why This Matters for the Hackathon
- Demonstrates both Apple and Google ecosystem coverage
- Shows the Pixel 9 Pro Fold / Gemma 4 / Google Fi stack as a cohesive Google hardware story
- Qualifies for the **LiteRT Special Technology Prize** ($10K) if on-device inference is demoed
- The Fold's 8" inner screen is perfect for demo presentations

---

## Tier Comparison Matrix

| | Tier 0 | Tier 1 | Tier 2 | Tier 3 | Tier 4 | Tier 5 |
|---|---|---|---|---|---|---|
| **Cost** | $0 | $100–300 | $499–799 | $1,399–2,999 | $20–100/mo | $1,799 + $20/mo |
| **Model** | (shared) | E2B | E4B Q8 | E4B full | E4B or Cloud | Cloud (future: E2B LiteRT) |
| **Transports** | Discord | Discord | Discord + iMessage | All three | All three | Mobile dashboard + WhatsApp |
| **Privacy** | Low | High | High | High | Medium | Medium (future: High w/ LiteRT) |
| **Setup difficulty** | None | Medium | Medium | Medium-High | High | Low (user) / High (dev) |
| **Target user** | Students | Hobbyists | Small biz / academic | Departments | Distributed teams | Mobile operators |
| **Best for** | Learning | Experimenting | **Production use** | Scale | Remote access | Field work & demos |

---

## Budget Guidance for Specific Audiences

### Chicago PoC Small Business Owner
**Recommended: Tier 2 ($499–$599)**  
A Mac Mini in your back office. Savitar helps with LLC filings, grant applications, customer FAQ, and inventory questions. Talk to it on iMessage from your phone. All data stays in your shop. Write off the hardware as a business expense.

### ISU Academic Department
**Recommended: Tier 2 or Tier 3 ($499–$1,399)**  
One Mac Mini per department or shared lab — **$499 with Apple Education pricing**. Students interact via Discord. Faculty manage via web UI. Memory packs loaded with course materials, research resources, and community partner information. ISU's Center for Civic Engagement can use it for reporting and coordination. Faculty and grad students with Pixel Folds can monitor from the field during community engagement activities.

### Community College ESL Program
**Recommended: Tier 0 (free) or Tier 1 ($100–$300)**  
Students join a Discord server. The instructor (or a grant-funded setup) runs the agent on a Raspberry Pi or donated Mac. Gemma 4's multilingual capabilities handle Spanish, Arabic, and other languages for immigrant learners.

### Rural Illinois Nonprofit
**Recommended: Tier 1 ($100–$300)**  
An old laptop, Ollama, and E2B. Runs in the office during business hours. Handles common questions about services, hours, and resources. Slow but private and free to operate.

### Illinois LLC / Startup
**Recommended: Tier 2 ($499–$599)**  
Same as small business. Use Savitar to draft grant applications (SBIR/STTR, Neighborhood Opportunity Fund, DCEO programs), track deadlines, and manage team communication across Discord and iMessage. The agent remembers your filing dates, your FEIN, and your grant history.

### Mobile Operator / Demo Presenter
**Recommended: Tier 5 ($1,799 + Fi)**  
You already have the Mac Mini running at home. The Pixel Fold lets you manage everything from the field. Walk into a hackathon demo, a community meeting, or a client pitch and show the full dashboard on an 8" folding screen. Google Fi keeps you connected anywhere. The Fi number doubles as your WhatsApp Business number for international contacts.

---

## Model Performance: Google Default vs Unsloth Quantized

Savitar's value proposition depends on running Gemma 4 **locally**. That means choosing a quantization level — a tradeoff between model size, memory footprint, and output quality. Two sources dominate:

1. **Google official weights** — full-precision BF16 and ggml-org Q8 GGUF files. Maximum quality. Large.
2. **Unsloth quantized** — community-optimized GGUF files at Q4_K_M, Q4_K_XL, Q2_K, and other levels. Smaller, faster, some quality loss.

### What Is Quantization?

Quantization reduces the numerical precision of model weights. Instead of storing each weight as a 16-bit float (BF16), you compress it to 8 bits (Q8), 4 bits (Q4), or even 2 bits (Q2). The result is a smaller file that uses less RAM and runs faster — but with reduced accuracy at each step down.

Think of it like JPEG compression for images. A high-quality JPEG at 95% looks indistinguishable from the original. At 50% you start to see artifacts. At 10% it's blurry but you can still tell what the picture is. The same principle applies to model quantization.

### Gemma 4 Model Sizes: Official vs Unsloth Quantized

#### E2B (2.3B effective parameters — edge/phone tier)

| Quantization | Provider | File Size | RAM Needed | Notes |
|---|---|---|---|---|
| BF16 (full) | Google / ggml-org | 9.3 GB | ~10 GB | Maximum quality, reference baseline |
| Q8_0 | ggml-org | 5.0 GB | 5–8 GB | Near-lossless, recommended minimum |
| Q8_0 | Unsloth | 5.1 GB | 5–8 GB | Functionally identical to ggml-org |
| Q4_K_M | Unsloth | 3.1 GB | ~4 GB | **Sweet spot for constrained hardware** |
| Q4_K_XL | Unsloth | 3.2 GB | ~4 GB | Slightly better than Q4_K_M |
| Q2_K | Unsloth | ~2.0 GB | ~3 GB | Noticeably degraded; last resort |

#### E4B (4.5B effective parameters — our Mac Mini primary)

| Quantization | Provider | File Size | RAM Needed | Notes |
|---|---|---|---|---|
| BF16 (full) | Google / ggml-org | ~16 GB | ~16 GB | Won't leave room for OS on 16GB machine |
| Q8_0 | ggml-org / Ollama | ~9.6 GB | 9–12 GB | **Our default**: high quality, fits on 16GB Mac Mini |
| Q4_K_M | Unsloth | ~5.5 GB | 5.5–6 GB | Good quality, leaves tons of headroom |
| Q4_K_XL | Unsloth | ~5.7 GB | ~6 GB | Minor improvement over Q4_K_M |
| Q2_K | Unsloth | ~3.5 GB | ~4 GB | Significant quality loss; only for extreme constraints |

#### 26B A4B MoE (3.8B active parameters)

| Quantization | Provider | File Size | RAM Needed | Notes |
|---|---|---|---|---|
| BF16 (full) | Google | ~52 GB | ~52 GB | Needs workstation-class hardware |
| Q8_0 | ggml-org | ~28 GB | 28–30 GB | Fits on 32GB Mac Mini M4 Pro |
| Q4_K_M | Unsloth | ~16 GB | 16–18 GB | **Just barely fits on 16GB Mac Mini** (tight) |

### How Much Quality Do You Actually Lose?

This is the key question. Google publishes benchmarks for the full-precision models. Quantization degrades these scores by varying amounts depending on the task type. Based on published community benchmarks and Unsloth's own testing:

**Estimated performance retention by quantization level:**

| Quantization | Estimated Quality Retention | What That Means |
|---|---|---|
| BF16 (full) | 100% (baseline) | Google's published benchmark scores |
| Q8_0 | ~98–99% | Near-lossless. You cannot tell the difference in conversation |
| Q4_K_M | ~92–95% | Slightly less nuanced on complex reasoning. Fine for conversation, Q&A, summaries |
| Q4_K_XL | ~93–96% | Unsloth's enhanced Q4. Marginally better than Q4_K_M |
| Q2_K | ~80–85% | Noticeable degradation. More repetition, less coherent long-form. Last resort |

**Concrete example with E4B benchmark scores (MMLU Pro):**

| Quantization | Estimated MMLU Pro | LiveCodeBench v6 | GPQA Diamond |
|---|---|---|---|
| BF16 (baseline) | 69.4% | 52.0% | 58.6% |
| Q8_0 | ~68–69% | ~51–52% | ~57–58% |
| Q4_K_M | ~64–66% | ~48–49% | ~54–56% |
| Q2_K | ~56–59% | ~42–44% | ~47–50% |

**The honest takeaway:** Q8 is essentially free in quality terms. Q4 costs you 5–8% on benchmarks but saves nearly half the RAM. Q2 costs you 15–20% and should only be used when RAM is the hard constraint.

### Dollars to Output Quality: The Performance Chart

This is what actually matters for a teacher or small business deciding what to buy.

```
DOLLARS-TO-QUALITY CHART: Gemma 4 on Local Hardware
====================================================

Hardware Cost        Model Config              RAM Used  Quality*  tok/s**  Monthly Cost
─────────────────────────────────────────────────────────────────────────────────────────
$0 (shared)          (someone else's E4B)       0 GB     95%       N/A      $0
$100 (Pi 5 8GB)      E2B Q4_K_M (Unsloth)       4 GB     55%      ~5-8     $0 (electricity)
$150 (old MacBook)   E2B Q8_0                    5 GB     59%      ~8-12    $0
$200 (Pi 5 + SSD)    E4B Q4_K_M (Unsloth)       6 GB     65%      ~3-5     $0
$499 (Mac Mini edu)  E4B Q8_0 (Ollama)          11 GB     68%      ~25-35   $0
$599 (Mac Mini)      E4B Q8_0 (Ollama)          11 GB     68%      ~25-35   $0
$599 (Mac Mini)      E4B Q4_K_M (Unsloth)        6 GB     65%      ~40-50   $0
$799 (Mac Mini 24G)  E4B BF16 (full)            16 GB     69%      ~20-25   $0
$1,399 (M4 Pro 24G)  26B A4B Q4_K_M (Unsloth)  18 GB     80%      ~15-20   $0
$1,999 (M4 Pro 36G)  26B A4B Q8_0              30 GB     82%      ~20-30   $0
$20/mo (cloud CPU)   E2B Q4 on VPS               4 GB     55%      ~3-5     $20/mo
$50/mo (cloud GPU)   E4B Q8 on GPU VM           12 GB     68%      ~40-60   $50/mo
$100/mo (Ollama Cloud) 31B cloud API              0 GB     85%      ~30-50   ~$100/mo

* Quality = estimated MMLU Pro score (benchmark proxy). Higher is better.
** tok/s = tokens per second generation speed on Mac Mini M4 16GB unless noted.
   Real-world perception: >20 tok/s feels instant; 10-20 is comfortable; <10 is noticeably slow.
```

**Key insight:** The **$499–$599 Mac Mini with E4B Q8** is the best value in the entire chart. You get 98% of full-precision quality, 25–35 tok/s (feels instant), zero monthly cost, and complete privacy. Going cheaper (Pi 5 + E2B Q4) saves $300–$400 but costs you roughly 13 quality points. Going more expensive (M4 Pro + 26B) gains you ~14 quality points for $800–$1,400 more.

The Unsloth Q4 quantizations are most valuable at **Tier 1** (budget hardware): they make the difference between "can't run E4B at all" and "can run E4B Q4 with room to spare." At Tier 2 (Mac Mini 16GB), you have enough RAM for Q8, so the Unsloth Q4 savings are nice-to-have (faster generation, more RAM for tools) but not critical.

### When to Use Unsloth Quantized vs Google Default

| Scenario | Recommendation | Why |
|---|---|---|
| Mac Mini M4 16GB, primary model | **Google/Ollama Q8** | Best quality that fits; RAM is sufficient |
| Mac Mini M4 16GB, need to run E4B + tools + large context | **Unsloth Q4_K_M** | Frees ~5GB for MCP tools and context |
| Raspberry Pi 5 8GB | **Unsloth Q4_K_M (E2B)** | Only way to fit a useful model |
| Old laptop 8GB | **Unsloth Q4_K_M (E2B)** | Same — Q8 E2B technically fits but leaves no room |
| Mac Mini M4 Pro 24GB | **Google/Ollama Q8 (E4B)** | Plenty of room; maximize quality |
| Mac Mini M4 Pro 36GB | **Google Q8 (26B A4B)** | MoE at full Q8; major quality jump |
| Phone / Pixel Fold (future) | **Unsloth Q4 or Q2 (E2B)** | LiteRT requires smallest possible artifact |

---

## Teacher Knowledge Packs: A Grounded Source of Truth

This is how Savitar turns from a generic chatbot into a **subject-specific teaching assistant** that any teacher can customize.

### The Problem

A bare language model knows a lot about a lot, but:
- It doesn't know *your* curriculum
- It doesn't know *your* school's specific policies
- It doesn't know the grant programs in *your* state
- It will confidently hallucinate answers about your course materials
- Different teachers need different knowledge for different classes

### The Solution: Memory Packs as a Local Knowledge Database

Savitar already has a memory pack system (code at `internal/memory/store.go` and `internal/memory/pack.go`). Each pack is:

```
~/.savitar/memory/<subject>/<name>.md
```

A pack is a simple Markdown file with frontmatter:

```markdown
---
name: us-history-unit-3
subject: history-101
source: teacher-authored
created: 2026-04-07T00:00:00Z
updated: 2026-04-07T00:00:00Z
---

## Unit 3: The Civil War and Reconstruction

### Key Dates
- 1861: Fort Sumter
- 1863: Emancipation Proclamation, Gettysburg
- 1865: Appomattox, Lincoln assassination
- 1865–1877: Reconstruction era

### Discussion Questions
1. How did the economic systems of North and South make conflict more likely?
2. What was the significance of the Emancipation Proclamation as a war measure vs a moral statement?

### Grading Rubric
- Thesis statement: 20 points
- Evidence from primary sources: 30 points
- Analysis and argumentation: 30 points
- Writing mechanics: 20 points

### Approved Sources
- Textbook: "The American Pageant" chapters 19–22
- Primary source collection: Lincoln-Douglas Debates
- Documentary: Ken Burns "The Civil War" episodes 1–5
```

### How It Works as a Grounded Source of Truth

When a student asks Savitar a question, the engine:

1. **Identifies the subject** from the conversation context (which Discord channel, which memory pack subject is loaded)
2. **Retrieves relevant packs** from the filesystem — no network call, no API, no cloud
3. **Injects pack content as context** before the model generates a response
4. **The model answers grounded in the pack** — it cites the teacher's rubric, the approved sources, the specific dates and facts the teacher loaded

This is **retrieval-augmented generation (RAG)** but with a critical difference: the retrieval database is a folder of Markdown files the teacher writes and controls. No vector database. No embedding pipeline. No cloud service. Just files on a Mac Mini.

### What a Teacher Can Do

| Action | How | Time Required |
|---|---|---|
| Create a new subject knowledge base | Create a folder, write Markdown files | 30 minutes for a unit |
| Add approved sources | List them in a pack file | 5 minutes |
| Add a grading rubric | Write it as a section in a pack | 10 minutes |
| Add FAQ answers | Write them as pack entries | 15 minutes |
| Update for next semester | Edit the Markdown files | 20 minutes |
| Share with another teacher | Copy the folder | 1 minute |
| Remove a pack | Delete the file | Instant |

### Teacher Knowledge Pack Examples

**ESL Program (Community College):**
```
~/.savitar/memory/esl-101/
├── vocabulary-unit-1.md        # Word lists, definitions, example sentences
├── grammar-present-tense.md    # Rules, exceptions, practice exercises
├── cultural-context-us.md      # US customs for new arrivals
├── local-resources-chicago.md  # Food banks, legal aid, healthcare in Spanish
└── approved-translations.md    # Teacher-verified translations for key terms
```

**Small Business Owner (Chicago PoC):**
```
~/.savitar/memory/my-llc/
├── filing-dates.md             # IL annual report due date, tax deadlines
├── grant-programs.md           # NOF, MBDA, SBIR eligibility
├── customer-faq.md             # Hours, services, pricing
├── vendor-contacts.md          # Supplier list, terms, reorder thresholds
└── neighborhood-resources.md   # Local chamber, SCORE mentor contacts
```

**ISU Academic Department:**
```
~/.savitar/memory/civic-engagement/
├── partner-orgs.md             # Community partner directory
├── reporting-template.md       # ISU CCE required fields
├── student-hours-policy.md     # Service learning credit requirements
├── irb-guidelines.md           # Research ethics quick reference
└── grant-deadlines.md          # Active grant timelines
```

### Why This Matters for the Hackathon

1. **No fine-tuning required** — the knowledge is injected as context, not baked into model weights
2. **Any teacher can do it** — if you can write a document, you can create a knowledge pack
3. **Instant updates** — edit a file, the next conversation uses the new version
4. **Sharable** — zip a folder, email it to another teacher running Savitar
5. **Verifiable** — the teacher can read exactly what the model has access to
6. **Private** — the knowledge never leaves the Mac Mini
7. **Composable** — load packs from multiple subjects for interdisciplinary work

---

## Privacy and Intellectual Property: What Running Locally Actually Means

### The Core Claim

When you run Gemma 4 on a Mac Mini via Ollama, your prompts, your knowledge packs, your conversations, and your students' questions **never leave your network**. No API call. No telemetry. No training data contribution. The model weights are Apache 2.0 licensed and the inference is entirely local.

This is fundamentally different from using ChatGPT, Claude, or Gemini through their web APIs, where your inputs are transmitted to a remote server (and may or may not be used for model training, depending on the provider's terms).

### Can You Own Your Prompts? The Legal Reality

The user asked: can teachers retain IP rights over their specific prompts?

**Short answer: Yes, with important nuances.**

Here's the current legal landscape as of April 2026:

#### What You CAN Protect

1. **Prompt libraries as trade secrets.** Under US trade secret law (Defend Trade Secrets Act), a curated library of prompts that you keep confidential, restrict access to, and derive commercial value from qualifies as a trade secret. This is the strongest protection available. As the law firm Klemchuk LLP notes: "Trade secret protection appears to offer the strongest legal framework for prompt libraries. Businesses should restrict access, use NDAs, and label prompt repositories as confidential."

2. **Prompts as copyrightable works (if sufficiently creative).** Simple prompts ("summarize this text") are not copyrightable — they're too utilitarian. But a carefully crafted prompt library with specific instructions, persona definitions, grading rubrics, and pedagogical structure likely meets the creativity threshold for copyright. The prompts themselves — not the AI's outputs — are your creative work.

3. **Knowledge packs as authored works.** The Markdown files a teacher writes are unambiguously copyrightable. They're original authored documents. If you write a unit guide, a grading rubric, approved source lists, and custom FAQ answers, those are your works. Period.

4. **Contractual protection.** Independent of copyright, you can use contracts (employment agreements, NDAs, licenses) to control how your prompts and knowledge packs are used.

#### What You CANNOT Protect (or Is Unsettled)

1. **AI-generated outputs are not copyrightable by the prompt author.** The US Copyright Office (January 2025 report) and the DC Circuit (*Thaler v. Perlmutter*, affirmed March 2025) are clear: copyright requires human authorship. "Prompts essentially function as instructions that convey unprotectable ideas" and the AI's output is not authored by the prompt writer. This means: your prompts are protectable, but the text the model generates from them probably isn't.

2. **Simple/functional prompts are not copyrightable.** "Explain photosynthesis to a 5th grader" is not creative enough for copyright. "You are a patient, encouraging tutor who breaks down cellular biology using sports metaphors, always checks for understanding before moving on, and flags when a student's question goes beyond the approved Unit 3 scope" probably is.

3. **The work-for-hire question.** If a teacher creates prompts as part of their job, the employer (school district, university) may own the copyright under work-for-hire doctrine. Teachers should check their employment agreements. This is the same issue as "who owns the lesson plans?" — it varies by state and contract.

### Why Running Locally Strengthens IP Protection

Here's where the Mac Mini deployment model becomes a **legal advantage**, not just a privacy feature:

| Protection | Cloud API | Local (Mac Mini + Ollama) |
|---|---|---|
| **Prompts sent to third party?** | Yes — sent to OpenAI/Google/Anthropic servers | No — never leaves your machine |
| **Prompts used for model training?** | Maybe — depends on ToS (opt-out available for some) | No — impossible, model is static local weights |
| **Trade secret viability** | Weakened — you've disclosed to a third party | Strong — you control all access |
| **Prompt confidentiality** | Requires trusting provider's data handling | Complete — you own the hardware |
| **Knowledge pack confidentiality** | Sent as context to remote API = disclosed | Stays on disk, never transmitted |
| **Contractual leverage** | Provider ToS may claim license to inputs | No third-party ToS applies |
| **Risk of subpoena/discovery** | Provider's servers may be subject to legal process | Only your hardware is at issue |
| **Student data (FERPA)** | Sending to cloud API may violate FERPA | Local processing is FERPA-compatible |

**The key legal point:** trade secret protection requires that you take "reasonable measures" to keep the secret. Running inference locally, storing knowledge packs on encrypted local storage, and never transmitting prompts to third parties meets this standard easily. Sending the same prompts to a cloud API undermines the "reasonable measures" argument.

### FERPA and Student Privacy

For schools specifically: the Family Educational Rights and Privacy Act (FERPA) restricts disclosure of student educational records to third parties. If a student asks Savitar about their grades, their accommodation status, or their academic standing, that information is an educational record. Sending it to a cloud API potentially constitutes a FERPA disclosure that requires parental consent.

Running Savitar locally eliminates this issue entirely. The student's question and the answer never leave the school's hardware.

### Practical Recommendations for Teachers

1. **Keep your knowledge packs in version control** (git) with access restricted to authorized staff
2. **Label your prompt libraries as "Confidential — Trade Secret"** in the file headers
3. **Check your employment agreement** for IP assignment clauses — especially for lesson plans and curriculum materials
4. **Do not upload your knowledge packs to cloud AI services** if you want to maintain trade secret protection
5. **Document your creative process** — if you ever need to assert copyright, showing the iterative development of your prompts strengthens the claim
6. **For FERPA compliance**, keep student-identifiable data in memory packs that never leave the local machine
7. **Consider licensing** if you want to share: Creative Commons for open sharing, or a custom license that allows educational use but prohibits commercial redistribution
