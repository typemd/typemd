package tui

import (
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
)

func pressDateKey(dp *datePicker, key string) (consumed, done, confirmed bool) {
	return dp.Update(tea.KeyPressMsg{Code: -1, Text: key})
}

func pressDateRune(dp *datePicker, r rune) (consumed, done, confirmed bool) {
	return dp.Update(tea.KeyPressMsg{Code: r, Text: string(r)})
}

func TestNewDatePicker_PreFillToday(t *testing.T) {
	dp := newDatePicker(time.Time{})
	today := time.Now()
	if dp.date.Year() != today.Year() || dp.date.Month() != today.Month() || dp.date.Day() != today.Day() {
		t.Errorf("expected today's date, got %s", dp.Value())
	}
	if dp.mode != datePickerSegment {
		t.Error("expected segment mode initially")
	}
	if dp.segment != segYear {
		t.Error("expected focus on year segment initially")
	}
}

func TestParseDatePickerValue(t *testing.T) {
	tv := time.Date(2025, 3, 15, 0, 0, 0, 0, time.Local)
	got := parseDatePickerValue(tv)
	if got.Year() != 2025 || got.Month() != 3 || got.Day() != 15 {
		t.Errorf("time.Time: expected 2025-03-15, got %v", got)
	}

	got = parseDatePickerValue("2024-12-25")
	if got.Year() != 2024 || got.Month() != 12 || got.Day() != 25 {
		t.Errorf("string: expected 2024-12-25, got %v", got)
	}

	got = parseDatePickerValue(nil)
	if !got.IsZero() {
		t.Errorf("nil: expected zero time, got %v", got)
	}

	got = parseDatePickerValue("not-a-date")
	if !got.IsZero() {
		t.Errorf("invalid string: expected zero time, got %v", got)
	}
}

func TestNewDatePicker_PreFillExisting(t *testing.T) {
	v := time.Date(2025, 3, 15, 0, 0, 0, 0, time.Local)
	dp := newDatePicker(v)
	if dp.Value() != "2025-03-15" {
		t.Errorf("expected 2025-03-15, got %s", dp.Value())
	}
}

func TestSegment_Navigation(t *testing.T) {
	dp := newDatePicker(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))

	if dp.segment != segYear {
		t.Error("expected year segment")
	}

	pressDateKey(dp, "right")
	if dp.segment != segMonth {
		t.Error("expected month segment after right")
	}

	pressDateKey(dp, "right")
	if dp.segment != segDay {
		t.Error("expected day segment after right")
	}

	pressDateKey(dp, "right")
	if dp.segment != segDay {
		t.Error("expected day segment (no wrap)")
	}

	pressDateKey(dp, "left")
	if dp.segment != segMonth {
		t.Error("expected month segment after left")
	}

	pressDateKey(dp, "left")
	if dp.segment != segYear {
		t.Error("expected year segment after left")
	}

	pressDateKey(dp, "left")
	if dp.segment != segYear {
		t.Error("expected year segment (no wrap)")
	}
}

func TestSegment_TabNavigation(t *testing.T) {
	dp := newDatePicker(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))

	pressDateKey(dp, "tab")
	if dp.segment != segMonth {
		t.Error("expected month after tab")
	}

	pressDateKey(dp, "shift+tab")
	if dp.segment != segYear {
		t.Error("expected year after shift+tab")
	}
}

