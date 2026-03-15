package cmd

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestStarterPicker_DefaultAllSelected(t *testing.T) {
	m := newStarterPicker()
	for _, item := range m.items {
		if !item.selected {
			t.Errorf("expected %q to be selected by default", item.name)
		}
	}
}

func TestStarterPicker_HasThreeItems(t *testing.T) {
	m := newStarterPicker()
	if len(m.items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(m.items))
	}
}

func keyPress(key string) tea.KeyPressMsg {
	if len(key) == 1 {
		return tea.KeyPressMsg{Code: rune(key[0])}
	}
	switch key {
	case "up":
		return tea.KeyPressMsg{Code: tea.KeyUp}
	case "down":
		return tea.KeyPressMsg{Code: tea.KeyDown}
	case "enter":
		return tea.KeyPressMsg{Code: tea.KeyEnter}
	case "esc":
		return tea.KeyPressMsg{Code: tea.KeyEscape}
	}
	return tea.KeyPressMsg{}
}

func sendKey(m starterPicker, key string) starterPicker {
	result, _ := m.Update(keyPress(key))
	return result.(starterPicker)
}

func TestStarterPicker_MoveDown(t *testing.T) {
	m := newStarterPicker()
	if m.cursor != 0 {
		t.Fatalf("expected cursor at 0, got %d", m.cursor)
	}
	m = sendKey(m, "down")
	if m.cursor != 1 {
		t.Errorf("expected cursor at 1, got %d", m.cursor)
	}
}

func TestStarterPicker_MoveUp(t *testing.T) {
	m := newStarterPicker()
	m = sendKey(m, "down")
	m = sendKey(m, "up")
	if m.cursor != 0 {
		t.Errorf("expected cursor at 0, got %d", m.cursor)
	}
}

func TestStarterPicker_MoveUpAtTop(t *testing.T) {
	m := newStarterPicker()
	m = sendKey(m, "up")
	if m.cursor != 0 {
		t.Errorf("cursor should stay at 0, got %d", m.cursor)
	}
}

func TestStarterPicker_MoveDownAtBottom(t *testing.T) {
	m := newStarterPicker()
	for i := 0; i < 10; i++ {
		m = sendKey(m, "down")
	}
	if m.cursor != len(m.items)-1 {
		t.Errorf("cursor should stay at %d, got %d", len(m.items)-1, m.cursor)
	}
}

func TestStarterPicker_ToggleSpace(t *testing.T) {
	m := newStarterPicker()
	if !m.items[0].selected {
		t.Fatal("expected first item to be selected")
	}
	m = sendKey(m, " ")
	if m.items[0].selected {
		t.Error("expected first item to be deselected after space")
	}
	m = sendKey(m, " ")
	if !m.items[0].selected {
		t.Error("expected first item to be re-selected after second space")
	}
}

func TestStarterPicker_SelectAll(t *testing.T) {
	m := newStarterPicker()
	// Deselect all first
	m = sendKey(m, "n")
	// Now select all
	m = sendKey(m, "a")
	for _, item := range m.items {
		if !item.selected {
			t.Errorf("expected %q to be selected after 'a'", item.name)
		}
	}
}

func TestStarterPicker_DeselectAll(t *testing.T) {
	m := newStarterPicker()
	m = sendKey(m, "n")
	for _, item := range m.items {
		if item.selected {
			t.Errorf("expected %q to be deselected after 'n'", item.name)
		}
	}
}

func TestStarterPicker_EnterConfirms(t *testing.T) {
	m := newStarterPicker()
	m = sendKey(m, "enter")
	if !m.done {
		t.Error("expected done to be true after enter")
	}
	names := m.selectedNames()
	if len(names) != 3 {
		t.Errorf("expected 3 selected names, got %d", len(names))
	}
}

func TestStarterPicker_EscSkips(t *testing.T) {
	m := newStarterPicker()
	m = sendKey(m, "esc")
	if !m.done {
		t.Error("expected done to be true after esc")
	}
	if !m.quitting {
		t.Error("expected quitting to be true after esc")
	}
	names := m.selectedNames()
	if len(names) != 0 {
		t.Errorf("expected 0 selected names after esc, got %d", len(names))
	}
}

func TestStarterPicker_SelectedNames(t *testing.T) {
	m := newStarterPicker()
	// Deselect note (index 1)
	m = sendKey(m, "down")
	m = sendKey(m, " ")
	names := m.selectedNames()
	if len(names) != 2 {
		t.Fatalf("expected 2 selected, got %d: %v", len(names), names)
	}
	// Should have idea and book
	have := make(map[string]bool)
	for _, n := range names {
		have[n] = true
	}
	if !have["idea"] || !have["book"] {
		t.Errorf("expected idea and book, got %v", names)
	}
}

func TestStarterPicker_VJKMovement(t *testing.T) {
	m := newStarterPicker()
	m = sendKey(m, "j")
	if m.cursor != 1 {
		t.Errorf("expected cursor at 1 after j, got %d", m.cursor)
	}
	m = sendKey(m, "k")
	if m.cursor != 0 {
		t.Errorf("expected cursor at 0 after k, got %d", m.cursor)
	}
}
