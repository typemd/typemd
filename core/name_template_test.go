package core

import (
	"testing"
	"time"
)

func TestConvertDateFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"YYYY-MM-DD", "2006-01-02"},
		{"YYYY-MM", "2006-01"},
		{"YYYY", "2006"},
		{"MM/DD/YYYY", "01/02/2006"},
		{"HH:mm:ss", "15:04:05"},
		{"YYYY-MM-DD HH:mm", "2006-01-02 15:04"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := convertDateFormat(tt.input)
			if got != tt.expected {
				t.Errorf("convertDateFormat(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestEvaluateNameTemplate(t *testing.T) {
	now := time.Date(2026, 3, 14, 9, 30, 0, 0, time.UTC)

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{"date placeholder", "日記 {{ date:YYYY-MM-DD }}", "日記 2026-03-14"},
		{"year-month only", "月報 {{ date:YYYY-MM }}", "月報 2026-03"},
		{"datetime", "會議 {{ date:YYYY-MM-DD HH:mm }}", "會議 2026-03-14 09:30"},
		{"no placeholders", "Weekly Review", "Weekly Review"},
		{"multiple placeholders", "{{ date:YYYY }}-{{ date:MM }}", "2026-03"},
		{"unknown placeholder type", "{{ foo:bar }}", "{{ foo:bar }}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EvaluateNameTemplate(tt.template, now)
			if got != tt.expected {
				t.Errorf("EvaluateNameTemplate(%q) = %q, want %q", tt.template, got, tt.expected)
			}
		})
	}
}
