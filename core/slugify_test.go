package core

import "testing"

func TestSlugify(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"simple spaces", "Some Thought", "some-thought"},
		{"mixed case", "Clean Code", "clean-code"},
		{"underscores", "my_great_idea", "my-great-idea"},
		{"special characters", "What's the plan?", "whats-the-plan"},
		{"consecutive spaces", "too   many   spaces", "too-many-spaces"},
		{"leading trailing whitespace", "  padded name  ", "padded-name"},
		{"idempotent", "clean-code", "clean-code"},
		{"numbers preserved", "Chapter 3 Notes", "chapter-3-notes"},
		{"non-ASCII letters preserved", "café latte", "café-latte"},
		{"CJK characters preserved", "我的日記", "我的日記"},
		{"only special chars", "!@#$%", ""},
		{"consecutive hyphens", "a--b---c", "a-b-c"},
		{"mixed separators", "hello_world foo-bar", "hello-world-foo-bar"},
		{"trailing special chars", "hello!", "hello"},
		{"leading special chars", "!hello", "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Slugify(tt.input)
			if got != tt.expected {
				t.Errorf("Slugify(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
