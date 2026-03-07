package core

import (
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
