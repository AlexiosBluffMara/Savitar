# Community Access Narrative

## Discord Is the New IRC — And That Matters

### The IRC Legacy

In the 1990s and early 2000s, IRC (Internet Relay Chat) was the backbone of technical community. Linux kernel developers coordinated patches on Freenode. Open-source projects lived and died by their channel activity. Academic research groups used IRC for asynchronous collaboration across time zones. The culture IRC fostered was distinctive: pseudonymous but persistent identity, channel-specific norms, operator hierarchies that rewarded contribution over credentials, and a deeply meritocratic ethos where your ideas mattered more than your title.

IRC was also radically accessible. A terminal, a network connection, and `irssi` — that was enough. No app store, no account verification, no phone number. Just connect and participate.

### Discord as IRC's Successor

Discord inherited IRC's best qualities and made them accessible to people who would never open a terminal. Servers replace networks. Channels replace rooms. Roles replace operator flags. But the culture is the same: communities form around shared interests, governance is local, and the tools are free.

For Savitar, Discord is the primary community surface because:

1. **Zero-cost entry** — anyone with a browser or phone can join
2. **Channel granularity** — separate spaces for different topics, roles, and access levels
3. **Bot integration** — first-class API for AI agents, no workarounds needed
4. **Voice channels** — real-time discussion when text isn't enough
5. **Already adopted** — students, gamers, open-source contributors, and small business communities are already there

The "Alexios Bluff Mara" Discord server uses this structure: `#main` for project coordination, `#artemis-talk` for agent interaction, `#rules` for governance, and role-based access to sensitive channels like `#secrets`.

### Identity and Exclusivity in Academic Networks

For academic and professional communities, open access is insufficient. When building a network for ISU students, faculty, and community partners, you need:

- **Verified identity** — real names tied to institutional email or verified community membership
- **Invitation-only gates** — a faculty advisor or community partner vouches for new members
- **Role-based access** — students see different channels than faculty, who see different channels than community partners
- **Audit trails** — who said what, when, to whom — required for academic integrity and grant reporting

Discord supports all of this through its role and permission system. Savitar enforces it through operator configuration: `requireMention`, `allowedChannels`, `respondOnlyToVerified` (planned).

The exclusivity isn't gatekeeping — it's trust-building. When Maria posts a question about her LLC filing in a verified channel and gets an answer from Savitar that cites Illinois statute, she needs to know that conversation is private to her community, not scraped for training data, and backed by a local model she can inspect.

---

## Three Streams, One Agent: Discord + iMessage + WhatsApp

### Why Three Platforms

| Platform | Audience | Geography | Strength |
|---|---|---|---|
| **Discord** | Students, developers, community orgs | Global, US-heavy | Rich channels, bots, voice, free |
| **iMessage** | American consumers, professionals | US-dominant | Native to iPhone, encrypted, trusted, zero-friction |
| **WhatsApp** | International communities, immigrant business owners | Global, dominant outside US | 2B+ users, encrypted, group chats, business accounts |

No single platform reaches everyone. A Chicago PoC small business owner might prefer iMessage for personal queries and WhatsApp for communicating with suppliers. An ISU student lives on Discord. A community college instructor uses WhatsApp groups for their ESL class.

The ecosystem coverage is deliberate:
- **Apple side** — Mac Mini runs the server, iMessage is the transport, Apple Education pricing makes it affordable for schools
- **Google side** — Pixel 9 Pro Fold is the mobile operator device, Google Fi provides connectivity, the web UI uses Google login, and the future LiteRT path puts Gemma 4 on-device
- **Open side** — Discord and WhatsApp work on everything, Ollama runs on Linux too

### The Unified Gateway

Savitar's gateway normalizes all three streams into a single message envelope:

```
InboundMessage {
  Surface:    "discord" | "imessage" | "whatsapp"
  SenderID:   platform-specific identifier
  Channel:    channel/thread/conversation ID
  Content:    text, image, or audio
  Timestamp:  UTC
  Identity:   verified name if available
}
```

The reply engine processes every message the same way:
1. Check for operator commands (help, status, ping)
2. Load conversation history for this sender
3. Load relevant memory packs by subject
4. Generate a reply using Gemma 4 via Ollama
5. Route the reply back through the originating surface

The web UI shows the operator a unified view: all conversations across all surfaces in one dashboard. Filter by platform, user, topic, or time. This is what judges see in the live demo.

### Why This Architecture Matters for "Good"

A community organization running Savitar doesn't need to choose between platforms. They serve everyone through the channels those people already use. The agent maintains context across surfaces — if Maria asks about her LLC on WhatsApp and then follows up on Discord, Savitar remembers.

This is digital equity in practice: **meet people where they are, on the tools they already have, in the language they prefer, with privacy they can verify.**

---

## The Web UI as the Hackathon Demo Surface

The hackathon requires a publicly accessible live demo. The web UI serves double duty:

1. **For judges** — a clean dashboard showing real conversations, model routing decisions, memory pack usage, and multi-surface activity
2. **For operators** — the daily interface for managing the agent, monitoring conversations, and updating knowledge

The web UI does *not* replace messaging. It's the operator's window into the agent that lives on Discord, WhatsApp, and iMessage. This distinction matters: users interact through platforms they know, operators manage through a purpose-built interface.

Authentication: Google login with explicit session validation, per ADR 0005. No anonymous access. Operator must be approved.

---

## The Mobile Companion: Pixel 9 Pro Fold on Google Fi

The web UI isn't only for desktops. The operator's primary mobile surface is a **Google Pixel 9 Pro Fold** on **Google Fi**. The Fold's 8" inner display shows the full operator dashboard in split-screen — conversations on one side, model routing on the other.

Why this combination matters:

1. **Field access** — an ISU faculty member at a community partner site can check on student interactions, update memory packs, and approve new users from their pocket
2. **Demo presentations** — unfold the Pixel at a hackathon table or community meeting and show the full dashboard without a laptop
3. **Google Fi international** — the Fi number serves as the WhatsApp Business number, connecting international communities
4. **Future on-device inference** — the Tensor G4 can run Gemma 4 E2B via LiteRT for completely offline, phone-local AI when connectivity is unreliable
5. **Two-ecosystem story** — the project deliberately spans Apple (Mac Mini server + iMessage) and Google (Pixel Fold + Fi + Android + LiteRT), showing that the architecture isn't locked to one vendor

This dual-ecosystem approach is a competitive advantage. Most hackathon entries pick one platform. Savitar bridges both — because the communities it serves use both.