func TestSegment_IncrementYear(t *testing.T) {
	dp := newDatePicker(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	dp.segment = segYear

	pressDateKey(dp, "up")
	if dp.date.Year() != 2026 {
		t.Errorf("expected 2026, got %d", dp.date.Year())
	}

	pressDateKey(dp, "down")
	pressDateKey(dp, "down")
	if dp.date.Year() != 2024 {
		t.Errorf("expected 2024, got %d", dp.date.Year())
	}
}

func TestSegment_IncrementMonthWithCarry(t *testing.T) {
	dp := newDatePicker(time.Date(2025, 12, 15, 0, 0, 0, 0, time.Local))
	dp.segment = segMonth

	pressDateKey(dp, "up")
	if dp.date.Month() != time.January || dp.date.Year() != 2026 {
		t.Errorf("expected 2026-01, got %d-%02d", dp.date.Year(), dp.date.Month())
	}

	pressDateKey(dp, "down")
	if dp.date.Month() != time.December || dp.date.Year() != 2025 {
		t.Errorf("expected 2025-12, got %d-%02d", dp.date.Year(), dp.date.Month())
	}
}

func TestSegment_IncrementDayWithCarry(t *testing.T) {
	dp := newDatePicker(time.Date(2025, 1, 31, 0, 0, 0, 0, time.Local))
	dp.segment = segDay

	pressDateKey(dp, "up")
	if dp.date.Month() != time.February || dp.date.Day() != 1 {
		t.Errorf("expected Feb 1, got %s", dp.Value())
	}

	pressDateKey(dp, "down")
	if dp.date.Month() != time.January || dp.date.Day() != 31 {
		t.Errorf("expected Jan 31, got %s", dp.Value())
	}
}

func TestSegment_DayClampOnMonthChange(t *testing.T) {
	dp := newDatePicker(time.Date(2025, 1, 31, 0, 0, 0, 0, time.Local))
	dp.segment = segMonth

	pressDateKey(dp, "up")
	if dp.date.Day() != 28 {
		t.Errorf("expected day clamped to 28, got %d", dp.date.Day())
	}
	if dp.date.Month() != time.February {
		t.Errorf("expected February, got %s", dp.date.Month())
	}
}

func TestSegment_LeapYear(t *testing.T) {
	dp := newDatePicker(time.Date(2024, 2, 29, 0, 0, 0, 0, time.Local))
	if dp.Value() != "2024-02-29" {
		t.Errorf("expected 2024-02-29, got %s", dp.Value())
	}

	dp.segment = segYear
	pressDateKey(dp, "up")
	if dp.date.Day() != 28 {
		t.Errorf("expected day clamped to 28, got %d", dp.date.Day())
	}
}

func TestSegment_DigitInput_Year(t *testing.T) {
	dp := newDatePicker(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	dp.segment = segYear

	pressDateRune(dp, '2')
	pressDateRune(dp, '0')
	pressDateRune(dp, '3')
	pressDateRune(dp, '0')

	if dp.date.Year() != 2030 {
		t.Errorf("expected year 2030, got %d", dp.date.Year())
	}
	if dp.segment != segMonth {
		t.Error("expected auto-advance to month segment")
	}
}

func TestSegment_DigitInput_Month(t *testing.T) {
	dp := newDatePicker(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	dp.segment = segMonth

	pressDateRune(dp, '0')
	pressDateRune(dp, '3')

	if dp.date.Month() != time.March {
		t.Errorf("expected March, got %s", dp.date.Month())
	}
	if dp.segment != segDay {
		t.Error("expected auto-advance to day segment")
	}
}

func TestSegment_DigitInput_MonthSingleDigit(t *testing.T) {
	dp := newDatePicker(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	dp.segment = segMonth

	pressDateRune(dp, '5')

	if dp.date.Month() != time.May {
		t.Errorf("expected May, got %s", dp.date.Month())
	}
	if dp.segment != segDay {
		t.Error("expected auto-advance to day segment")
	}
}

func TestSegment_DigitInput_Day(t *testing.T) {
	dp := newDatePicker(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	dp.segment = segDay

	pressDateRune(dp, '2')
	pressDateRune(dp, '5')

	if dp.date.Day() != 25 {
		t.Errorf("expected day 25, got %d", dp.date.Day())
	}
}

func TestSegment_DigitInput_DayClamped(t *testing.T) {
	dp := newDatePicker(time.Date(2025, 2, 10, 0, 0, 0, 0, time.Local))
	dp.segment = segDay

	pressDateRune(dp, '2')
	pressDateRune(dp, '9')

	if dp.date.Day() != 28 {
		t.Errorf("expected day clamped to 28, got %d", dp.date.Day())
	}
}

func TestSegment_ConfirmCancel(t *testing.T) {
	dp := newDatePicker(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	dp.segment = segYear
	pressDateKey(dp, "up")

	consumed, done, confirmed := pressDateKey(dp, "enter")
	if !consumed || !done || !confirmed {
		t.Error("expected consumed=true, done=true, confirmed=true")
	}
	if dp.Value() != "2026-06-15" {
		t.Errorf("expected 2026-06-15, got %s", dp.Value())
	}
}

func TestSegment_Cancel(t *testing.T) {
	dp := newDatePicker(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	dp.segment = segYear
	pressDateKey(dp, "up")

	consumed, done, confirmed := pressDateKey(dp, "esc")
	if !consumed || !done || confirmed {
		t.Error("expected consumed=true, done=true, confirmed=false")
	}
	if dp.Value() != "2025-06-15" {
		t.Errorf("expected original 2025-06-15, got %s", dp.Value())
	}
}

func TestModeToggle(t *testing.T) {
	dp := newDatePicker(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))

	if dp.mode != datePickerSegment {
		t.Error("expected segment mode initially")
	}

	pressDateKey(dp, "c")
	if dp.mode != datePickerCalendar {
		t.Error("expected calendar mode after c")
	}

	pressDateKey(dp, "c")
	if dp.mode != datePickerSegment {
		t.Error("expected segment mode after second c")
	}
}

func TestCalendar_DayNavigation(t *testing.T) {
	dp := newDatePicker(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	dp.mode = datePickerCalendar

	pressDateKey(dp, "right")
	if dp.date.Day() != 16 {
		t.Errorf("expected day 16, got %d", dp.date.Day())
	}

	pressDateKey(dp, "left")
	if dp.date.Day() != 15 {
		t.Errorf("expected day 15, got %d", dp.date.Day())
	}

	pressDateKey(dp, "down")
	if dp.date.Day() != 22 {
		t.Errorf("expected day 22, got %d", dp.date.Day())
	}

	pressDateKey(dp, "up")
	if dp.date.Day() != 15 {
		t.Errorf("expected day 15, got %d", dp.date.Day())
	}
}

func TestCalendar_VimNavigation(t *testing.T) {
	dp := newDatePicker(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	dp.mode = datePickerCalendar

	pressDateKey(dp, "l")
	if dp.date.Day() != 16 {
		t.Errorf("expected day 16, got %d", dp.date.Day())
	}

	pressDateKey(dp, "h")
	if dp.date.Day() != 15 {
		t.Errorf("expected day 15, got %d", dp.date.Day())
	}

	pressDateKey(dp, "j")
	if dp.date.Day() != 22 {
		t.Errorf("expected day 22, got %d", dp.date.Day())
	}

	pressDateKey(dp, "k")
	if dp.date.Day() != 15 {
		t.Errorf("expected day 15, got %d", dp.date.Day())
	}
}

func TestCalendar_MonthBoundary(t *testing.T) {
	dp := newDatePicker(time.Date(2025, 6, 30, 0, 0, 0, 0, time.Local))
	dp.mode = datePickerCalendar

	pressDateKey(dp, "right")
	if dp.date.Month() != time.July || dp.date.Day() != 1 {
		t.Errorf("expected Jul 1, got %s", dp.Value())
	}

	pressDateKey(dp, "left")
	if dp.date.Month() != time.June || dp.date.Day() != 30 {
		t.Errorf("expected Jun 30, got %s", dp.Value())
	}
}

func TestCalendar_MonthSwitch(t *testing.T) {
	dp := newDatePicker(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	dp.mode = datePickerCalendar

	pressDateKey(dp, "H")
	if dp.date.Month() != time.May {
		t.Errorf("expected May, got %s", dp.date.Month())
	}
	if dp.date.Day() != 15 {
		t.Errorf("expected day 15, got %d", dp.date.Day())
	}

	pressDateKey(dp, "L")
	if dp.date.Month() != time.June {
		t.Errorf("expected June, got %s", dp.date.Month())
	}
}

func TestCalendar_MonthSwitchWithDayClamp(t *testing.T) {
	dp := newDatePicker(time.Date(2025, 3, 31, 0, 0, 0, 0, time.Local))
	dp.mode = datePickerCalendar

	pressDateKey(dp, "H")
	if dp.date.Month() != time.February || dp.date.Day() != 28 {
		t.Errorf("expected Feb 28, got %s", dp.Value())
	}
}

func TestCalendar_JumpToToday(t *testing.T) {
	dp := newDatePicker(time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local))
	dp.mode = datePickerCalendar

	pressDateKey(dp, "t")
	today := time.Now()
	if dp.date.Year() != today.Year() || dp.date.Month() != today.Month() || dp.date.Day() != today.Day() {
		t.Errorf("expected today, got %s", dp.Value())
	}
}

func TestCalendar_ConfirmCancel(t *testing.T) {
	dp := newDatePicker(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	dp.mode = datePickerCalendar
	pressDateKey(dp, "right")

	consumed, done, confirmed := pressDateKey(dp, "enter")
	if !consumed || !done || !confirmed {
		t.Error("expected consumed=true, done=true, confirmed=true")
	}
	if dp.Value() != "2025-06-16" {
		t.Errorf("expected 2025-06-16, got %s", dp.Value())
	}
}

func TestCalendar_CancelRestores(t *testing.T) {
	dp := newDatePicker(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	dp.mode = datePickerCalendar
	pressDateKey(dp, "right")

	_, _, confirmed := pressDateKey(dp, "esc")
	if confirmed {
		t.Error("expected not confirmed on esc")
	}
	if dp.Value() != "2025-06-15" {
		t.Errorf("expected original 2025-06-15, got %s", dp.Value())
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
		{2024, time.February, 29},
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
	dp := newDatePicker(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	view := dp.View()
	if view == "" {
		t.Error("expected non-empty view")
	}
	if !containsStr(view, "Sun") {
		t.Errorf("expected day of week in view, got: %s", view)
	}
}

func TestPickerOverlay_InlineShowsSegments(t *testing.T) {
	dp := newDatePicker(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	dp.mode = datePickerCalendar
	view := dp.View()
	if !containsStr(view, "Sun") {
		t.Errorf("expected day of week in inline view during calendar mode, got: %s", view)
	}
}

func TestPickerOverlay_Output(t *testing.T) {
	dp := newDatePicker(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	dp.mode = datePickerCalendar
	cal := dp.PickerOverlay()
	if !containsStr(cal, "Jun 2025") {
		t.Errorf("expected 'Jun 2025' in picker overlay, got: %s", cal)
	}
	if !containsStr(cal, "Mo Tu We Th Fr Sa Su") {
		t.Errorf("expected weekday headers in picker overlay, got: %s", cal)
	}
}

func TestPickerOverlay_EmptyInSegmentMode(t *testing.T) {
	dp := newDatePicker(time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local))
	if dp.PickerOverlay() != "" {
		t.Error("expected empty overlay in segment mode")
	}
}

// containsStr is defined in type_editor_test.go
