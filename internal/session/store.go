package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type State struct {
	Version              int    `json:"version"`
	CurrentSurface       string `json:"currentSurface"`
	CurrentModelProfile  string `json:"currentModelProfile"`
	ActiveConversationID string `json:"activeConversationID,omitempty"`
	LastCommand          string `json:"lastCommand,omitempty"`
	UpdatedAt            string `json:"updatedAt,omitempty"`
}

type Report struct {
	Path   string
	Exists bool
	State  State
}

// Summary is a lightweight descriptor for a session file without full state.
type Summary struct {
	Path      string
	Name      string
	UpdatedAt time.Time
	Surface   string
	Exists    bool
}

type Store struct {
	path string
}

func NewStore(path string) Store {
	return Store{path: path}
}

func DefaultState() State {
	return State{
		Version:             1,
		CurrentSurface:      "cli",
		CurrentModelProfile: "local-default",
		UpdatedAt:           time.Now().UTC().Format(time.RFC3339),
	}
}

func (s Store) Load() (Report, error) {
	report := Report{
		Path:  s.path,
		State: DefaultState(),
	}

	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return report, nil
		}

		return Report{}, err
	}

	report.Exists = true
	if err := json.Unmarshal(data, &report.State); err != nil {
		return Report{}, err
	}

	return report, nil
}

func (s Store) Init(state State) (Report, error) {
	if state == (State{}) {
		state = DefaultState()
	}
	state.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return Report{}, err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return Report{}, err
	}

	data = append(data, '\n')
	if err := os.WriteFile(s.path, data, 0o600); err != nil {
		return Report{}, err
	}

	return Report{Path: s.path, Exists: true, State: state}, nil
}

// Save persists an updated state to disk without changing the path.
// It is safe to call even if the file has not been initialized yet.
func (s Store) Save(state State) error {
	state.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, append(data, '\n'), 0o600)
}

// List returns summaries for all *.json session files in the same directory
// as this store's path. Files that cannot be parsed are skipped.
func (s Store) List() ([]Summary, error) {
	dir := filepath.Dir(s.path)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var summaries []Summary
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var state State
		if err := json.Unmarshal(data, &state); err != nil {
			continue
		}
		updatedAt := time.Time{}
		if state.UpdatedAt != "" {
			updatedAt, _ = time.Parse(time.RFC3339, state.UpdatedAt)
		}
		name := strings.TrimSuffix(entry.Name(), ".json")
		summaries = append(summaries, Summary{
			Path:      path,
			Name:      name,
			UpdatedAt: updatedAt,
			Surface:   state.CurrentSurface,
			Exists:    true,
		})
	}

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].UpdatedAt.After(summaries[j].UpdatedAt)
	})
	return summaries, nil
}

// Prune deletes session files (other than the store's own path) that have not
// been updated within maxAge. It returns the count of pruned files.
func (s Store) Prune(maxAge time.Duration) (int, error) {
	if maxAge <= 0 {
		return 0, nil
	}
	summaries, err := s.List()
	if err != nil {
		return 0, err
	}

	cutoff := time.Now().UTC().Add(-maxAge)
	pruned := 0
	for _, sum := range summaries {
		if sum.Path == s.path {
			continue // never prune the primary session file
		}
		if !sum.UpdatedAt.IsZero() && sum.UpdatedAt.Before(cutoff) {
			if err := os.Remove(sum.Path); err == nil {
				pruned++
			}
		}
	}
	return pruned, nil
}
