package webui

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alexiosbluffmara/savitar/internal/gateway"
	"github.com/alexiosbluffmara/savitar/internal/models"
	"github.com/alexiosbluffmara/savitar/internal/operator"
	savitarruntime "github.com/alexiosbluffmara/savitar/internal/runtime"
)

var errDashboardRuntimeNotWired = errors.New("dashboard runtime not wired")

// RepoSource is a markdown document the prototype can surface when the local
// knowledge-pack store has not been initialized yet.
type RepoSource struct {
	Name      string    `json:"name"`
	Section   string    `json:"section"`
	Path      string    `json:"path"`
	UpdatedAt time.Time `json:"updatedAt"`
	SizeBytes int64     `json:"sizeBytes"`
}

// BudgetTier summarizes one hackathon deployment tier.
type BudgetTier struct {
	Key        string   `json:"key"`
	Name       string   `json:"name"`
	Budget     string   `json:"budget"`
	Hardware   string   `json:"hardware"`
	Audience   string   `json:"audience"`
	Routing    string   `json:"routing"`
	WhyItFits  string   `json:"whyItFits"`
	Highlights []string `json:"highlights"`
}

// DashboardData is the richer payload used by the HTML dashboard and JSON API.
type DashboardData struct {
	GeneratedAt   time.Time                          `json:"generatedAt"`
	AgentName     string                             `json:"agentName"`
	AuthMode      string                             `json:"authMode"`
	DemoMode      bool                               `json:"demoMode"`
	WebUIEnabled  bool                               `json:"webUIEnabled"`
	PublicBaseURL string                             `json:"publicBaseURL"`
	LocalProvider string                             `json:"localProvider"`
	LocalModel    string                             `json:"localModel"`
	Status        savitarruntime.Status              `json:"status"`
	Surfaces      []gateway.SurfaceConfig            `json:"surfaces"`
	ModelProfiles []models.Profile                   `json:"modelProfiles"`
	Integrations  []savitarruntime.IntegrationStatus `json:"integrations"`
	Snapshot      operator.Snapshot                  `json:"snapshot"`
	Reviews       []operator.ReviewItem              `json:"reviews"`
	KnowledgeMode string                             `json:"knowledgeMode"`
	RepoSources   []RepoSource                       `json:"repoSources,omitempty"`
	BudgetTiers   []BudgetTier                       `json:"budgetTiers"`
	Plan          []string                           `json:"plan"`
}

func (s Service) dashboardData(ctx context.Context) (DashboardData, error) {
	if s.Runtime == nil {
		return DashboardData{}, errDashboardRuntimeNotWired
	}

	status, err := s.Runtime.Status()
	if err != nil {
		return DashboardData{}, err
	}

	data := DashboardData{
		GeneratedAt:   time.Now().UTC(),
		AgentName:     s.Runtime.Persona().Name,
		AuthMode:      s.authMode(),
		DemoMode:      s.DemoMode,
		WebUIEnabled:  s.Runtime.Config().Config.WebUI.Enabled,
		PublicBaseURL: s.Runtime.Config().Config.WebUI.PublicBaseURL,
		LocalProvider: s.Runtime.Config().Config.Models.LocalDefault.Provider,
		LocalModel:    s.Runtime.Config().Config.Models.LocalDefault.Model,
		Status:        status,
		Surfaces:      s.Runtime.GatewayPlan().Surfaces,
		ModelProfiles: s.Runtime.ModelProfiles(),
		Integrations:  s.Runtime.IntegrationStatuses(),
		BudgetTiers:   defaultBudgetTiers(),
		Plan:          s.Runtime.Plan(),
	}

	if s.Operator != nil {
		snapshot, err := s.Operator.Snapshot(ctx)
		if err != nil && !errors.Is(err, operator.ErrNotImplemented) {
			return DashboardData{}, err
		}
		if err == nil {
			data.Snapshot = snapshot
		}

		reviews, err := s.Operator.PendingReviews(ctx)
		if err != nil && !errors.Is(err, operator.ErrNotImplemented) {
			return DashboardData{}, err
		}
		if err == nil {
			data.Reviews = reviews
			if data.Snapshot.PendingReviews == 0 {
				data.Snapshot.PendingReviews = len(reviews)
			}
		}
	}

	if len(data.Snapshot.KnowledgePacks) > 0 {
		data.KnowledgeMode = "disk-packs"
		return data, nil
	}

	data.KnowledgeMode = "repo-markdown-fallback"
	sources, err := repoMarkdownSources(s.Runtime.RootDir(), s.Runtime.Config().Config.Knowledge.RepoMarkdownDirs, 8)
	if err != nil {
		return DashboardData{}, err
	}
	data.RepoSources = sources
	return data, nil
}

func (s Service) authMode() string {
	if s.DemoMode {
		return "local-demo"
	}
	if s.Runtime == nil {
		return "unconfigured"
	}
	return s.Runtime.Config().Config.WebUI.AuthProvider
}

