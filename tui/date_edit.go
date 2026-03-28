package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type dateEditMode int

const (
	dateSegmentMode  dateEditMode = iota
	dateCalendarMode
)

type dateSegment int

const (
	segYear  dateSegment = iota
	segMonth
	segDay
)

// Styles for date picker rendering (package-level to avoid per-render allocation).
var (
	dateSegmentFocusStyle = lipgloss.NewStyle().Bold(true).Reverse(true)
	dateCalCursorStyle    = lipgloss.NewStyle().Bold(true).Reverse(true)
	dateCalTodayStyle     = lipgloss.NewStyle().Bold(true).Underline(true)
)

// dateEdit is a shared date editing widget used by both propEditor and cellEdit.
type dateEdit struct {
	mode     dateEditMode
	original time.Time
	date     time.Time

	// Segment mode state
	segment    dateSegment
	digitBuf   string
	digitCount int
}

// localToday returns today's date at midnight in local time.
func localToday() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
}

// parseDateValue extracts a time.Time from a property value that may be
// time.Time, string ("2006-01-02"), or nil.
func parseDateValue(v any) time.Time {
	switch val := v.(type) {
	case time.Time:
		return val
	case string:
		if t, err := time.Parse("2006-01-02", val); err == nil {
			return t
		}
	}
	return time.Time{}
}

// newDateEdit creates a date picker initialized with the given value.
// If value is zero, today's date is used.
func newDateEdit(value time.Time) *dateEdit {
	if value.IsZero() {
		value = localToday()
	}
	value = time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.Local)
	return &dateEdit{
		mode:     dateSegmentMode,
		original: value,
		date:     value,
		segment:  segYear,
	}
}

// Value returns the current date formatted as YYYY-MM-DD.
func (de *dateEdit) Value() string {
	return de.date.Format("2006-01-02")
}

// Mode returns the current edit mode.
func (de *dateEdit) Mode() dateEditMode {
	return de.mode
}

// Update handles a key press and returns whether the date picker consumed it,
// and whether editing is done (confirmed or cancelled).
// Returns: (consumed bool, done bool, confirmed bool)
func (de *dateEdit) Update(msg tea.KeyPressMsg) (bool, bool, bool) {
	switch de.mode {
	case dateSegmentMode:
		return de.updateSegment(msg)
	case dateCalendarMode:
		return de.updateCalendar(msg)
	}
	return false, false, false
}

func (de *dateEdit) updateSegment(msg tea.KeyPressMsg) (bool, bool, bool) {
	key := msg.String()

	switch key {
	case "enter":
		return true, true, true
	case "esc":
		de.date = de.original
		return true, true, false
	case "c":
		de.mode = dateCalendarMode
		de.clearDigitBuf()
		return true, false, false
	case "left", "shift+tab":
		de.clearDigitBuf()
		if de.segment > segYear {
			de.segment--
		}
		return true, false, false
	case "right", "tab":
		de.clearDigitBuf()
		if de.segment < segDay {
			de.segment++
		}
		return true, false, false
	case "up":
		de.clearDigitBuf()
		de.incrementSegment(1)
		return true, false, false
	case "down":
		de.clearDigitBuf()
		de.incrementSegment(-1)
		return true, false, false
	default:
		if len(key) == 1 && key[0] >= '0' && key[0] <= '9' {
			de.handleDigit(key[0])
			return true, false, false
		}
	}
	return false, false, false
}

func (de *dateEdit) updateCalendar(msg tea.KeyPressMsg) (bool, bool, bool) {
	key := msg.String()

	switch key {
	case "enter":
		return true, true, true
	case "esc":
		de.date = de.original
		return true, true, false
	case "c":
		de.mode = dateSegmentMode
		de.segment = segYear
		de.clearDigitBuf()
		return true, false, false
	case "left", "h":
		de.date = de.date.AddDate(0, 0, -1)
		return true, false, false
	case "right", "l":
		de.date = de.date.AddDate(0, 0, 1)
		return true, false, false
	case "up", "k":
		de.date = de.date.AddDate(0, 0, -7)
		return true, false, false
	case "down", "j":
		de.date = de.date.AddDate(0, 0, 7)
		return true, false, false
	case "H":
		de.prevMonth()
		return true, false, false
	case "L":
		de.nextMonth()
		return true, false, false
	case "t":
		de.date = localToday()
		return true, false, false
	}
	return false, false, false
}

