package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alexiosbluffmara/savitar/internal/shell"
)

// Summary is the structured output of a repository analysis run.
type Summary struct {
	URL        string     `json:"url"`
	CloneDir   string     `json:"cloneDir"`
	CommitHash string     `json:"commitHash"`
	Branch     string     `json:"branch"`
	LogSnippet string     `json:"logSnippet"`
	FileTree   string     `json:"fileTree"`
	Stats      RepoStats  `json:"stats"`
	Provenance Provenance `json:"provenance"`
}

// RepoStats holds counts from the repository.
type RepoStats struct {
	CommitCount int `json:"commitCount"`
	FileCount   int `json:"fileCount"`
}

// Provenance records who ran the analysis and when.
type Provenance struct {
	SourceURL  string `json:"sourceURL"`
	CommitHash string `json:"commitHash"`
	AnalyzedAt string `json:"analyzedAt"`
	AnalyzedBy string `json:"analyzedBy"`
}

// Analyzer runs git-backed repository analysis under a shell policy.
type Analyzer struct {
	policy    *shell.Policy
	outputDir string
}

// NewAnalyzer creates an Analyzer. outputDir is where summaries are written.
func NewAnalyzer(outputDir string) *Analyzer {
	return &Analyzer{
		policy:    shell.DefaultPolicy(),
		outputDir: outputDir,
	}
}

// Analyze clones (or re-uses) a repository, collects metadata, writes a
// JSON summary to outputDir, and returns the summary.
func (a *Analyzer) Analyze(ctx context.Context, repoURL string) (*Summary, error) {
	if strings.TrimSpace(repoURL) == "" {
		return nil, fmt.Errorf("repo URL is required")
	}

	// Create a temp work directory for the clone.
	workDir, err := os.MkdirTemp("", "savitar-repo-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(workDir)

	cloneDir := filepath.Join(workDir, "repo")

	// Clone the repository.
	cloneResult := a.policy.Run(ctx, "git", "clone", "--depth=50", "--", repoURL, cloneDir)
	if cloneResult.Err != nil {
		return nil, fmt.Errorf("git clone: %w", cloneResult.Err)
	}
	if cloneResult.ExitCode != 0 {
		return nil, fmt.Errorf("git clone failed (exit %d): %s", cloneResult.ExitCode, strings.TrimSpace(cloneResult.Stderr))
	}

	policy := &shell.Policy{WorkDir: cloneDir}

	// Get HEAD commit hash.
	hashResult := policy.Run(ctx, "git", "rev-parse", "HEAD")
	commitHash := strings.TrimSpace(hashResult.Stdout)

	// Get current branch.
	branchResult := policy.Run(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	branch := strings.TrimSpace(branchResult.Stdout)

	// Get commit log (last 20 entries, one-line format).
	logResult := policy.Run(ctx, "git", "log", "--oneline", "-20")
	logSnippet := strings.TrimSpace(logResult.Stdout)

	// Count commits.
	countResult := policy.Run(ctx, "git", "rev-list", "--count", "HEAD")
	commitCount := parseIntLine(countResult.Stdout)

	// List top-level files (non-recursive for readability).
	findResult := policy.Run(ctx, "find", cloneDir, "-maxdepth", "3",
		"-not", "-path", "*/\\.git/*",
		"-not", "-name", ".git",
		"-type", "f",
	)
	fileLines := nonEmptyLines(findResult.Stdout)
	// Strip the clone dir prefix for readability.
	trimmed := make([]string, 0, len(fileLines))
	for _, line := range fileLines {
		rel := strings.TrimPrefix(line, cloneDir+"/")
		trimmed = append(trimmed, rel)
	}
	fileTree := strings.Join(trimmed, "\n")

	summary := &Summary{
		URL:        repoURL,
		CloneDir:   cloneDir,
		CommitHash: commitHash,
		Branch:     branch,
		LogSnippet: logSnippet,
		FileTree:   fileTree,
		Stats: RepoStats{
			CommitCount: commitCount,
			FileCount:   len(fileLines),
		},
		Provenance: Provenance{
			SourceURL:  repoURL,
			CommitHash: commitHash,
			AnalyzedAt: time.Now().UTC().Format(time.RFC3339),
			AnalyzedBy: "savitar",
		},
	}

	if err := a.write(summary); err != nil {
		return summary, fmt.Errorf("write summary: %w", err)
	}

	return summary, nil
}

// write serializes the summary to a JSON file in outputDir.
func (a *Analyzer) write(summary *Summary) error {
	if a.outputDir == "" {
		return nil
	}
	if err := os.MkdirAll(a.outputDir, 0o755); err != nil {
		return err
	}

	slug := repoSlug(summary.URL)
	name := fmt.Sprintf("%s-%s.json", slug, summary.Provenance.AnalyzedAt[:10])
	path := filepath.Join(a.outputDir, name)

	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func repoSlug(url string) string {
	s := strings.TrimSuffix(url, ".git")
	s = strings.TrimPrefix(s, "https://")
	s = strings.TrimPrefix(s, "http://")
	s = strings.TrimPrefix(s, "git@")
	s = strings.ReplaceAll(s, ":", "/")
	s = strings.ReplaceAll(s, "/", "-")
	s = strings.Trim(s, "-")
	if len(s) > 60 {
		s = s[len(s)-60:]
	}
	return s
}

func parseIntLine(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	var n int
	_, _ = fmt.Sscanf(s, "%d", &n)
	return n
}

func nonEmptyLines(s string) []string {
	var out []string
	for _, line := range strings.Split(s, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			out = append(out, line)
		}
	}
	return out
}
