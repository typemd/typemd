package tui

import (
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
)

func pressDateKey(de *dateEdit, key string) (consumed, done, confirmed bool) {
	return de.Update(tea.KeyPressMsg{Code: -1, Text: key})
}

func pressDateRune(de *dateEdit, r rune) (consumed, done, confirmed bool) {
	return de.Update(tea.KeyPressMsg{Code: r, Text: string(r)})
}

func TestNewDateEdit_PreFillToday(t *testing.T) {
	de := newDateEdit(time.Time{})
	today := time.Now()
	if de.date.Year() != today.Year() || de.date.Month() != today.Month() || de.date.Day() != today.Day() {
		t.Errorf("expected today's date, got %s", de.Value())
	}
	if de.mode != dateSegmentMode {
		t.Error("expected segment mode initially")
	}
	if de.segment != segYear {
		t.Error("expected focus on year segment initially")
	}
}

func TestParseDateValue(t *testing.T) {
	// time.Time
	tv := time.Date(2025, 3, 15, 0, 0, 0, 0, time.Local)
	got := parseDateValue(tv)
	if got.Year() != 2025 || got.Month() != 3 || got.Day() != 15 {
		t.Errorf("time.Time: expected 2025-03-15, got %v", got)
	}

	// string
	got = parseDateValue("2024-12-25")
	if got.Year() != 2024 || got.Month() != 12 || got.Day() != 25 {
		t.Errorf("string: expected 2024-12-25, got %v", got)
	}

	// nil → zero
	got = parseDateValue(nil)
	if !got.IsZero() {
		t.Errorf("nil: expected zero time, got %v", got)
	}

	// invalid string → zero
	got = parseDateValue("not-a-date")
	if !got.IsZero() {
		t.Errorf("invalid string: expected zero time, got %v", got)
	}
}

func TestNewDateEdit_PreFillExisting(t *testing.T) {
	v := time.Date(2025, 3, 15, 0, 0, 0, 0, time.Local)
	de := newDateEdit(v)
	if de.Value() != "2025-03-15" {
		t.Errorf("expected 2025-03-15, got %s", de.Value())
	}
}

func TestSegment_Navigation(t *testing.T) {
	de := newDateEdit(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))

	// Start on year
	if de.segment != segYear {
		t.Error("expected year segment")
	}

	// Right → month
	pressDateKey(de, "right")
	if de.segment != segMonth {
		t.Error("expected month segment after right")
	}

	// Right → day
	pressDateKey(de, "right")
	if de.segment != segDay {
		t.Error("expected day segment after right")
	}

	// Right at day → stays on day (no wrap)
	pressDateKey(de, "right")
	if de.segment != segDay {
		t.Error("expected day segment (no wrap)")
	}

	// Left → month
	pressDateKey(de, "left")
	if de.segment != segMonth {
		t.Error("expected month segment after left")
	}

	// Left → year
	pressDateKey(de, "left")
	if de.segment != segYear {
		t.Error("expected year segment after left")
	}

	// Left at year → stays on year
	pressDateKey(de, "left")
	if de.segment != segYear {
		t.Error("expected year segment (no wrap)")
	}
}

func TestSegment_TabNavigation(t *testing.T) {
	de := newDateEdit(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))

	pressDateKey(de, "tab")
	if de.segment != segMonth {
		t.Error("expected month after tab")
	}

	pressDateKey(de, "shift+tab")
	if de.segment != segYear {
		t.Error("expected year after shift+tab")
	}
}

