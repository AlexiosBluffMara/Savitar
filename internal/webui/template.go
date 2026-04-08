package webui

import (
	"fmt"
	"html/template"
	"time"
)

var dashboardTemplate = template.Must(template.New("dashboard").Funcs(template.FuncMap{
	"formatTime": func(ts time.Time) string {
		if ts.IsZero() {
			return "not yet"
		}
		return ts.Local().Format("Jan 2 2006 15:04 MST")
	},
	"formatUsage": func(multiplier float64) string {
		if multiplier == 0 {
			return "local-first"
		}
		return fmt.Sprintf("%.2fx", multiplier)
	},
	"formatBytes": func(size int64) string {
		switch {
		case size >= 1<<30:
			return fmt.Sprintf("%.1f GB", float64(size)/(1<<30))
		case size >= 1<<20:
			return fmt.Sprintf("%.1f MB", float64(size)/(1<<20))
		case size >= 1<<10:
			return fmt.Sprintf("%.1f KB", float64(size)/(1<<10))
		case size > 0:
			return fmt.Sprintf("%d B", size)
		default:
			return "0 B"
		}
	},
	"yesNo": func(value bool) string {
		if value {
			return "enabled"
		}
		return "planned"
	},
}).Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{.AgentName}} Operator Console</title>
  <style>
    :root {
      --paper: #f4ede2;
      --paper-strong: #fffaf3;
      --ink: #1d2423;
      --muted: #5c6763;
      --line: rgba(29, 36, 35, 0.12);
      --accent: #b64e2d;
      --accent-soft: rgba(182, 78, 45, 0.14);
      --secondary: #224d4d;
      --shadow: 0 24px 50px rgba(40, 31, 18, 0.12);
      --radius: 22px;
    }

    * { box-sizing: border-box; }
    body {
      margin: 0;
      color: var(--ink);
      font-family: "Avenir Next", "Helvetica Neue", sans-serif;
      background:
        radial-gradient(circle at top left, rgba(182, 78, 45, 0.18), transparent 32%),
        radial-gradient(circle at bottom right, rgba(34, 77, 77, 0.18), transparent 28%),
        linear-gradient(180deg, #efe5d6 0%, #f8f2ea 45%, #f4ede2 100%);
    }

    .page {
      max-width: 1280px;
      margin: 0 auto;
      padding: 32px 20px 48px;
    }

    .hero, .grid.two {
      display: grid;
      grid-template-columns: repeat(2, minmax(0, 1fr));
      gap: 20px;
    }

    .hero {
      margin-bottom: 20px;
      align-items: stretch;
    }

    .panel {
      border: 1px solid var(--line);
      border-radius: var(--radius);
      background: rgba(255, 250, 243, 0.82);
      backdrop-filter: blur(10px);
      box-shadow: var(--shadow);
      padding: 22px;
    }

    .eyebrow {
      margin: 0 0 10px;
      color: var(--accent);
      font-size: 0.82rem;
      font-weight: 700;
      letter-spacing: 0.18em;
      text-transform: uppercase;
    }

    h1, h2, h3 {
      margin: 0;
      font-weight: 700;
      letter-spacing: -0.03em;
    }

    h1 {
      font-size: clamp(2.5rem, 6vw, 4.5rem);
      line-height: 0.98;
      max-width: 10ch;
    }

    .lede, .section-copy, .surface-copy, .integration-copy, .source-copy, .small-copy, .meta-note {
      color: var(--muted);
      line-height: 1.6;
    }

    .lede {
      margin: 18px 0 0;
      max-width: 54ch;
      font-size: 1.04rem;
    }

    .kicker, .stat-label {
      color: var(--muted);
      font-size: 0.82rem;
      text-transform: uppercase;
      letter-spacing: 0.1em;
    }

    .meta-value, .stat-value {
      font-weight: 700;
      line-height: 1;
    }

    .meta-value {
      font-size: 1.7rem;
      margin-bottom: 6px;
    }

    .meta-list, .stats, .surface-grid, .source-list, .run-list, .review-list, .integrations-grid, .tier-highlights {
      list-style: none;
      margin: 0;
      padding: 0;
    }

    .meta-list {
      display: grid;
      gap: 12px;
    }

    .meta-list div {
      display: flex;
      justify-content: space-between;
      gap: 16px;
      padding-top: 12px;
      border-top: 1px solid var(--line);
    }

    .meta-list dd {
      margin: 0;
      text-align: right;
      font-weight: 700;
    }

    .stats {
      display: grid;
      grid-template-columns: repeat(5, minmax(0, 1fr));
      gap: 14px;
      margin-bottom: 20px;
    }

    .stat-card, .surface-card, .integration-card, .source-item, .run-item, .review-item, .tier-panel {
      border: 1px solid var(--line);
      border-radius: 18px;
      background: rgba(255, 255, 255, 0.62);
      padding: 16px;
    }

    .stat-value {
      font-size: 2rem;
      margin: 8px 0;
    }

    .section-head {
      display: flex;
      justify-content: space-between;
      gap: 16px;
      align-items: baseline;
      margin-bottom: 18px;
    }

    .surface-grid, .integrations-grid {
      display: grid;
      grid-template-columns: repeat(2, minmax(0, 1fr));
      gap: 14px;
    }

    .status-pill, .chip, .tier-button {
      display: inline-flex;
      align-items: center;
      border-radius: 999px;
      font-size: 0.8rem;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: 0.08em;
    }

    .status-pill {
      padding: 7px 11px;
      margin-bottom: 12px;
    }

    .status-pill.on {
      background: rgba(34, 77, 77, 0.12);
      color: var(--secondary);
    }

    .status-pill.off {
      background: var(--accent-soft);
      color: var(--accent);
    }

    .surface-title, .integration-title, .tier-name {
      font-size: 1.15rem;
      margin-bottom: 6px;
    }

    .profiles-table {
      width: 100%;
      border-collapse: collapse;
      margin-top: 8px;
    }

    .profiles-table th, .profiles-table td {
      text-align: left;
      padding: 13px 10px;
      border-bottom: 1px solid var(--line);
      vertical-align: top;
    }

    .profiles-table th {
      color: var(--muted);
      font-size: 0.84rem;
      text-transform: uppercase;
      letter-spacing: 0.08em;
    }

    .profiles-table td strong {
      display: block;
      margin-bottom: 4px;
    }

    .source-list, .run-list, .review-list {
      display: grid;
      gap: 12px;
      margin-top: 16px;
    }

    .source-meta, .run-meta, .review-meta, .tier-meta {
      display: flex;
      flex-wrap: wrap;
      gap: 10px;
      color: var(--muted);
      font-size: 0.88rem;
      margin-top: 10px;
    }

    .chip {
      padding: 6px 10px;
      background: rgba(29, 36, 35, 0.06);
    }

    .tier-controls {
      display: flex;
      flex-wrap: wrap;
      gap: 10px;
      margin: 20px 0 16px;
    }

    .tier-button {
      appearance: none;
      border: 1px solid var(--line);
      background: rgba(255, 255, 255, 0.7);
      color: var(--ink);
      padding: 10px 16px;
      cursor: pointer;
    }

    .tier-button.is-active {
      border-color: rgba(29, 36, 35, 0.22);
      background: rgba(255, 255, 255, 0.96);
    }

    .tier-panel { display: none; }
    .tier-panel.is-active { display: block; }

    .tier-highlights {
      display: grid;
      grid-template-columns: repeat(2, minmax(0, 1fr));
      gap: 10px;
      margin-top: 14px;
    }

    .tier-highlights li {
      border-radius: 16px;
      background: rgba(34, 77, 77, 0.08);
      padding: 12px 14px;
      line-height: 1.5;
    }

    .plan-list, .guardrail-list {
      color: var(--muted);
      line-height: 1.6;
      margin-top: 16px;
    }

    .footer {
      display: flex;
      justify-content: space-between;
      gap: 16px;
      align-items: flex-end;
      margin-top: 20px;
      color: var(--muted);
      font-size: 0.92rem;
    }

    .footer a {
      color: var(--secondary);
      font-weight: 700;
    }

    @media (max-width: 1100px) {
      .stats { grid-template-columns: repeat(3, minmax(0, 1fr)); }
    }

    @media (max-width: 920px) {
      .hero, .grid.two, .surface-grid, .integrations-grid, .tier-highlights {
        grid-template-columns: 1fr;
      }
      .stats { grid-template-columns: repeat(2, minmax(0, 1fr)); }
      .footer, .section-head {
        flex-direction: column;
        align-items: flex-start;
      }
    }

    @media (max-width: 640px) {
      .page { padding: 20px 14px 40px; }
      .stats { grid-template-columns: 1fr; }
      h1 { max-width: none; }
    }
  </style>
</head>
<body>
  <main class="page">
    <section class="hero">
      <div class="panel">
        <p class="eyebrow">Gemma 4 Good prototype</p>
        <h1>{{.AgentName}} Operator Console</h1>
        <p class="lede">A local-first dashboard for the current Savitar runtime: model routing, knowledge sources, integration readiness, and hackathon budget framing in one place. This build is intentionally honest about what is live versus what still needs authentication and production hardening.</p>
      </div>
      <aside class="panel">
        <div class="kicker">Launch mode</div>
        <div class="meta-value">{{.AuthMode}}</div>
        <p class="meta-note">{{if .DemoMode}}Local demo mode is active. This is appropriate for a laptop or Mac Mini demo, not for public exposure.{{else}}Authenticated mode is expected here once Google OAuth and signed sessions are implemented.{{end}}</p>
        <dl class="meta-list">
          <div><dt>Primary model</dt><dd>{{.LocalModel}}</dd></div>
          <div><dt>Provider</dt><dd>{{.LocalProvider}}</dd></div>
          <div><dt>Configured public URL</dt><dd>{{.PublicBaseURL}}</dd></div>
          <div><dt>Gateway plan</dt><dd>{{.Status.EnabledSurfaces}} / {{.Status.TotalSurfaces}} surfaces enabled</dd></div>
        </dl>
      </aside>
    </section>

    <section class="stats">
      <li class="stat-card">
        <div class="stat-label">Knowledge view</div>
        <div class="stat-value">{{if .Snapshot.KnowledgePacks}}{{len .Snapshot.KnowledgePacks}}{{else}}{{len .RepoSources}}{{end}}</div>
        <div class="small-copy">{{if .Snapshot.KnowledgePacks}}On-disk packs from the configured knowledge store.{{else}}Repo markdown fallback because the local pack store is empty in this checkout.{{end}}</div>
      </li>
      <li class="stat-card">
        <div class="stat-label">Review queue</div>
        <div class="stat-value">{{.Snapshot.PendingReviews}}</div>
        <div class="small-copy">Pending operator review items discovered from the local review directory.</div>
      </li>
      <li class="stat-card">
        <div class="stat-label">MCP reach</div>
        <div class="stat-value">{{.Status.EnabledMCPServers}}</div>
        <div class="small-copy">Configured MCP servers in the local runtime plan.</div>
      </li>
      <li class="stat-card">
        <div class="stat-label">Workspace assets</div>
        <div class="stat-value">{{.Status.AgentCount}} / {{.Status.SkillCount}}</div>
        <div class="small-copy">Discovered agents and skills that shape Savitar’s local behavior.</div>
      </li>
      <li class="stat-card">
        <div class="stat-label">Last refresh</div>
        <div class="stat-value">{{formatTime .GeneratedAt}}</div>
        <div class="small-copy">Rendered from current runtime state and local repository files.</div>
      </li>
    </section>

    <section class="grid two">
      <article class="panel">
        <div class="section-head">
          <div>
            <h2>Surface status</h2>
            <p class="section-copy">The hackathon story depends on one coherent system: Discord today, the operator web UI now, and later transports only after the core loop is stable.</p>
          </div>
        </div>
        <ul class="surface-grid">
          {{range .Surfaces}}
          <li class="surface-card">
            <div class="status-pill {{if .Enabled}}on{{else}}off{{end}}">{{yesNo .Enabled}}</div>
            <div class="surface-title">{{.Name}}</div>
            <div class="surface-copy">{{.Kind}} surface, identity {{.Identity}}.</div>
            <div class="source-meta"><span class="chip">{{.Details}}</span></div>
          </li>
          {{end}}
        </ul>
      </article>

      <article class="panel">
        <div class="section-head">
          <div>
            <h2>Model routing</h2>
            <p class="section-copy">Gemma 4 stays central. The local route is the default story, while the managed lanes remain clearly separated for more complex work.</p>
          </div>
        </div>
        <table class="profiles-table">
          <thead>
            <tr><th>Lane</th><th>Model</th><th>Usage</th></tr>
          </thead>
          <tbody>
            {{range .ModelProfiles}}
            <tr>
              <td><strong>{{.Name}}</strong><span class="small-copy">{{.Provider}}</span></td>
              <td><strong>{{.Model}}</strong><span class="small-copy">{{.Purpose}}</span></td>
              <td>{{formatUsage .UsageMultiplier}}</td>
            </tr>
            {{end}}
          </tbody>
        </table>
      </article>
    </section>

    <section class="grid two">
      <article class="panel">
        <div class="section-head">
          <div>
            <h2>Knowledge view</h2>
            <p class="section-copy">The prototype prefers on-disk packs when they exist. In a fresh checkout, it falls back to the repo markdown directories the runtime already searches for grounded evidence.</p>
          </div>
        </div>
        {{if .Snapshot.KnowledgePacks}}
        <ul class="source-list">
          {{range .Snapshot.KnowledgePacks}}
          <li class="source-item">
            <h3>{{.Name}}</h3>
            <p class="source-copy">Subject {{.Subject}}</p>
            <div class="source-meta">
              <span class="chip">{{formatTime .UpdatedAt}}</span>
              <span class="chip">{{formatBytes .SizeBytes}}</span>
              <span class="chip">{{.Path}}</span>
            </div>
          </li>
          {{end}}
        </ul>
        {{else}}
        <ul class="source-list">
          {{range .RepoSources}}
          <li class="source-item">
            <h3>{{.Name}}</h3>
            <p class="source-copy">{{.Section}} source</p>
            <div class="source-meta">
              <span class="chip">{{.Path}}</span>
              <span class="chip">{{formatTime .UpdatedAt}}</span>
              <span class="chip">{{formatBytes .SizeBytes}}</span>
            </div>
          </li>
          {{end}}
        </ul>
        {{end}}
      </article>

      <article class="panel">
        <div class="section-head">
          <div>
            <h2>Recent runtime state</h2>
            <p class="section-copy">This section stays truthful: it shows local session-backed state and the local review queue when present, rather than inventing live conversations.</p>
          </div>
        </div>
        {{if .Snapshot.RecentRuns}}
        <ul class="run-list">
          {{range .Snapshot.RecentRuns}}
          <li class="run-item">
            <h3>{{.ID}}</h3>
            <p class="source-copy">Surface {{.Surface}}, route {{.Route}}</p>
            <div class="run-meta">
              <span class="chip">conversation {{if .ConversationID}}{{.ConversationID}}{{else}}&lt;none&gt;{{end}}</span>
              <span class="chip">{{formatTime .StartedAt}}</span>
            </div>
          </li>
          {{end}}
        </ul>
        {{else}}
        <p class="section-copy">No durable run records are available yet. Initialize the local session or wire the orchestrator recorder to populate this panel.</p>
        {{end}}

        {{if .Reviews}}
        <ul class="review-list">
          {{range .Reviews}}
          <li class="review-item">
            <h3>{{.Reason}}</h3>
            <p class="source-copy">{{.Preview}}</p>
            <div class="review-meta">
              <span class="chip">{{.Surface}}</span>
              <span class="chip">{{formatTime .CreatedAt}}</span>
              <span class="chip">{{.ConversationID}}</span>
            </div>
          </li>
          {{end}}
        </ul>
        {{else}}
        <p class="section-copy">The local review queue is empty. The JSON API still supports future operator decisions once the orchestrator starts writing review items.</p>
        {{end}}
      </article>
    </section>

    <section class="panel">
      <div class="section-head">
        <div>
          <h2>Hackathon budget selector</h2>
          <p class="section-copy">The dashboard includes the pricing story the judges need to see: Mac Mini as the default local server, Pixel Fold as the mobile operator surface, and cloud only as an explicit tradeoff.</p>
        </div>
      </div>
      <div class="tier-controls">
        {{range $index, $tier := .BudgetTiers}}
        <button class="tier-button {{if eq $index 1}}is-active{{end}}" type="button" data-tier-button="{{$tier.Key}}">{{$tier.Name}}</button>
        {{end}}
      </div>
      {{range $index, $tier := .BudgetTiers}}
      <article class="tier-panel {{if eq $index 1}}is-active{{end}}" data-tier-panel="{{$tier.Key}}">
        <h3 class="tier-name">{{$tier.Name}}</h3>
        <div class="tier-meta">
          <span class="chip">{{$tier.Budget}}</span>
          <span class="chip">{{$tier.Hardware}}</span>
          <span class="chip">{{$tier.Audience}}</span>
        </div>
        <p class="section-copy">{{$tier.WhyItFits}}</p>
        <p class="small-copy"><strong>Routing:</strong> {{$tier.Routing}}</p>
        <ul class="tier-highlights">
          {{range $tier.Highlights}}
          <li>{{.}}</li>
          {{end}}
        </ul>
      </article>
      {{end}}
    </section>

    <section class="grid two">
      <article class="panel">
        <div class="section-head">
          <div>
            <h2>Execution plan</h2>
            <p class="section-copy">The operator console is only useful if it reflects the repo’s current build order rather than hand-waving beyond what exists.</p>
          </div>
        </div>
        <ol class="plan-list">{{range .Plan}}<li>{{.}}</li>{{end}}</ol>
      </article>

      <article class="panel">
        <div class="section-head">
          <div>
            <h2>Guardrails</h2>
            <p class="section-copy">These constraints keep the prototype compatible with the current ADRs and public-surface rules.</p>
          </div>
        </div>
        <ul class="guardrail-list">
          <li>Public exposure is still blocked on Google OAuth, signed sessions, rate limiting, and audit logging.</li>
          <li>Demo mode is explicit and local-only. It exists for hackathon validation and operator UX iteration, not for public deployment.</li>
          <li>The dashboard shows real runtime and repository state. When data does not exist yet, it says so instead of fabricating activity.</li>
          <li>Repo markdown fallback keeps the prototype useful even before the local knowledge store is populated on this checkout.</li>
        </ul>
      </article>
    </section>

    <section class="panel">
      <div class="section-head">
        <div>
          <h2>Integration readiness</h2>
          <p class="section-copy">The prototype pulls the same provider readiness data the CLI already exposes, so operators can see what is configured before a demo starts.</p>
        </div>
      </div>
      <ul class="integrations-grid">
        {{range .Integrations}}
        <li class="integration-card">
          <div class="status-pill {{if .Enabled}}on{{else}}off{{end}}">{{yesNo .Enabled}}</div>
          <div class="integration-title">{{.Name}}</div>
          <p class="integration-copy">{{.Details}}</p>
          <div class="source-meta">
            <span class="chip">auth {{.AuthSource}}</span>
            <span class="chip">credential {{if .CredentialPresent}}ready{{else}}missing{{end}}</span>
          </div>
        </li>
        {{end}}
      </ul>
    </section>

    <footer class="footer">
      <div>Rendered {{formatTime .GeneratedAt}}. Local operator prototype only.</div>
      <div><a href="/api/dashboard">Open JSON dashboard payload</a></div>
    </footer>
  </main>
  <script>
    const buttons = Array.from(document.querySelectorAll('[data-tier-button]'));
    const panels = Array.from(document.querySelectorAll('[data-tier-panel]'));

    function activateTier(key) {
      buttons.forEach((button) => {
        button.classList.toggle('is-active', button.dataset.tierButton === key);
      });
      panels.forEach((panel) => {
        panel.classList.toggle('is-active', panel.dataset.tierPanel === key);
      });
    }

    buttons.forEach((button) => {
      button.addEventListener('click', () => activateTier(button.dataset.tierButton));
    });
  </script>
</body>
</html>`))
