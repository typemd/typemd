package widget

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ToastLevel represents the severity level of a toast notification.
type ToastLevel int

const (
	ToastInfo    ToastLevel = iota
	ToastWarning
	ToastError
)

// ToastItem represents a single item to display in a toast notification.
type ToastItem struct {
	Message string
	Group   string // messages with same group are counted, not listed
}

// ToastDismissMsg is sent by tea.Tick to auto-dismiss. Seq prevents stale ticks
// from dismissing newer toasts.
type ToastDismissMsg struct {
	Seq int
}

// ToastModel manages toast notification display state.
type ToastModel struct {
	// Config
	DurationMs   int    // default 3000
	DismissKey   string // default "esc"
	ShowWarnings bool   // default true
	ShowSuccess  bool   // default false

	// State
	active bool
	level  ToastLevel
	lines  []string // rendered message lines
	seq    int      // monotonic counter for dismiss correlation
}

// NewToastModel creates a ToastModel with default configuration.
func NewToastModel() ToastModel {
	return ToastModel{
		DurationMs:   3000,
		DismissKey:   "esc",
		ShowWarnings: true,
		ShowSuccess:  false,
	}
}

// Active returns whether a toast is currently visible.
func (t *ToastModel) Active() bool {
	return t.active
}

func levelPrefix(level ToastLevel) string {
	switch level {
	case ToastInfo:
		return "ℹ"
	case ToastWarning:
		return "⚠"
	case ToastError:
		return "✗"
	default:
		return "ℹ"
	}
}

// Show activates a toast with the given level and items. It processes items
// (aggregating by group), increments the seq counter, and returns a tea.Tick
// cmd for auto-dismiss. Returns nil if the toast is suppressed by config.
func (t *ToastModel) Show(level ToastLevel, items []ToastItem) tea.Cmd {
	// Check suppression config
	if level == ToastInfo && !t.ShowSuccess {
		return nil
	}
	if level == ToastWarning && !t.ShowWarnings {
		return nil
	}

	prefix := levelPrefix(level)

	// Aggregate items by group
	groupCounts := make(map[string]int)
	groupOrder := []string{}
	var ungrouped []string

	for _, item := range items {
		if item.Group == "" {
			ungrouped = append(ungrouped, item.Message)
		} else {
			if groupCounts[item.Group] == 0 {
				groupOrder = append(groupOrder, item.Group)
			}
			groupCounts[item.Group]++
		}
	}

	// Build rendered lines
	var lines []string

	// Ungrouped items: each rendered individually
	for _, msg := range ungrouped {
		lines = append(lines, fmt.Sprintf("%s %s", prefix, msg))
	}

	// Grouped items: aggregated count
	for _, group := range groupOrder {
		count := groupCounts[group]
		lines = append(lines, fmt.Sprintf("%s %d %s", prefix, count, group))
	}

	// Update state
	t.active = true
	t.level = level
	t.lines = lines
	t.seq++

	// Return tick cmd for auto-dismiss
	duration := time.Duration(t.DurationMs) * time.Millisecond
	seq := t.seq
	return tea.Tick(duration, func(time.Time) tea.Msg {
		return ToastDismissMsg{Seq: seq}
	})
}

// Update handles ToastDismissMsg and tea.KeyPressMsg. Returns the updated model,
// an optional cmd, and whether the message was consumed by the toast.
func (t ToastModel) Update(msg tea.Msg) (ToastModel, tea.Cmd, bool) {
	switch msg := msg.(type) {
	case ToastDismissMsg:
		if t.active && msg.Seq == t.seq {
			t.active = false
			t.lines = nil
		}
		return t, nil, false

	case tea.KeyPressMsg:
		if !t.active {
			return t, nil, false
		}
		if msg.String() == t.DismissKey {
			t.active = false
			t.lines = nil
			return t, nil, true
		}
		return t, nil, false
	}

	return t, nil, false
}

// View returns the rendered toast content. Returns an empty string if not active.
func (t *ToastModel) View() string {
	if !t.active {
		return ""
	}

	content := strings.Join(t.lines, "\n")

	var colorCode string
	switch t.level {
	case ToastInfo:
		colorCode = "2" // green
	case ToastWarning:
		colorCode = "3" // yellow
	case ToastError:
		colorCode = "1" // red
	}

	fg := lipgloss.Color(colorCode)
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(fg).
		Foreground(fg).
		Padding(0, 1)

	return style.Render(content)
}

// Overlay composites the toast onto a background positioned in the bottom-right corner.
// Returns background unchanged if not active.
func (t *ToastModel) Overlay(background string, termW, termH int) string {
	popup := t.View()
	if popup == "" {
		return background
	}

	popupW := lipgloss.Width(popup)
	popupH := lipgloss.Height(popup)

	x := termW - popupW - 2
	y := termH - popupH - 2

	return OverlayAt(background, popup, x, y, termW, termH)
}