func TestSegment_IncrementYear(t *testing.T) {
	de := newDateEdit(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	de.segment = segYear

	pressDateKey(de, "up")
	if de.date.Year() != 2026 {
		t.Errorf("expected 2026, got %d", de.date.Year())
	}

	pressDateKey(de, "down")
	pressDateKey(de, "down")
	if de.date.Year() != 2024 {
		t.Errorf("expected 2024, got %d", de.date.Year())
	}
}

func TestSegment_IncrementMonthWithCarry(t *testing.T) {
	de := newDateEdit(time.Date(2025, 12, 15, 0, 0, 0, 0, time.Local))
	de.segment = segMonth

	// Dec + 1 → Jan next year
	pressDateKey(de, "up")
	if de.date.Month() != time.January || de.date.Year() != 2026 {
		t.Errorf("expected 2026-01, got %d-%02d", de.date.Year(), de.date.Month())
	}

	// Jan - 1 → Dec previous year
	pressDateKey(de, "down")
	if de.date.Month() != time.December || de.date.Year() != 2025 {
		t.Errorf("expected 2025-12, got %d-%02d", de.date.Year(), de.date.Month())
	}
}

func TestSegment_IncrementDayWithCarry(t *testing.T) {
	// Jan 31 + 1 day
	de := newDateEdit(time.Date(2025, 1, 31, 0, 0, 0, 0, time.Local))
	de.segment = segDay

	pressDateKey(de, "up")
	if de.date.Month() != time.February || de.date.Day() != 1 {
		t.Errorf("expected Feb 1, got %s", de.Value())
	}

	// Feb 1 - 1 day → Jan 31
	pressDateKey(de, "down")
	if de.date.Month() != time.January || de.date.Day() != 31 {
		t.Errorf("expected Jan 31, got %s", de.Value())
	}
}

func TestSegment_DayClampOnMonthChange(t *testing.T) {
	// Jan 31, change month to Feb → clamp to 28
	de := newDateEdit(time.Date(2025, 1, 31, 0, 0, 0, 0, time.Local))
	de.segment = segMonth

	pressDateKey(de, "up")
	if de.date.Day() != 28 {
		t.Errorf("expected day clamped to 28, got %d", de.date.Day())
	}
	if de.date.Month() != time.February {
		t.Errorf("expected February, got %s", de.date.Month())
	}
}

func TestSegment_LeapYear(t *testing.T) {
	// 2024 is leap year, Feb 29 exists
	de := newDateEdit(time.Date(2024, 2, 29, 0, 0, 0, 0, time.Local))
	if de.Value() != "2024-02-29" {
		t.Errorf("expected 2024-02-29, got %s", de.Value())
	}

	// Change year to 2025 (non-leap) → clamp to Feb 28
	de.segment = segYear
	pressDateKey(de, "up")
	if de.date.Day() != 28 {
		t.Errorf("expected day clamped to 28, got %d", de.date.Day())
	}
}

func TestSegment_DigitInput_Year(t *testing.T) {
	de := newDateEdit(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	de.segment = segYear

	pressDateRune(de, '2')
	pressDateRune(de, '0')
	pressDateRune(de, '3')
	pressDateRune(de, '0')

	if de.date.Year() != 2030 {
		t.Errorf("expected year 2030, got %d", de.date.Year())
	}
	// Should auto-advance to month
	if de.segment != segMonth {
		t.Error("expected auto-advance to month segment")
	}
}

func TestSegment_DigitInput_Month(t *testing.T) {
	de := newDateEdit(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	de.segment = segMonth

	pressDateRune(de, '0')
	pressDateRune(de, '3')

	if de.date.Month() != time.March {
		t.Errorf("expected March, got %s", de.date.Month())
	}
	if de.segment != segDay {
		t.Error("expected auto-advance to day segment")
	}
}

func TestSegment_DigitInput_MonthSingleDigit(t *testing.T) {
	de := newDateEdit(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	de.segment = segMonth

	// Typing "5" should auto-advance (5 > 1, so can't be first digit of 10-12)
	pressDateRune(de, '5')

	if de.date.Month() != time.May {
		t.Errorf("expected May, got %s", de.date.Month())
	}
	if de.segment != segDay {
		t.Error("expected auto-advance to day segment")
	}
}

func TestSegment_DigitInput_Day(t *testing.T) {
	de := newDateEdit(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	de.segment = segDay

	pressDateRune(de, '2')
	pressDateRune(de, '5')

	if de.date.Day() != 25 {
		t.Errorf("expected day 25, got %d", de.date.Day())
	}
}

func TestSegment_DigitInput_DayClamped(t *testing.T) {
	// Feb in non-leap year, type 29 → clamped to 28
	de := newDateEdit(time.Date(2025, 2, 10, 0, 0, 0, 0, time.Local))
	de.segment = segDay

	pressDateRune(de, '2')
	pressDateRune(de, '9')

	if de.date.Day() != 28 {
		t.Errorf("expected day clamped to 28, got %d", de.date.Day())
	}
}

func TestSegment_ConfirmCancel(t *testing.T) {
	de := newDateEdit(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	de.segment = segYear
	pressDateKey(de, "up") // year → 2026

	// Confirm
	consumed, done, confirmed := pressDateKey(de, "enter")
	if !consumed || !done || !confirmed {
		t.Error("expected consumed=true, done=true, confirmed=true")
	}
	if de.Value() != "2026-06-15" {
		t.Errorf("expected 2026-06-15, got %s", de.Value())
	}
}

func TestSegment_Cancel(t *testing.T) {
	de := newDateEdit(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	de.segment = segYear
	pressDateKey(de, "up") // year → 2026

	consumed, done, confirmed := pressDateKey(de, "esc")
	if !consumed || !done || confirmed {
		t.Error("expected consumed=true, done=true, confirmed=false")
	}
	if de.Value() != "2025-06-15" {
		t.Errorf("expected original 2025-06-15, got %s", de.Value())
	}
}

func TestModeToggle(t *testing.T) {
	de := newDateEdit(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))

	if de.mode != dateSegmentMode {
		t.Error("expected segment mode initially")
	}

	pressDateKey(de, "c")
	if de.mode != dateCalendarMode {
		t.Error("expected calendar mode after c")
	}

	pressDateKey(de, "c")
	if de.mode != dateSegmentMode {
		t.Error("expected segment mode after second c")
	}
}

func TestCalendar_DayNavigation(t *testing.T) {
	de := newDateEdit(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	de.mode = dateCalendarMode

	pressDateKey(de, "right")
	if de.date.Day() != 16 {
		t.Errorf("expected day 16, got %d", de.date.Day())
	}

	pressDateKey(de, "left")
	if de.date.Day() != 15 {
		t.Errorf("expected day 15, got %d", de.date.Day())
	}

	pressDateKey(de, "down") // +7 days
	if de.date.Day() != 22 {
		t.Errorf("expected day 22, got %d", de.date.Day())
	}

	pressDateKey(de, "up") // -7 days
	if de.date.Day() != 15 {
		t.Errorf("expected day 15, got %d", de.date.Day())
	}
}

func TestCalendar_VimNavigation(t *testing.T) {
	de := newDateEdit(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	de.mode = dateCalendarMode

	pressDateKey(de, "l")
	if de.date.Day() != 16 {
		t.Errorf("expected day 16, got %d", de.date.Day())
	}

	pressDateKey(de, "h")
	if de.date.Day() != 15 {
		t.Errorf("expected day 15, got %d", de.date.Day())
	}

	pressDateKey(de, "j")
	if de.date.Day() != 22 {
		t.Errorf("expected day 22, got %d", de.date.Day())
	}

	pressDateKey(de, "k")
	if de.date.Day() != 15 {
		t.Errorf("expected day 15, got %d", de.date.Day())
	}
}

func TestCalendar_MonthBoundary(t *testing.T) {
	// June 30 + right → July 1
	de := newDateEdit(time.Date(2025, 6, 30, 0, 0, 0, 0, time.Local))
	de.mode = dateCalendarMode

	pressDateKey(de, "right")
	if de.date.Month() != time.July || de.date.Day() != 1 {
		t.Errorf("expected Jul 1, got %s", de.Value())
	}

	// July 1 - left → June 30
	pressDateKey(de, "left")
	if de.date.Month() != time.June || de.date.Day() != 30 {
		t.Errorf("expected Jun 30, got %s", de.Value())
	}
}

func TestCalendar_MonthSwitch(t *testing.T) {
	de := newDateEdit(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	de.mode = dateCalendarMode

	pressDateKey(de, "H") // prev month
	if de.date.Month() != time.May {
		t.Errorf("expected May, got %s", de.date.Month())
	}
	if de.date.Day() != 15 {
		t.Errorf("expected day 15, got %d", de.date.Day())
	}

	pressDateKey(de, "L") // next month
	if de.date.Month() != time.June {
		t.Errorf("expected June, got %s", de.date.Month())
	}
}

func TestCalendar_MonthSwitchWithDayClamp(t *testing.T) {
	// Jan 31, H → prev month clamps day
	de := newDateEdit(time.Date(2025, 3, 31, 0, 0, 0, 0, time.Local))
	de.mode = dateCalendarMode

	pressDateKey(de, "H") // Feb has 28 days
	if de.date.Month() != time.February || de.date.Day() != 28 {
		t.Errorf("expected Feb 28, got %s", de.Value())
	}
}

func TestCalendar_JumpToToday(t *testing.T) {
	de := newDateEdit(time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local))
	de.mode = dateCalendarMode

	pressDateKey(de, "t")
	today := time.Now()
	if de.date.Year() != today.Year() || de.date.Month() != today.Month() || de.date.Day() != today.Day() {
		t.Errorf("expected today, got %s", de.Value())
	}
}

func TestCalendar_ConfirmCancel(t *testing.T) {
	de := newDateEdit(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	de.mode = dateCalendarMode
	pressDateKey(de, "right") // → Jun 16

	consumed, done, confirmed := pressDateKey(de, "enter")
	if !consumed || !done || !confirmed {
		t.Error("expected consumed=true, done=true, confirmed=true")
	}
	if de.Value() != "2025-06-16" {
		t.Errorf("expected 2025-06-16, got %s", de.Value())
	}
}

func TestCalendar_CancelRestores(t *testing.T) {
	de := newDateEdit(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	de.mode = dateCalendarMode
	pressDateKey(de, "right") // → Jun 16

	_, _, confirmed := pressDateKey(de, "esc")
	if confirmed {
		t.Error("expected not confirmed on esc")
	}
	if de.Value() != "2025-06-15" {
		t.Errorf("expected original 2025-06-15, got %s", de.Value())
	}
}

func TestDaysInMonth(t *testing.T) {
	tests := []struct {
		year  int
		month time.Month
		want  int
	}{
		{2025, time.January, 31},
		{2025, time.February, 28},
		{2024, time.February, 29}, // leap year
		{2025, time.April, 30},
		{2025, time.December, 31},
	}
	for _, tt := range tests {
		got := daysInMonth(tt.year, tt.month)
		if got != tt.want {
			t.Errorf("daysInMonth(%d, %s) = %d, want %d", tt.year, tt.month, got, tt.want)
		}
	}
}

func TestViewSegment_Output(t *testing.T) {
	de := newDateEdit(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	view := de.View()
	// Should contain the date parts and day of week
	if view == "" {
		t.Error("expected non-empty view")
	}
	// Should contain "Sun" (June 15 2025 is Sunday)
	if !containsStr(view, "Sun") {
		t.Errorf("expected day of week in view, got: %s", view)
	}
}

func TestViewCalendar_Output(t *testing.T) {
	de := newDateEdit(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	de.mode = dateCalendarMode
	view := de.View()
	// Should contain month header and weekday headers
	if !containsStr(view, "Jun 2025") {
		t.Errorf("expected 'Jun 2025' in calendar view, got: %s", view)
	}
	if !containsStr(view, "Mo Tu We Th Fr Sa Su") {
		t.Errorf("expected weekday headers in calendar view, got: %s", view)
	}
}

// containsStr is defined in type_editor_test.go