func (de *dateEdit) incrementSegment(delta int) {
	y, m, d := de.date.Year(), de.date.Month(), de.date.Day()

	switch de.segment {
	case segYear:
		y += delta
		if y < 1 {
			y = 1
		}
		if y > 9999 {
			y = 9999
		}
	case segMonth:
		m += time.Month(delta)
		if m > 12 {
			m = 1
			y++
		} else if m < 1 {
			m = 12
			y--
		}
	case segDay:
		de.date = de.date.AddDate(0, 0, delta)
		return
	}

	maxDay := daysInMonth(y, m)
	if d > maxDay {
		d = maxDay
	}
	de.date = time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

func (de *dateEdit) handleDigit(digit byte) {
	de.digitBuf += string(digit)
	de.digitCount++

	y, m, d := de.date.Year(), de.date.Month(), de.date.Day()

	switch de.segment {
	case segYear:
		if de.digitCount >= 4 {
			val, _ := strconv.Atoi(de.digitBuf)
			if val < 1 {
				val = 1
			}
			if val > 9999 {
				val = 9999
			}
			y = val
			de.clearDigitBuf()
			de.segment = segMonth
		} else {
			return
		}
	case segMonth:
		val, _ := strconv.Atoi(de.digitBuf)
		if de.digitCount >= 2 || val > 1 {
			if val < 1 {
				val = 1
			}
			if val > 12 {
				val = 12
			}
			m = time.Month(val)
			de.clearDigitBuf()
			de.segment = segDay
		} else {
			return
		}
	case segDay:
		val, _ := strconv.Atoi(de.digitBuf)
		maxDay := daysInMonth(y, m)
		if de.digitCount >= 2 || val > 3 {
			if val < 1 {
				val = 1
			}
			if val > maxDay {
				val = maxDay
			}
			d = val
			de.clearDigitBuf()
		} else {
			return
		}
	}

	maxDay := daysInMonth(y, m)
	if d > maxDay {
		d = maxDay
	}
	de.date = time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

func (de *dateEdit) prevMonth() {
	y, m, d := de.date.Year(), de.date.Month(), de.date.Day()
	m--
	if m < 1 {
		m = 12
		y--
	}
	maxDay := daysInMonth(y, m)
	if d > maxDay {
		d = maxDay
	}
	de.date = time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

func (de *dateEdit) nextMonth() {
	y, m, d := de.date.Year(), de.date.Month(), de.date.Day()
	m++
	if m > 12 {
		m = 1
		y++
	}
	maxDay := daysInMonth(y, m)
	if d > maxDay {
		d = maxDay
	}
	de.date = time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

func (de *dateEdit) clearDigitBuf() {
	de.digitBuf = ""
	de.digitCount = 0
}

// View returns the inline portion of the date picker.
// In both modes, this shows the segmented date display.
// In calendar mode, the calendar grid is rendered separately via CalendarOverlay().
func (de *dateEdit) View() string {
	return de.viewSegment()
}

// CalendarOverlay returns the calendar grid content for overlay rendering.
// Returns empty string when not in calendar mode.
func (de *dateEdit) CalendarOverlay() string {
	if de.mode != dateCalendarMode {
		return ""
	}
	return de.viewCalendar()
}

func (de *dateEdit) viewSegment() string {
	y, m, d := de.date.Year(), de.date.Month(), de.date.Day()
	dow := de.date.Weekday().String()[:3]

	yearStr := fmt.Sprintf("%04d", y)
	monthStr := fmt.Sprintf("%02d", int(m))
	dayStr := fmt.Sprintf("%02d", d)

	switch de.segment {
	case segYear:
		if de.digitCount > 0 {
			yearStr = de.digitBuf + strings.Repeat("_", 4-de.digitCount)
		}
	case segMonth:
		if de.digitCount > 0 {
			monthStr = de.digitBuf + strings.Repeat("_", 2-de.digitCount)
		}
	case segDay:
		if de.digitCount > 0 {
			dayStr = de.digitBuf + strings.Repeat("_", 2-de.digitCount)
		}
	}

	var parts [3]string
	parts[0] = yearStr
	parts[1] = monthStr
	parts[2] = dayStr

	switch de.segment {
	case segYear:
		parts[0] = dateSegmentFocusStyle.Render(yearStr)
	case segMonth:
		parts[1] = dateSegmentFocusStyle.Render(monthStr)
	case segDay:
		parts[2] = dateSegmentFocusStyle.Render(dayStr)
	}

	return fmt.Sprintf("%s-%s-%s  %s", parts[0], parts[1], parts[2], dow)
}

func (de *dateEdit) viewCalendar() string {
	y, m := de.date.Year(), de.date.Month()
	today := localToday()

	// Grid width: "Mo Tu We Th Fr Sa Su" = 20 chars
	const gridWidth = 20

	var b strings.Builder

	// Month header — centered over the grid
	header := fmt.Sprintf("%s %d", m.String()[:3], y)
	pad := (gridWidth - len(header)) / 2
	if pad < 0 {
		pad = 0
	}
	b.WriteString(strings.Repeat(" ", pad) + header + "\n")
	b.WriteString("Mo Tu We Th Fr Sa Su\n")

	first := time.Date(y, m, 1, 0, 0, 0, 0, time.Local)
	offset := int(first.Weekday())
	if offset == 0 {
		offset = 7
	}
	offset--

	lastDay := daysInMonth(y, m)

	for i := 0; i < offset; i++ {
		b.WriteString("   ")
	}

	for day := 1; day <= lastDay; day++ {
		dayDate := time.Date(y, m, day, 0, 0, 0, 0, time.Local)
		dayStr := fmt.Sprintf("%2d", day)

		isCursor := de.date.Year() == y && de.date.Month() == m && de.date.Day() == day
		isToday := dayDate.Equal(today)

		if isCursor {
			b.WriteString(dateCalCursorStyle.Render(dayStr))
		} else if isToday {
			b.WriteString(dateCalTodayStyle.Render(dayStr))
		} else {
			b.WriteString(dayStr)
		}

		weekday := int(dayDate.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		if weekday == 7 && day < lastDay {
			b.WriteString("\n")
		} else if day < lastDay {
			b.WriteString(" ")
		}
	}

	return b.String()
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.Local).Day()
}
