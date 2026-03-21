package widget

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestShowSingleItem(t *testing.T) {
	tm := NewToastModel()
	tm.ShowSuccess = true
	cmd := tm.Show(ToastInfo, []ToastItem{{Message: "saved"}})

	if !tm.Active() {
		t.Fatal("expected toast to be active after Show")
	}
	if cmd == nil {
		t.Fatal("expected non-nil cmd (tick) from Show")
	}
	view := tm.View()
	if !strings.Contains(view, "ℹ") {
		t.Errorf("expected info prefix ℹ in view, got: %s", view)
	}
	if !strings.Contains(view, "saved") {
		t.Errorf("expected message 'saved' in view, got: %s", view)
	}
}

func TestShowGroupAggregation(t *testing.T) {
	tm := NewToastModel()
	tm.ShowWarnings = true
	items := []ToastItem{
		{Message: "ref1", Group: "unresolved refs"},
		{Message: "ref2", Group: "unresolved refs"},
		{Message: "ref3", Group: "unresolved refs"},
	}
	tm.Show(ToastWarning, items)

	if !tm.Active() {
		t.Fatal("expected toast to be active")
	}
	view := tm.View()
	if !strings.Contains(view, "3") {
		t.Errorf("expected count '3' in view, got: %s", view)
	}
	if !strings.Contains(view, "unresolved refs") {
		t.Errorf("expected group name 'unresolved refs' in view, got: %s", view)
	}
	if !strings.Contains(view, "⚠") {
		t.Errorf("expected warning prefix ⚠ in view, got: %s", view)
	}
}

func TestShowMixedGroupedAndUngrouped(t *testing.T) {
	tm := NewToastModel()
	tm.ShowWarnings = true
	items := []ToastItem{
		{Message: "orphan found", Group: ""},
		{Message: "ref1", Group: "unresolved refs"},
		{Message: "ref2", Group: "unresolved refs"},
	}
	tm.Show(ToastWarning, items)

	if !tm.Active() {
		t.Fatal("expected toast to be active")
	}
	view := tm.View()
	// Should have the ungrouped message
	if !strings.Contains(view, "orphan found") {
		t.Errorf("expected ungrouped message 'orphan found' in view, got: %s", view)
	}
	// Should have the aggregated group count
	if !strings.Contains(view, "2") {
		t.Errorf("expected count '2' for grouped items in view, got: %s", view)
	}
	if !strings.Contains(view, "unresolved refs") {
		t.Errorf("expected group name 'unresolved refs' in view, got: %s", view)
	}
}

func TestShowInfoWithShowSuccessFalse(t *testing.T) {
	tm := NewToastModel()
	tm.ShowSuccess = false // default
	cmd := tm.Show(ToastInfo, []ToastItem{{Message: "saved"}})

	if tm.Active() {
		t.Fatal("expected toast NOT to be active when ShowSuccess=false for Info level")
	}
	if cmd != nil {
		t.Fatal("expected nil cmd when toast is suppressed")
	}
}

func TestShowWarningWithShowWarningsFalse(t *testing.T) {
	tm := NewToastModel()
	tm.ShowWarnings = false
	cmd := tm.Show(ToastWarning, []ToastItem{{Message: "warning"}})

	if tm.Active() {
		t.Fatal("expected toast NOT to be active when ShowWarnings=false for Warning level")
	}
	if cmd != nil {
		t.Fatal("expected nil cmd when toast is suppressed")
	}
}

func TestShowErrorAlwaysShows(t *testing.T) {
	tm := NewToastModel()
	tm.ShowSuccess = false
	tm.ShowWarnings = false
	cmd := tm.Show(ToastError, []ToastItem{{Message: "something broke"}})

	if !tm.Active() {
		t.Fatal("expected toast to be active for Error level regardless of config")
	}
	if cmd == nil {
		t.Fatal("expected non-nil cmd for Error toast")
	}
	view := tm.View()
	if !strings.Contains(view, "✗") {
		t.Errorf("expected error prefix ✗ in view, got: %s", view)
	}
	if !strings.Contains(view, "something broke") {
		t.Errorf("expected message in view, got: %s", view)
	}
}

func TestUpdateDismissMsgMatchingSeq(t *testing.T) {
	tm := NewToastModel()
	tm.ShowSuccess = true
	tm.Show(ToastInfo, []ToastItem{{Message: "hello"}})

	if !tm.Active() {
		t.Fatal("expected toast to be active before dismiss")
	}

	seq := tm.seq
	tm, _, _ = tm.Update(ToastDismissMsg{Seq: seq})

	if tm.Active() {
		t.Fatal("expected toast to be dismissed after matching seq")
	}
}

