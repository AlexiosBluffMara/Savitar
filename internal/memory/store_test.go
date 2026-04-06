package memory

import (
	"os"
	"testing"
)

func TestStoreWriteAndRead(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(dir)

	pack := New("test-pack", "testing", "This is the body of the test pack.")
	if err := store.Write(pack); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}

	packs, err := store.ReadBySubject("testing")
	if err != nil {
		t.Fatalf("ReadBySubject returned error: %v", err)
	}
	if len(packs) != 1 {
		t.Fatalf("expected 1 pack, got %d", len(packs))
	}
	if packs[0].Name != "test-pack" {
		t.Errorf("expected name %q, got %q", "test-pack", packs[0].Name)
	}
	if packs[0].Subject != "testing" {
		t.Errorf("expected subject %q, got %q", "testing", packs[0].Subject)
	}
	if packs[0].Body != "This is the body of the test pack." {
		t.Errorf("unexpected body: %q", packs[0].Body)
	}
}

func TestStoreListAll(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(dir)

	_ = store.Write(New("pack-a", "subjectA", "body A"))
	_ = store.Write(New("pack-b", "subjectB", "body B"))

	metas, err := store.ListAll()
	if err != nil {
		t.Fatalf("ListAll returned error: %v", err)
	}
	if len(metas) != 2 {
		t.Errorf("expected 2 metas, got %d", len(metas))
	}
}

func TestStoreReadMissingSubjectReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(dir)

	packs, err := store.ReadBySubject("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(packs) != 0 {
		t.Errorf("expected 0 packs, got %d", len(packs))
	}
}

func TestStoreGet(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(dir)

	original := New("my-pack", "mysubject", "hello body")
	_ = store.Write(original)

	retrieved, err := store.Get("mysubject", "my-pack")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if retrieved.Body != original.Body {
		t.Errorf("body mismatch: got %q, want %q", retrieved.Body, original.Body)
	}
}

func TestPackMarkdownRoundTrip(t *testing.T) {
	original := New("round-trip", "test", "Some content here.")
	original.Source = "https://example.com"

	md := original.Markdown()
	parsed := FromMarkdown(md)

	if parsed.Name != original.Name {
		t.Errorf("name: got %q, want %q", parsed.Name, original.Name)
	}
	if parsed.Subject != original.Subject {
		t.Errorf("subject: got %q, want %q", parsed.Subject, original.Subject)
	}
	if parsed.Body != original.Body {
		t.Errorf("body: got %q, want %q", parsed.Body, original.Body)
	}
}

func TestSanitize(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{"Hello World", "hello-world"},
		{"repo/analysis", "repo-analysis"},
		{"  spaced  ", "spaced"},
		{"multiple---hyphens", "multiple-hyphens"},
	}
	for _, c := range cases {
		got := sanitize(c.in)
		if got != c.out {
			t.Errorf("sanitize(%q) = %q, want %q", c.in, got, c.out)
		}
	}
}

func init() {
	// Silence unused import warning.
	_ = os.DevNull
}
