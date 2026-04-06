package memory

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Store reads and writes memory packs to a directory tree at
// <rootDir>/<subject>/<name>.md
type Store struct {
	rootDir string
}

// NewStore creates a Store rooted at rootDir.
func NewStore(rootDir string) *Store {
	return &Store{rootDir: rootDir}
}

// Write serializes pack to disk. The pack name is used as the file base name
// (sanitized). Directories are created as needed.
func (s *Store) Write(pack Pack) error {
	if strings.TrimSpace(pack.Name) == "" {
		return fmt.Errorf("pack name is required")
	}
	if strings.TrimSpace(pack.Subject) == "" {
		return fmt.Errorf("pack subject is required")
	}

	subjectDir := filepath.Join(s.rootDir, sanitize(pack.Subject))
	if err := os.MkdirAll(subjectDir, 0o755); err != nil {
		return err
	}

	path := filepath.Join(subjectDir, sanitize(pack.Name)+".md")
	return os.WriteFile(path, []byte(pack.Markdown()), 0o644)
}

// ReadBySubject loads all packs whose subject matches subject (case-insensitive).
func (s *Store) ReadBySubject(subject string) ([]Pack, error) {
	subjectDir := filepath.Join(s.rootDir, sanitize(subject))
	entries, err := os.ReadDir(subjectDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	packs := make([]Pack, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		path := filepath.Join(subjectDir, entry.Name())
		pack, err := readFile(path)
		if err != nil {
			continue
		}
		packs = append(packs, pack)
	}

	// Sort by updated time, newest first.
	sort.Slice(packs, func(i, j int) bool {
		return packs[i].UpdatedAt.After(packs[j].UpdatedAt)
	})
	return packs, nil
}

// ListAll returns a summary of all packs across all subjects.
func (s *Store) ListAll() ([]PackMeta, error) {
	subjectDirs, err := os.ReadDir(s.rootDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var metas []PackMeta
	for _, subjectEntry := range subjectDirs {
		if !subjectEntry.IsDir() {
			continue
		}
		subject := subjectEntry.Name()
		entries, err := os.ReadDir(filepath.Join(s.rootDir, subject))
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}
			path := filepath.Join(s.rootDir, subject, entry.Name())
			pack, err := readFile(path)
			if err != nil {
				continue
			}
			info, _ := entry.Info()
			size := int64(0)
			if info != nil {
				size = info.Size()
			}
			metas = append(metas, PackMeta{
				Name:      pack.Name,
				Subject:   pack.Subject,
				UpdatedAt: pack.UpdatedAt,
				Path:      path,
				SizeBytes: size,
			})
		}
	}

	sort.Slice(metas, func(i, j int) bool {
		return metas[i].UpdatedAt.After(metas[j].UpdatedAt)
	})
	return metas, nil
}

// Get reads a single pack by subject and name.
func (s *Store) Get(subject, name string) (Pack, error) {
	path := filepath.Join(s.rootDir, sanitize(subject), sanitize(name)+".md")
	return readFile(path)
}

// PackMeta is a lightweight summary of a pack without the full body.
type PackMeta struct {
	Name      string
	Subject   string
	UpdatedAt time.Time
	Path      string
	SizeBytes int64
}

func readFile(path string) (Pack, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Pack{}, err
	}
	return FromMarkdown(string(data)), nil
}

// sanitize turns a string into a safe filesystem name: lowercase, hyphens only.
func sanitize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || r == '_' || r == '.':
			b.WriteRune('-')
		default:
			b.WriteRune('-')
		}
	}
	result := b.String()
	// Collapse multiple hyphens.
	for strings.Contains(result, "--") {
		result = strings.ReplaceAll(result, "--", "-")
	}
	return strings.Trim(result, "-")
}