func TestUpdateDismissMsgNonMatchingSeq(t *testing.T) {
	tm := NewToastModel()
	tm.ShowSuccess = true
	tm.Show(ToastInfo, []ToastItem{{Message: "hello"}})

	if !tm.Active() {
		t.Fatal("expected toast to be active")
	}

	tm, _, _ = tm.Update(ToastDismissMsg{Seq: tm.seq + 100})

	if !tm.Active() {
		t.Fatal("expected toast to remain active with non-matching seq")
	}
}

func TestUpdateDismissKeyPress(t *testing.T) {
	tm := NewToastModel()
	tm.ShowSuccess = true
	tm.Show(ToastInfo, []ToastItem{{Message: "hello"}})

	if !tm.Active() {
		t.Fatal("expected toast to be active")
	}

	msg := tea.KeyPressMsg{Code: tea.KeyEscape}
	var consumed bool
	tm, _, consumed = tm.Update(msg)

	if tm.Active() {
		t.Fatal("expected toast to be dismissed after Esc key")
	}
	if !consumed {
		t.Fatal("expected consumed=true when dismiss key is pressed on active toast")
	}
}

func TestUpdateOtherKeyPressWhenActive(t *testing.T) {
	tm := NewToastModel()
	tm.ShowSuccess = true
	tm.Show(ToastInfo, []ToastItem{{Message: "hello"}})

	msg := tea.KeyPressMsg{Code: 'a', Text: "a"}
	var consumed bool
	tm, _, consumed = tm.Update(msg)

	if consumed {
		t.Fatal("expected consumed=false for non-dismiss key when active")
	}
	if !tm.Active() {
		t.Fatal("expected toast to remain active after non-dismiss key")
	}
}

func TestUpdateDismissKeyWhenNotActive(t *testing.T) {
	tm := NewToastModel()

	msg := tea.KeyPressMsg{Code: tea.KeyEscape}
	var consumed bool
	tm, _, consumed = tm.Update(msg)

	if consumed {
		t.Fatal("expected consumed=false when toast is not active")
	}
}

func TestViewWhenNotActive(t *testing.T) {
	tm := NewToastModel()
	view := tm.View()
	if view != "" {
		t.Errorf("expected empty string when not active, got: %q", view)
	}
}

func TestViewWhenActive(t *testing.T) {
	tm := NewToastModel()
	tm.ShowSuccess = true
	tm.Show(ToastInfo, []ToastItem{{Message: "done"}})

	view := tm.View()
	if view == "" {
		t.Fatal("expected non-empty view when active")
	}
	if !strings.Contains(view, "ℹ") {
		t.Errorf("expected info prefix in view, got: %s", view)
	}
	if !strings.Contains(view, "done") {
		t.Errorf("expected message 'done' in view, got: %s", view)
	}
}

func TestOverlayWhenNotActive(t *testing.T) {
	tm := NewToastModel()
	bg := "background content here"
	result := tm.Overlay(bg, 80, 24)
	if result != bg {
		t.Errorf("expected background unchanged when not active, got: %q", result)
	}
}

func TestShowResetsPreviousToast(t *testing.T) {
	tm := NewToastModel()
	tm.ShowWarnings = true
	tm.ShowSuccess = true

	// Show first toast
	tm.Show(ToastWarning, []ToastItem{{Message: "first"}})
	if !tm.Active() {
		t.Fatal("expected first toast to be active")
	}
	firstSeq := tm.seq

	// Show second toast (should replace)
	tm.Show(ToastInfo, []ToastItem{{Message: "second"}})
	if !tm.Active() {
		t.Fatal("expected second toast to be active")
	}
	secondSeq := tm.seq

	if secondSeq <= firstSeq {
		t.Errorf("expected seq to increment, got first=%d second=%d", firstSeq, secondSeq)
	}

	// Old dismiss msg should NOT dismiss new toast
	tm, _, _ = tm.Update(ToastDismissMsg{Seq: firstSeq})
	if !tm.Active() {
		t.Fatal("expected toast to remain active after stale dismiss msg")
	}

	view := tm.View()
	if !strings.Contains(view, "second") {
		t.Errorf("expected second message in view, got: %s", view)
	}
	if strings.Contains(view, "first") {
		t.Errorf("expected first message NOT in view, got: %s", view)
	}
}
