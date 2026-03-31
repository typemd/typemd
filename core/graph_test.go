package core

import "testing"

func TestDotQuote(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`hello`, `"hello"`},
		{`say "hi"`, `"say \"hi\""`},
		{`path\to`, `"path\\to"`},
		{`normal`, `"normal"`},
	}
	for _, tt := range tests {
		got := dotQuote(tt.input)
		if got != tt.want {
			t.Errorf("dotQuote(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestDotID(t *testing.T) {
	got := dotID("book/clean-code-01abc")
	want := `"book/clean-code-01abc"`
	if got != want {
		t.Errorf("dotID = %q, want %q", got, want)
	}
}
