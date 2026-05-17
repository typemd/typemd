package core

import (
	"strings"
	"testing"
)

func TestParseWikiLinks(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected []WikiLink
	}{
		{
			name:     "no links",
			body:     "This is plain text with no links.",
			expected: nil,
		},
		{
			name: "single link",
			body: "See [[book/golang-in-action]] for details.",
			expected: []WikiLink{
				{Target: "book/golang-in-action", DisplayText: ""},
			},
		},
		{
			name: "link with display text",
			body: "See [[book/golang-in-action|Go in Action]] for details.",
			expected: []WikiLink{
				{Target: "book/golang-in-action", DisplayText: "Go in Action"},
			},
		},
		{
			name: "multiple links",
			body: "Read [[book/clean-code]] and [[person/robert-martin]].",
			expected: []WikiLink{
				{Target: "book/clean-code", DisplayText: ""},
				{Target: "person/robert-martin", DisplayText: ""},
			},
		},
		{
			name: "duplicate links deduplicated",
			body: "See [[book/clean-code]] and again [[book/clean-code]].",
			expected: []WikiLink{
				{Target: "book/clean-code", DisplayText: ""},
			},
		},
		{
			name:     "empty brackets ignored",
			body:     "This [[]] is empty.",
			expected: nil,
		},
		{
			name: "link in multiline content",
			body: `First paragraph.

See [[book/clean-code]] for more.

Also check [[person/john-doe|John Doe]].`,
			expected: []WikiLink{
				{Target: "book/clean-code", DisplayText: ""},
				{Target: "person/john-doe", DisplayText: "John Doe"},
			},
		},
		{
			name:     "single bracket prefix ignored in target",
			body:     "This [[[nested]]] is weird.",
			expected: []WikiLink{{Target: "[nested", DisplayText: ""}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseWikiLinks(tt.body)
			if len(got) != len(tt.expected) {
				t.Fatalf("got %d links, want %d: %v", len(got), len(tt.expected), got)
			}
			for i, link := range got {
				if link.Target != tt.expected[i].Target {
					t.Errorf("link[%d].Target = %q, want %q", i, link.Target, tt.expected[i].Target)
				}
				if link.DisplayText != tt.expected[i].DisplayText {
					t.Errorf("link[%d].DisplayText = %q, want %q", i, link.DisplayText, tt.expected[i].DisplayText)
				}
			}
		})
	}
}

func TestRenderWikiLinks(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		expected string
	}{
		{
			name:     "no links unchanged",
			body:     "Plain text with no links.",
			expected: "Plain text with no links.",
		},
		{
			name:     "link with display text",
			body:     "By [[person/gene-kim-01jqr3k5mpbvn8e0f2g7h9txyz|Gene Kim]].",
			expected: "By Gene Kim.",
		},
		{
			name:     "link without display text strips ULID",
			body:     "See [[book/the-phoenix-project-01jqr3k5mpbvn8e0f2g7h9txyz]].",
			expected: "See book/the-phoenix-project.",
		},
		{
			name:     "link without ULID kept as-is",
			body:     "See [[book/clean-code]].",
			expected: "See book/clean-code.",
		},
		{
			name:     "multiple links",
			body:     "By [[person/gene-kim-01jqr3k5mpbvn8e0f2g7h9txyz|Gene Kim]] about [[book/the-phoenix-project-01jqr3k5mpbvn8e0f2g7h9txyz]].",
			expected: "By Gene Kim about book/the-phoenix-project.",
		},
		{
			name:     "multiline body",
			body:     "First line.\n\nSee [[book/clean-code|Clean Code]] for more.\n\nAlso [[person/john-doe]].",
			expected: "First line.\n\nSee Clean Code for more.\n\nAlso person/john-doe.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RenderWikiLinks(tt.body)
			if got != tt.expected {
				t.Errorf("RenderWikiLinks() =\n%q\nwant:\n%q", got, tt.expected)
			}
		})
	}
}

