# Live Demo Plan

This file captures the broader demo concept drafted before the rewrite reset.

The current target demo should be read through [docs/adr/0006-rewrite-product-direction.md](../adr/0006-rewrite-product-direction.md) and [docs/roadmap/0006-rewrite-blueprint.md](../roadmap/0006-rewrite-blueprint.md): one working Discord surface, a truthful knowledge workflow, and an authenticated operator console before later transport expansion.

## Architecture Overview

The live demo runs at a publicly accessible URL and shows judges a real, working Savitar instance connected to live messaging surfaces.

```
┌──────────────┐   ┌──────────────┐   ┌──────────────┐
│   Discord     │   │   iMessage    │   │   WhatsApp    │
│   (live)      │   │   (planned)   │   │   (planned)   │
└──────┬───────┘   └──────┬───────┘   └──────┬───────┘
       │                  │                  │
       ▼                  ▼                  ▼
┌─────────────────────────────────────────────────────┐
│                  Savitar Gateway                     │
│          Normalized Message Envelope                 │
└──────────────────────┬──────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────┐
│                  Reply Engine                        │
│   Operator commands │ Memory packs │ Model routing   │
└──────────────────────┬──────────────────────────────┘
                       │
          ┌────────────┼────────────┐
          ▼            ▼            ▼
   ┌────────────┐ ┌─────────┐ ┌────────────┐
   │ Gemma 4    │ │ Gemma 4 │ │ Gemma 4    │
   │ E4B (local)│ │ E2B     │ │ 31B (cloud)│
   │ via Ollama │ │ (fast)  │ │ via Ollama │
   └────────────┘ └─────────┘ └────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────┐
│                    Web UI                            │
│   Conversation view │ Model routing │ Memory packs   │
│   Operator dashboard │ Session monitor               │
│   (This is what judges see at the live demo URL)     │
└─────────────────────────────────────────────────────┘
```

---

## Demo Walkthrough for Judges (3 scenarios)

### Scenario 1: Small Business Grant Discovery

**User context:** Maria, a Chicago small business owner, messages Savitar on Discord.

```
Maria: Hey, I just formed my LLC in Illinois. Are there any grants I can
       apply for? I'm a Latina woman running a food truck on the South Side.

Savitar: Great news — you're eligible for several programs:

1. **Neighborhood Opportunity Fund** — grants up to $250K for businesses
   on Chicago's South Side corridors. Rolling deadline.
2. **Comcast RISE** — $10K grant for women- or minority-owned businesses,
   plus marketing support and free internet. Quarterly deadline.
3. **Amber Grant for Women** — $10K monthly, $25K annual bonus. Open now.
4. **We Rise Together** — up to $12,500 for storefront expansion.

Want me to pull up eligibility details for any of these?
```

**What the web UI shows:** The conversation in real-time, the memory packs that were accessed (illinois-grants, chicago-programs), the model used (gemma4:e4b local), and the response latency.

### Scenario 2: Academic Collaboration

**User context:** An ISU student messages on Discord.

```
Student: I'm doing a capstone project on community health data in
         McLean County. Does ISU have any partnerships I could build on?

Savitar: ISU's Center for Civic Engagement tracks community partnerships
in the Collaboratory database. Relevant existing partnerships:

- **Chestnut Health Systems** — hosted ~90 ISU students in 2024,
  focused on career development and community care
- **Public health promotion** partners through the CCE network

For your capstone, I'd suggest contacting the CCE at (309) 438-1100
or ceceserve@IllinoisState.edu to register your project in
Collaboratory. That also makes you eligible for the Civic Engagement
Champion Award.

Want me to draft a project proposal template for your advisor?
```

### Scenario 3: Multi-Surface, Multi-Language

**User context:** A community college ESL instructor uses WhatsApp to coordinate.

```
Instructor (WhatsApp): My student Ahmed needs help with his FAFSA
                       application but his English is limited. Can you
                       help him in Arabic?

Savitar (WhatsApp): Of course. I'll switch to Arabic when talking with
                    Ahmed. Have him message me directly or I can help
                    translate your instructions.

Ahmed (WhatsApp): مرحبا، أحتاج مساعدة في تقديم طلب المساعدة المالية

Savitar (WhatsApp): مرحبا أحمد! سأساعدك في طلب FAFSA خطوة بخطوة.
                    أولاً، هل لديك رقم FSA ID؟
```

**What the web UI shows:** Same unified dashboard, now showing conversations from multiple surfaces in multiple languages, all routed through the same Gemma 4 instance.

---

## Demo URL Architecture

### Minimum Viable Demo (for submission)

```
Mac Mini M4 (home network)
├── Ollama (gemma4:e4b, gemma4:e2b)
├── Savitar runtime (Go binary)
├── Discord bot (live, connected)
├── Web UI (Go HTTP server on :8080)
└── Cloudflare Tunnel → public URL

Pixel 9 Pro Fold (Google Fi, mobile companion)
├── Web UI in Chrome (operator dashboard on 8" display)
├── Discord app (see bot conversations live)
├── WhatsApp (Fi number as Business endpoint)
└── Future: Gemma 4 E2B via LiteRT (on-device)
```

**Cloudflare Tunnel** (free tier) exposes the web UI to a public URL without opening router ports. Judges visit `https://savitar-demo.example.com` and see the operator dashboard. No login required for read-only demo view (a special demo mode).

The **Pixel Fold** is shown in the video as the mobile operator surface — unfold it and you see the same dashboard the judges see at the demo URL.

### Enhanced Demo (with more transports)

Same as above plus:
- iMessage bridge (native macOS, no additional cost)
- WhatsApp bridge (requires WhatsApp Business API or community bridge)

---

## What Judges See at the Demo URL

1. **Dashboard** — live conversation feed from all surfaces
2. **Model routing panel** — which model handled which request, local vs cloud
3. **Memory browser** — subject-area packs loaded, provenance metadata
4. **Budget calculator** — interactive tier selector showing hardware cost vs capability (including Apple Education $499 pricing)
5. **Architecture diagram** — the system map above, with live status indicators
6. **Mobile view** — screenshot or live view from Pixel 9 Pro Fold showing the same dashboard on a mobile device

---

## Technical Requirements for Demo

| Component | Status | Needed By |
|---|---|---|
| Discord bot | Shipped | Now |
| Web UI backend | Not started | Apr 20 |
| Web UI frontend | Not started | Apr 20 |
| Cloudflare Tunnel | Not started | Apr 27 |
| Demo mode (read-only, no auth) | Not started | Apr 27 |
| iMessage bridge | Contracted | May 4 |
| WhatsApp bridge | Contracted | May 4 |
| Memory packs loaded | Not started | Apr 13 |
| Pixel Fold mobile dashboard | Not started | Apr 27 |
| Google Fi WhatsApp number | Not started | May 4 |
