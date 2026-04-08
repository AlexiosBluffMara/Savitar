package respond

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/alexiosbluffmara/savitar/internal/chat"
	savitarruntime "github.com/alexiosbluffmara/savitar/internal/runtime"
)

var (
	duckResultLinkPattern    = regexp.MustCompile(`nofollow" class="result__a" href="([^"]+)"[^>]*>(.*?)</a>`)
	duckResultSnippetPattern = regexp.MustCompile(`<a[^>]*class="result__snippet"[^>]*>(.*?)</a>`)
	htmlTagPattern           = regexp.MustCompile(`<.*?>`)
)

type repoSearcher interface {
	Search(context.Context, string) (repoEvidence, error)
}

type liveLookup interface {
	Lookup(context.Context, string) (liveEvidence, error)
}

type repoEvidence struct {
	Context string
	Sources []string
}

type liveEvidence struct {
	ToolContext chat.ToolContext
	Sources     []string
}

type markdownRepoSearcher struct {
	rt *savitarruntime.Runtime
}

func (s markdownRepoSearcher) Search(_ context.Context, query string) (repoEvidence, error) {
	if s.rt == nil {
		return repoEvidence{}, nil
	}
	lookupQuery := normalizeRepoSearchQuery(query)
	if lookupQuery == "" {
		lookupQuery = query
	}
	result, err := s.rt.RepoMarkdownSearch(lookupQuery)
	if err != nil || len(result.Results) == 0 {
		return repoEvidence{}, err
	}
	sources := make([]string, 0, len(result.Results))
	for _, match := range result.Results {
		sources = append(sources, "Local repo: "+match.SourceLabel())
	}
	return repoEvidence{Context: result.ContextString(), Sources: sources}, nil
}

