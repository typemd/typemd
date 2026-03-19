package widget

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestCenteredPopup(t *testing.T) {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)

	result := CenteredPopup("hello", style, 80, 24)
	if result == "" {
		t.Fatal("CenteredPopup returned empty string")
	}
	if !strings.Contains(result, "hello") {
		t.Error("CenteredPopup output does not contain content")
	}
}