func repoMarkdownSources(rootDir string, configured []string, limit int) ([]RepoSource, error) {
	seen := map[string]bool{}
	sources := make([]RepoSource, 0, len(configured))

	appendFile := func(path string, info os.FileInfo) {
		rel, err := filepath.Rel(rootDir, path)
		if err != nil {
			rel = path
		}
		rel = filepath.ToSlash(rel)
		if seen[rel] {
			return
		}
		seen[rel] = true

		section := filepath.ToSlash(filepath.Dir(rel))
		if section == "." {
			section = "repo-root"
		}
		sources = append(sources, RepoSource{
			Name:      strings.TrimSuffix(filepath.Base(rel), filepath.Ext(rel)),
			Section:   section,
			Path:      rel,
			UpdatedAt: info.ModTime(),
			SizeBytes: info.Size(),
		})
	}

	for _, entry := range configured {
		candidate := resolveRuntimePath(rootDir, entry)
		info, err := os.Stat(candidate)
		if err != nil {
			continue
		}

		if !info.IsDir() {
			if strings.EqualFold(filepath.Ext(candidate), ".md") {
				appendFile(candidate, info)
			}
			continue
		}

		err = filepath.Walk(candidate, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil || info == nil || info.IsDir() {
				return nil
			}
			if !strings.EqualFold(filepath.Ext(path), ".md") {
				return nil
			}
			appendFile(path, info)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	sort.Slice(sources, func(i, j int) bool {
		if sources[i].UpdatedAt.Equal(sources[j].UpdatedAt) {
			return sources[i].Path < sources[j].Path
		}
		return sources[i].UpdatedAt.After(sources[j].UpdatedAt)
	})

	if limit > 0 && len(sources) > limit {
		sources = sources[:limit]
	}
	return sources, nil
}

func defaultBudgetTiers() []BudgetTier {
	return []BudgetTier{
		{
			Key:       "tier-1",
			Name:      "Tier 1: Solo Builder",
			Budget:    "$100-$300",
			Hardware:  "Old Mac, Linux desktop, or spare local host",
			Audience:  "students, hobbyists, and small private pilots",
			Routing:   "local-only gemma4:e2b on Ollama",
			WhyItFits: "This is the honest floor: one low-cost box, one lightweight Gemma lane, one transport, and grounded knowledge from disk.",
			Highlights: []string{
				"Savitar CLI and Discord can run from a small local box with no hosted-model dependency",
				"Best for experimentation, light educator prep, and low-volume private use",
				"Grounded markdown packs still matter more than raw model size here",
				"Useful baseline for proving Savitar is local by design",
			},
		},
		{
			Key:       "tier-2",
			Name:      "Tier 2: Standard Operator Host",
			Budget:    "$499-$799",
			Hardware:  "Mac Mini M4 with 16GB unified memory",
			Audience:  "small businesses, community orgs, educators, and LLC operators",
			Routing:   "gemma4:e4b-it-q8_0 primary with gemma4:e2b as the lighter local lane",
			WhyItFits: "This is the recommended hackathon lane because it keeps Gemma 4 local, supports the operator web UI, and aligns with the Mac Mini story.",
			Highlights: []string{
				"Full Savitar runtime with the local Gemma 4 E4B lane front and center",
				"Strongest fit for the Ollama special-technology prize",
				"Credible for Chicago small-business operators and educator pilots",
				"This is the target host for the public hackathon demo",
			},
		},
		{
			Key:       "tier-3",
			Name:      "Tier 3: Community Hub",
			Budget:    "$1,200-$2,500",
			Hardware:  "Higher-memory Mac Mini, Mac Studio, or another always-on local host",
			Audience:  "educator teams, community labs, and shared operator deployments",
			Routing:   "the same local Gemma 4 lanes, with more concurrency and larger knowledge packs",
			WhyItFits: "This tier keeps the same local story while giving the operator enough headroom for shared use and deeper evidence packs.",
			Highlights: []string{
				"Supports multiple educator or community channels without changing the product story",
				"Attached SSD storage makes larger local packs and run history practical",
				"Better for shared review queues and operator oversight",
				"Still keeps inference local and auditable",
			},
		},
		{
			Key:       "tier-4",
			Name:      "Tier 4: Educator Live Mode",
			Budget:    "Tier 2 or Tier 3 plus an existing laptop or foldable phone",
			Hardware:  "Local host plus a portable dashboard surface",
			Audience:  "educators, workshop operators, and live moderators",
			Routing:   "Mac Mini stays the inference box while the dashboard follows the operator into the room",
			WhyItFits: "This tier turns the same local system into a strong Future of Education demo without pretending the phone is the server.",
			Highlights: []string{
				"The operator can run a live, reviewable question stream against curated class materials",
				"Responsive dashboard doubles as the live demo surface",
				"Keeps the education story concrete and easy to explain in 3 minutes",
				"Pairs naturally with local evidence packs and operator review",
			},
		},
	}
}