type duckDuckGoLookup struct {
	httpClient httpClient
	enabled    bool
	timeout    time.Duration
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

func newLiveLookup(rt *savitarruntime.Runtime) liveLookup {
	if rt == nil {
		return nil
	}
	cfg := rt.Config().Config.Knowledge
	if !cfg.EnableLiveWebLookup {
		return nil
	}
	timeout := time.Duration(cfg.WebLookupTimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 8 * time.Second
	}
	provider := strings.TrimSpace(strings.ToLower(cfg.LiveWebProvider))
	if provider == "" || provider == "duckduckgo-html" || provider == "duckduckgo-json" || provider == "duckduckgo" {
		return duckDuckGoLookup{
			httpClient: &http.Client{Timeout: timeout},
			enabled:    true,
			timeout:    timeout,
		}
	}
	return nil
}

func (l duckDuckGoLookup) Lookup(ctx context.Context, query string) (liveEvidence, error) {
	if !l.enabled {
		return liveEvidence{}, nil
	}
	lookupQuery := normalizeLiveLookupQuery(query)
	if lookupQuery == "" {
		lookupQuery = strings.TrimSpace(query)
	}
	result, err := l.lookupDuckDuckGoJSON(ctx, lookupQuery)
	if err != nil {
		result, err = l.lookupWikipedia(ctx, lookupQuery)
		if err != nil {
			return liveEvidence{}, err
		}
	}
	contextText := fmt.Sprintf("Live web lookup (%s): %s\n%s", result.URL, result.Title, result.Snippet)
	return liveEvidence{
		ToolContext: chat.ToolContext{
			ServerName: "web",
			ToolName:   "duckduckgo-json",
			Result:     contextText,
		},
		Sources: []string{fmt.Sprintf("Live web: %s (%s)", result.Title, result.URL)},
	}, nil
}

func (l duckDuckGoLookup) lookupDuckDuckGoJSON(ctx context.Context, query string) (duckDuckGoResult, error) {
	searchURL := "https://api.duckduckgo.com/?q=" + url.QueryEscape(query) + "&format=json&no_html=1&skip_disambig=1"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
	if err != nil {
		return duckDuckGoResult{}, err
	}
	req.Header.Set("User-Agent", "savitar-prototype/0.1")

	resp, err := l.httpClient.Do(req)
	if err != nil {
		return duckDuckGoResult{}, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return duckDuckGoResult{}, fmt.Errorf("duckduckgo instant answer returned HTTP %d", resp.StatusCode)
	}
	var decoded struct {
		Heading      string `json:"Heading"`
		Abstract     string `json:"Abstract"`
		AbstractText string `json:"AbstractText"`
		AbstractURL  string `json:"AbstractURL"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return duckDuckGoResult{}, err
	}
	title := strings.TrimSpace(decoded.Heading)
	if title == "" {
		title = strings.TrimSpace(query)
	}
	snippet := strings.TrimSpace(decoded.AbstractText)
	if snippet == "" {
		snippet = strings.TrimSpace(decoded.Abstract)
	}
	if snippet == "" || strings.TrimSpace(decoded.AbstractURL) == "" {
		return duckDuckGoResult{}, fmt.Errorf("duckduckgo instant answer returned no abstract result")
	}
	return duckDuckGoResult{Title: title, Snippet: snippet, URL: strings.TrimSpace(decoded.AbstractURL)}, nil
}

func (l duckDuckGoLookup) lookupWikipedia(ctx context.Context, query string) (duckDuckGoResult, error) {
	searchURL := "https://en.wikipedia.org/w/api.php?action=opensearch&search=" + url.QueryEscape(query) + "&limit=1&namespace=0&format=json"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
	if err != nil {
		return duckDuckGoResult{}, err
	}
	req.Header.Set("User-Agent", "savitar-prototype/0.1")
	resp, err := l.httpClient.Do(req)
	if err != nil {
		return duckDuckGoResult{}, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return duckDuckGoResult{}, fmt.Errorf("wikipedia opensearch returned HTTP %d", resp.StatusCode)
	}
	var decoded []any
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return duckDuckGoResult{}, err
	}
	if len(decoded) < 4 {
		return duckDuckGoResult{}, fmt.Errorf("wikipedia opensearch returned incomplete result")
	}
	titles, _ := decoded[1].([]any)
	urls, _ := decoded[3].([]any)
	if len(titles) == 0 || len(urls) == 0 {
		return duckDuckGoResult{}, fmt.Errorf("wikipedia opensearch returned no titles")
	}
	title, _ := titles[0].(string)
	pageURL, _ := urls[0].(string)
	if title == "" || pageURL == "" {
		return duckDuckGoResult{}, fmt.Errorf("wikipedia opensearch returned empty title or URL")
	}
	return duckDuckGoResult{Title: title, Snippet: "Wikipedia search result for the requested topic.", URL: pageURL}, nil
}

type duckDuckGoResult struct {
	Title   string
	Snippet string
	URL     string
}

func parseDuckDuckGoResult(body string) (duckDuckGoResult, error) {
	linkMatch := duckResultLinkPattern.FindStringSubmatch(body)
	if len(linkMatch) != 3 {
		return duckDuckGoResult{}, fmt.Errorf("no duckduckgo result link found")
	}
	snippetMatch := duckResultSnippetPattern.FindStringSubmatch(body)
	title := cleanHTMLText(linkMatch[2])
	snippet := ""
	if len(snippetMatch) == 2 {
		snippet = cleanHTMLText(snippetMatch[1])
	}
	resolvedURL := decodeDuckDuckGoURL(linkMatch[1])
	if title == "" || resolvedURL == "" {
		return duckDuckGoResult{}, fmt.Errorf("duckduckgo result was incomplete")
	}
	return duckDuckGoResult{Title: title, Snippet: snippet, URL: resolvedURL}, nil
}

func cleanHTMLText(value string) string {
	value = htmlTagPattern.ReplaceAllString(value, "")
	value = html.UnescapeString(value)
	value = strings.Join(strings.Fields(value), " ")
	return strings.TrimSpace(value)
}

func decodeDuckDuckGoURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if strings.HasPrefix(raw, "//") {
		raw = "https:" + raw
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	uddg := parsed.Query().Get("uddg")
	if uddg == "" {
		return raw
	}
	decoded, err := url.QueryUnescape(uddg)
	if err != nil {
		return uddg
	}
	return decoded
}

func normalizeLiveLookupQuery(input string) string {
	normalized := strings.TrimSpace(input)
	if normalized == "" {
		return ""
	}
	for _, separator := range []string{" and ", " but ", ",", ";"} {
		parts := strings.SplitN(normalized, separator, 2)
		if len(parts) == 2 {
			normalized = parts[0]
			break
		}
	}
	lower := strings.ToLower(normalized)
	for _, prefix := range []string{"who is ", "what is ", "when is ", "where is ", "why is ", "how does ", "how do ", "tell me about ", "look up ", "lookup ", "search for ", "search ", "find ", "explain ", "summarize "} {
		if strings.HasPrefix(lower, prefix) {
			normalized = strings.TrimSpace(normalized[len(prefix):])
			break
		}
	}
	normalized = strings.Trim(normalized, " ?.!")
	return normalized
}

func normalizeRepoSearchQuery(input string) string {
	normalized := strings.TrimSpace(input)
	if normalized == "" {
		return ""
	}
	for _, separator := range []string{" and ", " but ", ";"} {
		parts := strings.SplitN(normalized, separator, 2)
		if len(parts) == 2 {
			candidate := strings.TrimSpace(parts[1])
			if len(strings.Fields(candidate)) >= 2 {
				normalized = candidate
				break
			}
		}
	}
	lower := strings.ToLower(normalized)
	for _, prefix := range []string{"who is ", "what is ", "when is ", "where is ", "why is ", "how does ", "how do ", "tell me about ", "look up ", "lookup ", "search for ", "search ", "find ", "explain ", "summarize "} {
		if strings.HasPrefix(lower, prefix) {
			normalized = strings.TrimSpace(normalized[len(prefix):])
			break
		}
	}
	normalized = strings.Trim(normalized, " ?.!")
	return normalized
}

func shouldUseLiveLookup(input string) bool {
	lower := strings.ToLower(strings.TrimSpace(input))
	if lower == "" {
		return false
	}
	if strings.Contains(lower, "?") {
		return true
	}
	for _, prefix := range []string{"who ", "what ", "when ", "where ", "why ", "how ", "tell me about ", "look up ", "lookup ", "search ", "search for ", "find ", "explain ", "summarize ", "latest ", "current "} {
		if strings.HasPrefix(lower, prefix) {
			return true
		}
	}
	return false
}

func sourceFooters(sources []string) []string {
	if len(sources) == 0 {
		return nil
	}
	seen := map[string]bool{}
	out := make([]string, 0, len(sources))
	for _, source := range sources {
		source = strings.TrimSpace(source)
		if source == "" || seen[source] {
			continue
		}
		seen[source] = true
		out = append(out, source)
	}
	return out
}