func TestRenderWikiLinksStyled(t *testing.T) {
	style := func(s string) string {
		return "[" + s + "]"
	}

	body := "By [[person/gene-kim-01jqr3k5mpbvn8e0f2g7h9txyz|Gene Kim]] about [[book/clean-code]]."
	got := RenderWikiLinksStyled(body, style)
	expected := "By [Gene Kim] about [book/clean-code]."
	if got != expected {
		t.Errorf("RenderWikiLinksStyled() =\n%q\nwant:\n%q", got, expected)
	}
}

func TestRenderWikiLinksStyled_nilStyle(t *testing.T) {
	body := "See [[book/clean-code|Clean Code]]."
	got := RenderWikiLinksStyled(body, nil)
	if !strings.Contains(got, "Clean Code") {
		t.Errorf("expected display text, got: %q", got)
	}
}

func TestResolveWikiLinkTarget(t *testing.T) {
	// Use 26-char ULID suffixes to match ulidSuffixPattern
	const (
		cleanCodeID = "book/clean-code-01jqr3k5mpbvn8e0f2g7h9tx00"
		golangID    = "book/golang-01jqr3k5mpbvn8e0f2g7h9tx01"
		myNoteID    = "note/my-note-01jqr3k5mpbvn8e0f2g7h9tx02"
		johnDoeID   = "person/john-doe-01jqr3k5mpbvn8e0f2g7h9tx03"
	)
	diskIDs := map[string]bool{
		cleanCodeID: true,
		golangID:    true,
		myNoteID:    true,
		johnDoeID:   true,
	}
	nameIndex := map[string]map[string][]string{
		"book": {
			"clean-code": {cleanCodeID},
			"golang":     {golangID},
		},
		"note": {
			"my-note": {myNoteID},
		},
		"person": {
			"john-doe": {johnDoeID},
		},
	}

	tests := []struct {
		name       string
		target     string
		sourceType string
		wantID     string
		wantChange bool
		wantErr    bool
	}{
		{
			name:       "full ID resolves exactly",
			target:     cleanCodeID,
			sourceType: "note",
			wantID:     cleanCodeID,
		},
		{
			name:       "full ID not found",
			target:     "book/nonexistent-01jqr3k5mpbvn8e0f2g7h9tx99",
			sourceType: "note",
			wantID:     "",
		},
		{
			name:       "type-qualified name resolves",
			target:     "book/clean-code",
			sourceType: "note",
			wantID:     cleanCodeID,
			wantChange: true,
		},
		{
			name:       "type-qualified name not found",
			target:     "book/nonexistent",
			sourceType: "note",
			wantErr:    true,
		},
		{
			name:       "same-type shorthand resolves",
			target:     "my-note",
			sourceType: "note",
			wantID:     myNoteID,
			wantChange: true,
		},
		{
			name:       "same-type shorthand not found",
			target:     "nonexistent",
			sourceType: "note",
			wantErr:    true,
		},
		{
			name:       "empty target",
			target:     "",
			sourceType: "note",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := resolveWikiLinkTarget(tt.target, tt.sourceType, diskIDs, nameIndex)
			if res.resolvedID != tt.wantID {
				t.Errorf("resolvedID = %q, want %q", res.resolvedID, tt.wantID)
			}
			if res.changed != tt.wantChange {
				t.Errorf("changed = %v, want %v", res.changed, tt.wantChange)
			}
			if (res.err != nil) != tt.wantErr {
				t.Errorf("err = %v, wantErr = %v", res.err, tt.wantErr)
			}
		})
	}
}

func TestResolveWikiLinkTarget_Ambiguous(t *testing.T) {
	const (
		golangIntroID = "book/golang-intro-01jqr3k5mpbvn8e0f2g7h9tx04"
		golangGuideID = "book/golang-guide-01jqr3k5mpbvn8e0f2g7h9tx05"
	)
	diskIDs := map[string]bool{
		golangIntroID: true,
		golangGuideID: true,
	}
	nameIndex := map[string]map[string][]string{
		"book": {
			"golang": {golangIntroID, golangGuideID},
		},
	}

	res := resolveWikiLinkTarget("book/golang", "note", diskIDs, nameIndex)
	if res.resolvedID != "" {
		t.Errorf("resolvedID = %q, want empty", res.resolvedID)
	}
	if res.err == nil {
		t.Fatal("expected error for ambiguous match")
	}
	ame, ok := res.err.(*AmbiguousMatchError)
	if !ok {
		t.Fatalf("expected AmbiguousMatchError, got %T", res.err)
	}
	if len(ame.Matches) != 2 {
		t.Errorf("matches = %d, want 2", len(ame.Matches))
	}
}

