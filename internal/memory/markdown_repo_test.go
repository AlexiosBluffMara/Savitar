package memory

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSearchRepoMarkdownFindsRelevantSectionAndGraph(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("# Home\n\nSee the [guide](docs/guide.md).\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	guide := "# Guide\n\n## Installation\n\nThis section explains installation for Savitar operators.\n\n## Usage\n\nUse the operator console.\n"
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte(guide), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := SearchRepoMarkdown(root, []string{"README.md", "docs"}, "installation operator", 3, 4)
	if err != nil {
		t.Fatalf("SearchRepoMarkdown returned error: %v", err)
	}
	if len(result.Results) == 0 {
		t.Fatal("expected at least one markdown result")
	}
	if result.Results[0].Path != "docs/guide.md" {
		t.Fatalf("unexpected top result path: %q", result.Results[0].Path)
	}
	if !strings.Contains(strings.ToLower(result.Results[0].Snippet), "installation") {
		t.Fatalf("expected installation snippet, got %q", result.Results[0].Snippet)
	}
	graph := result.GraphString()
	if !strings.Contains(graph, "README.md") {
		t.Fatalf("expected graph to mention README.md, got %q", graph)
	}
	if !strings.Contains(graph, "README.md -> docs/guide.md") {
		t.Fatalf("expected graph to include README to guide edge, got %q", graph)
	}
}

func TestSearchRepoMarkdownRequiresQuery(t *testing.T) {
	_, err := SearchRepoMarkdown(t.TempDir(), nil, "", 3, 4)
	if err == nil {
		t.Fatal("expected query validation error")
	}
}
