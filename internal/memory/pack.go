package memory

import (
	"fmt"
	"strings"
	"time"
)

// Pack is a named, versioned unit of durable knowledge with provenance.
type Pack struct {
	Name      string    `json:"name"`
	Subject   string    `json:"subject"`
	Source    string    `json:"source,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Body      string    `json:"body"`
}

// frontmatter serializes the Pack metadata as a simple key: value header
// block suitable for embedding at the top of a Markdown file.
func (p Pack) frontmatter() string {
	lines := []string{
		"---",
		fmt.Sprintf("name: %s", p.Name),
		fmt.Sprintf("subject: %s", p.Subject),
	}
	if p.Source != "" {
		lines = append(lines, fmt.Sprintf("source: %s", p.Source))
	}
	lines = append(lines,
		fmt.Sprintf("created: %s", p.CreatedAt.UTC().Format(time.RFC3339)),
		fmt.Sprintf("updated: %s", p.UpdatedAt.UTC().Format(time.RFC3339)),
		"---",
	)
	return strings.Join(lines, "\n")
}

// Markdown returns the complete file content: frontmatter + body.
func (p Pack) Markdown() string {
	return p.frontmatter() + "\n\n" + strings.TrimSpace(p.Body) + "\n"
}

// parseFrontmatter extracts key-value pairs from a "--- ... ---" frontmatter block.
// It returns the metadata map and the body text after the closing "---".
func parseFrontmatter(content string) (map[string]string, string) {
	meta := map[string]string{}
	if !strings.HasPrefix(content, "---") {
		return meta, content
	}

	// Find closing "---"
	rest := content[3:]
	rest = strings.TrimLeft(rest, "\r\n")
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return meta, content
	}

	header := rest[:end]
	body := rest[end+4:] // skip "\n---"
	body = strings.TrimLeft(body, "\r\n")

	for _, line := range strings.Split(header, "\n") {
		line = strings.TrimSpace(line)
		if idx := strings.Index(line, ": "); idx > 0 {
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+2:])
			meta[key] = val
		}
	}
	return meta, body
}

// FromMarkdown parses a Pack from a Markdown file produced by Pack.Markdown.
func FromMarkdown(content string) Pack {
	meta, body := parseFrontmatter(content)
	p := Pack{
		Name:    meta["name"],
		Subject: meta["subject"],
		Source:  meta["source"],
		Body:    strings.TrimSpace(body),
	}
	if t, err := time.Parse(time.RFC3339, meta["created"]); err == nil {
		p.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339, meta["updated"]); err == nil {
		p.UpdatedAt = t
	}
	return p
}

// New creates a new Pack with the current time set for both timestamps.
func New(name, subject, body string) Pack {
	now := time.Now().UTC()
	return Pack{
		Name:      name,
		Subject:   subject,
		CreatedAt: now,
		UpdatedAt: now,
		Body:      body,
	}
}