func TestExpandWikiLinksInBody(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		resolutions map[string]string
		wantBody    string
		wantCount   int
	}{
		{
			name:        "no resolutions",
			body:        "See [[clean-code]].",
			resolutions: map[string]string{},
			wantBody:    "See [[clean-code]].",
			wantCount:   0,
		},
		{
			name: "single shorthand expanded",
			body: "See [[clean-code]].",
			resolutions: map[string]string{
				"clean-code": "book/clean-code-01abc",
			},
			wantBody:  "See [[book/clean-code-01abc]].",
			wantCount: 1,
		},
		{
			name: "display text preserved",
			body: "By [[clean-code|Clean Code]].",
			resolutions: map[string]string{
				"clean-code": "book/clean-code-01abc",
			},
			wantBody:  "By [[book/clean-code-01abc|Clean Code]].",
			wantCount: 1,
		},
		{
			name: "multiple links in one body",
			body: "See [[clean-code]] and [[golang]].",
			resolutions: map[string]string{
				"clean-code": "book/clean-code-01abc",
				"golang":     "book/golang-01def",
			},
			wantBody:  "See [[book/clean-code-01abc]] and [[book/golang-01def]].",
			wantCount: 2,
		},
		{
			name: "unresolved target not modified",
			body: "See [[clean-code]] and [[unknown]].",
			resolutions: map[string]string{
				"clean-code": "book/clean-code-01abc",
			},
			wantBody:  "See [[book/clean-code-01abc]] and [[unknown]].",
			wantCount: 1,
		},
		{
			name: "same target twice",
			body: "First [[clean-code]] and second [[clean-code]].",
			resolutions: map[string]string{
				"clean-code": "book/clean-code-01abc",
			},
			wantBody:  "First [[book/clean-code-01abc]] and second [[book/clean-code-01abc]].",
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody, gotCount := expandWikiLinksInBody(tt.body, tt.resolutions)
			if gotBody != tt.wantBody {
				t.Errorf("body =\n%q\nwant:\n%q", gotBody, tt.wantBody)
			}
			if gotCount != tt.wantCount {
				t.Errorf("count = %d, want %d", gotCount, tt.wantCount)
			}
		})
	}
}

// ── Wiki-link resolution via alias tests ───────────────────────────────────

func TestResolveWikiLinkTarget_ByAlias(t *testing.T) {
	const golangID = "book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9tx00"
	diskIDs := map[string]bool{golangID: true}

	aliasSlug := Slugify("Go 語言")
	nameIndex := map[string]map[string][]string{
		"book": {
			"golang-in-action": {golangID},
			aliasSlug:          {golangID},
		},
	}

	res := resolveWikiLinkTarget("book/"+aliasSlug, "note", diskIDs, nameIndex)
	if res.resolvedID != golangID {
		t.Errorf("resolvedID = %q, want %q", res.resolvedID, golangID)
	}
	if !res.changed {
		t.Error("expected changed = true for alias resolution")
	}
}

func TestResolveWikiLinkTarget_AliasLowerPriorityThanName(t *testing.T) {
	const (
		bookAID = "book/book-a-01jqr3k5mpbvn8e0f2g7h9tx01"
		bookBID = "book/book-b-01jqr3k5mpbvn8e0f2g7h9tx02"
	)
	diskIDs := map[string]bool{bookAID: true, bookBID: true}

	nameIndex := map[string]map[string][]string{
		"book": {
			"clean-code": {bookAID, bookBID},
		},
	}

	res := resolveWikiLinkTarget("book/clean-code", "note", diskIDs, nameIndex)
	if res.resolvedID != "" {
		t.Errorf("resolvedID = %q, want empty (ambiguous)", res.resolvedID)
	}
	if res.err == nil {
		t.Error("expected error for ambiguous alias/name collision")
	}
}
