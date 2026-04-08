package memory

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	headingPattern = regexp.MustCompile(`^(#{1,6})\s+(.+?)\s*$`)
	linkPattern    = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	tokenPattern   = regexp.MustCompile(`[a-z0-9]{3,}`)
)

// MarkdownSearchResult is a ranked set of markdown excerpts plus a compact
// graph summary that can be injected into an LLM prompt or shown in a CLI.
type MarkdownSearchResult struct {
	Query   string
	Results []MarkdownSectionMatch
	Graph   MarkdownGraph
}

// MarkdownSectionMatch is a single ranked local-repository excerpt.
type MarkdownSectionMatch struct {
	Path        string
	HeadingPath []string
	Snippet     string
	Score       int
}

// SourceLabel returns the human-readable label for a search match.
func (m MarkdownSectionMatch) SourceLabel() string {
	parts := []string{filepath.ToSlash(m.Path)}
	if len(m.HeadingPath) > 0 {
		parts = append(parts, strings.Join(m.HeadingPath, " > "))
	}
	return strings.Join(parts, " :: ")
}

// MarkdownGraph is the compact document-and-link graph derived from the local
// markdown repository for the top-ranked search results.
type MarkdownGraph struct {
	Nodes []MarkdownGraphNode
	Edges []MarkdownGraphEdge
}

// MarkdownGraphNode is a document node with a few representative headings.
type MarkdownGraphNode struct {
	Path     string
	Headings []string
}

// MarkdownGraphEdge is a directed internal markdown link between two docs.
type MarkdownGraphEdge struct {
	From string
	To   string
}

// ContextString renders the ranked excerpts plus the graph in a compact format
// suitable for LLM consumption.
func (r MarkdownSearchResult) ContextString() string {
	if len(r.Results) == 0 {
		return ""
	}
	lines := []string{"Local repository markdown evidence:"}
	for _, result := range r.Results {
		lines = append(lines, fmt.Sprintf("- %s", result.SourceLabel()))
		lines = append(lines, fmt.Sprintf("  %s", result.Snippet))
	}
	if graph := r.GraphString(); graph != "" {
		lines = append(lines, "")
		lines = append(lines, graph)
	}
	return strings.Join(lines, "\n")
}

// GraphString renders a compact markdown graph summary that gives the LLM a
// navigable view of which docs and internal links are near the query.
func (r MarkdownSearchResult) GraphString() string {
	if len(r.Graph.Nodes) == 0 {
		return ""
	}
	lines := []string{"Repository graph:"}
	for _, node := range r.Graph.Nodes {
		lines = append(lines, fmt.Sprintf("- %s", filepath.ToSlash(node.Path)))
		for _, heading := range node.Headings {
			lines = append(lines, fmt.Sprintf("  heading: %s", heading))
		}
	}
	for _, edge := range r.Graph.Edges {
		lines = append(lines, fmt.Sprintf("- edge: %s -> %s", filepath.ToSlash(edge.From), filepath.ToSlash(edge.To)))
	}
	return strings.Join(lines, "\n")
}

type markdownDoc struct {
	Path     string
	Headings []string
	Sections []markdownSection
	Links    []string
}

type markdownSection struct {
	HeadingPath []string
	Text        string
	Score       int
}

// SearchRepoMarkdown scans the configured markdown files and directories,
// returns the highest-ranked excerpts for query, and builds a lightweight doc
// graph for the top-ranked files.
func SearchRepoMarkdown(rootDir string, includePaths []string, query string, maxResults int, maxGraphEdges int) (MarkdownSearchResult, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return MarkdownSearchResult{}, fmt.Errorf("query is required")
	}
	if maxResults <= 0 {
		maxResults = 3
	}
	if maxGraphEdges <= 0 {
		maxGraphEdges = 6
	}

	docs, err := loadMarkdownDocs(rootDir, includePaths)
	if err != nil {
		return MarkdownSearchResult{}, err
	}

	tokens := tokenizeQuery(query)
	if len(tokens) == 0 {
		tokens = []string{strings.ToLower(query)}
	}

	results := rankSections(docs, query, tokens, maxResults)
	graph := buildGraph(docs, results, maxGraphEdges)

	return MarkdownSearchResult{Query: query, Results: results, Graph: graph}, nil
}

