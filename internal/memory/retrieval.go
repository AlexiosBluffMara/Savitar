package memory

import (
	"context"
	"time"
)

// Query describes the context request the orchestrator sends to the knowledge
// layer before composing a reply. Later passes can enrich it with per-surface
// policy, operator overrides, or freshness requirements.
type Query struct {
	Subject        string
	Terms          []string
	ConversationID string
	MaxPacks       int
	MaxExcerpts    int
	RequireSources bool
}

// Excerpt is the minimal evidence fragment passed forward to composition. The
// goal is to keep the retrieval layer responsible for what was found while the
// reply layer decides how much of it should be quoted or cited.
type Excerpt struct {
	Subject   string
	PackName  string
	Text      string
	Source    string
	Locator   string
	Score     float64
	UpdatedAt time.Time
}

// RetrievalResult is the normalized output of one knowledge lookup.
type RetrievalResult struct {
	Query    Query
	Excerpts []Excerpt
	Packs    []PackMeta
	Subjects []string
}

// SourceRecord tracks the origin material behind a pack, excerpt, or later
// live evidence import. It gives future indexing and operator review work a
// stable record that can be shown without re-parsing the original content.
type SourceRecord struct {
	ID             string
	Subject        string
	Kind           string
	Title          string
	Origin         string
	Locator        string
	LastCheckedAt  time.Time
	ProvenanceNote string
}

// Retriever resolves a query into ranked excerpts and related pack metadata.
type Retriever interface {
	Retrieve(context.Context, Query) (RetrievalResult, error)
}

// Indexer updates whatever search or ranking structure the knowledge layer uses
// to make later retrieval fast and reviewable.
type Indexer interface {
	IndexPack(context.Context, Pack) error
	DeletePack(context.Context, string, string) error
}

// SourceCatalog stores source provenance separately from pack bodies so later
// tooling can answer questions like "where did this knowledge come from?"
// without reparsing markdown files.
type SourceCatalog interface {
	UpsertSource(context.Context, SourceRecord) error
	ListSources(context.Context, string) ([]SourceRecord, error)
}
