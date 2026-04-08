package webui

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alexiosbluffmara/savitar/internal/gateway"
	"github.com/alexiosbluffmara/savitar/internal/operator"
	savitarruntime "github.com/alexiosbluffmara/savitar/internal/runtime"
	"github.com/alexiosbluffmara/savitar/internal/session"
)

// DemoAuthenticator is an explicit local-only authenticator for the current
// prototype. It should never be treated as a public auth mode.
type DemoAuthenticator struct {
	Email string
	Roles []string
}

// SessionFromRequest returns an always-approved local demo session.
func (a DemoAuthenticator) SessionFromRequest(*http.Request) (Session, error) {
	email := strings.TrimSpace(a.Email)
	if email == "" {
		email = "demo@savitar.local"
	}
	roles := append([]string(nil), a.Roles...)
	if len(roles) == 0 {
		roles = []string{"demo-operator"}
	}
	return Session{Email: email, Approved: true, Roles: roles}, nil
}

// RuntimeStore adapts the current runtime state into the operator.Store
// interface so the prototype can show real data without waiting for the full
// orchestrator backend.
type RuntimeStore struct {
	Runtime *savitarruntime.Runtime
}

// NewRuntimeStore returns an operator store backed by the runtime's current
// local state.
func NewRuntimeStore(rt *savitarruntime.Runtime) *RuntimeStore {
	return &RuntimeStore{Runtime: rt}
}

// DashboardSnapshot returns the current operator dashboard payload.
func (s *RuntimeStore) DashboardSnapshot(ctx context.Context) (operator.Snapshot, error) {
	_ = ctx
	if s == nil || s.Runtime == nil {
		return operator.Snapshot{}, operator.ErrNotImplemented
	}

	report, err := s.Runtime.Session()
	if err != nil {
		return operator.Snapshot{}, err
	}
	packs, err := s.Runtime.MemoryList()
	if err != nil {
		return operator.Snapshot{}, err
	}
	reviews, err := s.PendingApprovals(context.Background())
	if err != nil {
		return operator.Snapshot{}, err
	}

	return operator.Snapshot{
		GeneratedAt:    time.Now().UTC(),
		RecentRuns:     recentRunsFromSession(report),
		PendingReviews: len(reviews),
		KnowledgePacks: packs,
	}, nil
}

// PendingApprovals reads any queued review items from the configured review
// directory. Missing directories are treated as an empty queue.
func (s *RuntimeStore) PendingApprovals(ctx context.Context) ([]operator.ReviewItem, error) {
	_ = ctx
	if s == nil || s.Runtime == nil {
		return nil, operator.ErrNotImplemented
	}

	dir := resolveRuntimePath(s.Runtime.RootDir(), s.Runtime.Config().Config.Operator.ReviewQueueDir)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	items := make([]operator.ReviewItem, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var item operator.ReviewItem
		if err := json.Unmarshal(data, &item); err != nil {
			continue
		}
		items = append(items, item)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})
	return items, nil
}

// RecordOperatorDecision persists a decision into a local decisions directory so
// demo actions are inspectable even before the full operator backend exists.
func (s *RuntimeStore) RecordOperatorDecision(ctx context.Context, decision operator.Decision) error {
	_ = ctx
	if s == nil || s.Runtime == nil {
		return operator.ErrNotImplemented
	}
	if decision.DecidedAt.IsZero() {
		decision.DecidedAt = time.Now().UTC()
	}

	dir := filepath.Join(resolveRuntimePath(s.Runtime.RootDir(), s.Runtime.Config().Config.Operator.ReviewQueueDir), "decisions")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	name := safeFileSegment(decision.ReviewID)
	if name == "" {
		name = "decision"
	}
	path := filepath.Join(dir, fmt.Sprintf("%s-%s.json", decision.DecidedAt.Format("20060102T150405Z"), name))
	data, err := json.MarshalIndent(decision, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func recentRunsFromSession(report session.Report) []operator.RunSummary {
	if !report.Exists {
		return nil
	}
	updatedAt, _ := time.Parse(time.RFC3339, report.State.UpdatedAt)
	if updatedAt.IsZero() {
		updatedAt = time.Now().UTC()
	}
	return []operator.RunSummary{
		{
			ID:             "active-session",
			Surface:        gateway.Surface(report.State.CurrentSurface),
			ConversationID: report.State.ActiveConversationID,
			StartedAt:      updatedAt,
			Route:          report.State.CurrentModelProfile,
		},
	}
}

func resolveRuntimePath(rootDir string, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Clean(filepath.Join(rootDir, path))
}

func safeFileSegment(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var builder strings.Builder
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z':
			builder.WriteRune(r)
		case r >= '0' && r <= '9':
			builder.WriteRune(r)
		default:
			builder.WriteRune('-')
		}
	}
	result := builder.String()
	for strings.Contains(result, "--") {
		result = strings.ReplaceAll(result, "--", "-")
	}
	return strings.Trim(result, "-")
}