func loadMarkdownDocs(rootDir string, includePaths []string) ([]markdownDoc, error) {
	if len(includePaths) == 0 {
		includePaths = []string{"README.md", "docs/adr", "docs/roadmap", "docs/hackathon"}
	}

	seen := map[string]bool{}
	var docs []markdownDoc
	for _, includePath := range includePaths {
		resolved := filepath.Join(rootDir, includePath)
		info, err := os.Stat(resolved)
		if err != nil {
			continue
		}
		if !info.IsDir() {
			if strings.HasSuffix(strings.ToLower(resolved), ".md") {
				doc, err := parseMarkdownDoc(rootDir, resolved)
				if err == nil && !seen[doc.Path] {
					seen[doc.Path] = true
					docs = append(docs, doc)
				}
			}
			continue
		}
		err = filepath.WalkDir(resolved, func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() {
				return nil
			}
			if !strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
				return nil
			}
			doc, err := parseMarkdownDoc(rootDir, path)
			if err != nil || seen[doc.Path] {
				return nil
			}
			seen[doc.Path] = true
			docs = append(docs, doc)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return docs, nil
}

func parseMarkdownDoc(rootDir string, path string) (markdownDoc, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return markdownDoc{}, err
	}
	rel, err := filepath.Rel(rootDir, path)
	if err != nil {
		rel = path
	}
	rel = filepath.ToSlash(rel)
	content := strings.ReplaceAll(string(data), "\r\n", "\n")
	title := filepath.Base(rel)

	stack := []string{}
	sections := []markdownSection{}
	current := markdownSection{HeadingPath: []string{title}}
	links := []string{}
	headings := []string{}

	flush := func() {
		trimmed := strings.TrimSpace(current.Text)
		if trimmed == "" {
			return
		}
		current.Text = trimmed
		current.HeadingPath = append([]string(nil), current.HeadingPath...)
		sections = append(sections, current)
	}

	for _, line := range strings.Split(content, "\n") {
		for _, match := range linkPattern.FindAllStringSubmatch(line, -1) {
			target := normalizeLinkTarget(rel, match[2])
			if target != "" {
				links = append(links, target)
			}
		}

		if match := headingPattern.FindStringSubmatch(line); len(match) == 3 {
			flush()
			level := len(match[1])
			text := strings.TrimSpace(match[2])
			if level <= len(stack) {
				stack = stack[:level-1]
			}
			stack = append(stack, text)
			headings = append(headings, text)
			current = markdownSection{HeadingPath: append([]string{title}, stack...)}
			continue
		}

		if current.Text != "" {
			current.Text += "\n"
		}
		current.Text += line
	}
	flush()

	if len(headings) == 0 {
		headings = append(headings, title)
	}

	return markdownDoc{Path: rel, Headings: headings, Sections: sections, Links: uniqueSorted(links)}, nil
}

func rankSections(docs []markdownDoc, query string, tokens []string, maxResults int) []MarkdownSectionMatch {
	type scored struct {
		match MarkdownSectionMatch
	}
	var scoredResults []scored
	queryLower := strings.ToLower(query)
	for _, doc := range docs {
		pathLower := strings.ToLower(doc.Path)
		for _, section := range doc.Sections {
			headingLower := strings.ToLower(strings.Join(section.HeadingPath, " "))
			textLower := strings.ToLower(section.Text)
			score := 0
			for _, token := range tokens {
				score += 8 * strings.Count(pathLower, token)
				score += 5 * strings.Count(headingLower, token)
				score += strings.Count(textLower, token)
			}
			if strings.Contains(textLower, queryLower) {
				score += 4
			}
			if score <= 0 {
				continue
			}
			scoredResults = append(scoredResults, scored{match: MarkdownSectionMatch{
				Path:        doc.Path,
				HeadingPath: append([]string(nil), section.HeadingPath...),
				Snippet:     buildSnippet(section.Text, tokens, 280),
				Score:       score,
			}})
		}
	}

	sort.Slice(scoredResults, func(i, j int) bool {
		if scoredResults[i].match.Score == scoredResults[j].match.Score {
			return scoredResults[i].match.SourceLabel() < scoredResults[j].match.SourceLabel()
		}
		return scoredResults[i].match.Score > scoredResults[j].match.Score
	})

	results := make([]MarkdownSectionMatch, 0, maxResults)
	seen := map[string]bool{}
	for _, result := range scoredResults {
		key := result.match.SourceLabel()
		if seen[key] {
			continue
		}
		seen[key] = true
		results = append(results, result.match)
		if len(results) >= maxResults {
			break
		}
	}
	return results
}

func buildGraph(docs []markdownDoc, results []MarkdownSectionMatch, maxGraphEdges int) MarkdownGraph {
	if len(results) == 0 {
		return MarkdownGraph{}
	}
	docByPath := map[string]markdownDoc{}
	reverseLinks := map[string][]string{}
	for _, doc := range docs {
		docByPath[doc.Path] = doc
		for _, link := range doc.Links {
			reverseLinks[link] = append(reverseLinks[link], doc.Path)
		}
	}

	docOrder := []string{}
	selected := map[string]bool{}
	for _, result := range results {
		if selected[result.Path] {
			continue
		}
		selected[result.Path] = true
		docOrder = append(docOrder, result.Path)
	}
	for _, path := range append([]string(nil), docOrder...) {
		doc := docByPath[path]
		for _, link := range doc.Links {
			if !selected[link] {
				selected[link] = true
				docOrder = append(docOrder, link)
			}
		}
		for _, source := range reverseLinks[path] {
			if !selected[source] {
				selected[source] = true
				docOrder = append(docOrder, source)
			}
		}
	}

	nodes := make([]MarkdownGraphNode, 0, len(docOrder))
	edges := make([]MarkdownGraphEdge, 0, maxGraphEdges)
	seenEdge := map[string]bool{}
	for _, path := range docOrder {
		doc, ok := docByPath[path]
		if !ok {
			continue
		}
		headings := append([]string(nil), doc.Headings...)
		if len(headings) > 4 {
			headings = headings[:4]
		}
		nodes = append(nodes, MarkdownGraphNode{Path: path, Headings: headings})
		for _, link := range doc.Links {
			if len(edges) >= maxGraphEdges {
				break
			}
			if _, ok := docByPath[link]; !ok {
				continue
			}
			key := path + "->" + link
			if seenEdge[key] {
				continue
			}
			seenEdge[key] = true
			edges = append(edges, MarkdownGraphEdge{From: path, To: link})
		}
	}

	return MarkdownGraph{Nodes: nodes, Edges: edges}
}

func buildSnippet(text string, tokens []string, limit int) string {
	cleaned := strings.Join(strings.Fields(text), " ")
	if cleaned == "" {
		return ""
	}
	if limit <= 0 {
		limit = 280
	}
	cleanedLower := strings.ToLower(cleaned)
	start := 0
	for _, token := range tokens {
		if idx := strings.Index(cleanedLower, token); idx >= 0 {
			start = idx
			if start > 80 {
				start -= 80
			} else {
				start = 0
			}
			break
		}
	}
	end := start + limit
	if end > len(cleaned) {
		end = len(cleaned)
	}
	snippet := cleaned[start:end]
	if start > 0 {
		snippet = "..." + snippet
	}
	if end < len(cleaned) {
		snippet += "..."
	}
	return snippet
}

func tokenizeQuery(query string) []string {
	matches := tokenPattern.FindAllString(strings.ToLower(query), -1)
	return uniqueSorted(matches)
}

func normalizeLinkTarget(docPath string, target string) string {
	target = strings.TrimSpace(target)
	if target == "" || strings.HasPrefix(target, "#") || strings.Contains(target, "://") || strings.HasPrefix(target, "mailto:") {
		return ""
	}
	if idx := strings.IndexAny(target, "#?"); idx >= 0 {
		target = target[:idx]
	}
	if target == "" {
		return ""
	}
	base := filepath.Dir(filepath.FromSlash(docPath))
	resolved := filepath.Clean(filepath.Join(base, filepath.FromSlash(target)))
	return filepath.ToSlash(resolved)
}

func uniqueSorted(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	seen := map[string]bool{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}
