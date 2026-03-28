package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type datePickerMode int

const (
	datePickerSegment  datePickerMode = iota
	datePickerCalendar
)

type dpSegment int

const (
	segYear  dpSegment = iota
	segMonth
	segDay
)

var (
	dpSegmentFocusStyle = lipgloss.NewStyle().Bold(true).Reverse(true)
	dpCursorStyle       = lipgloss.NewStyle().Bold(true).Reverse(true)
	dpTodayStyle        = lipgloss.NewStyle().Bold(true).Underline(true)
)

const (
	dpGridWidth  = 20 // "Mo Tu We Th Fr Sa Su"
	dpGridHeight = 8  // header + weekday header + 6 week rows
)

// datePicker is a shared date editing widget used by both propEditor and cellEdit.
type datePicker struct {
	mode     datePickerMode
	original time.Time
	date     time.Time

	segment    dpSegment
	digitBuf   string
	digitCount int
}

func localToday() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
}

// parseDatePickerValue extracts a time.Time from a property value that may be
// time.Time, string ("2006-01-02"), or nil.
func parseDatePickerValue(v any) time.Time {
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

// newDatePicker creates a date picker initialized with the given value.
// If value is zero, today's date is used.
func newDatePicker(value time.Time) *datePicker {
	if value.IsZero() {
		value = localToday()
	}
	value = time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.Local)
	return &datePicker{
		mode:     datePickerSegment,
		original: value,
		date:     value,
		segment:  segYear,
	}
}

func (dp *datePicker) Value() string {
	return dp.date.Format("2006-01-02")
}

func (dp *datePicker) Mode() datePickerMode {
	return dp.mode
}

// Update handles a key press and returns whether the date picker consumed it,
// and whether editing is done (confirmed or cancelled).
func (dp *datePicker) Update(msg tea.KeyPressMsg) (bool, bool, bool) {
	switch dp.mode {
	case datePickerSegment:
		return dp.updateSegment(msg)
	case datePickerCalendar:
		return dp.updatePickerCalendar(msg)
	}
	return false, false, false
}

func (dp *datePicker) updateSegment(msg tea.KeyPressMsg) (bool, bool, bool) {
	key := msg.String()

	switch key {
	case "enter":
		return true, true, true
	case "esc":
		dp.date = dp.original
		return true, true, false
	case "c":
		dp.mode = datePickerCalendar
		dp.clearDigitBuf()
		return true, false, false
	case "left", "shift+tab":
		dp.clearDigitBuf()
		if dp.segment > segYear {
			dp.segment--
		}
		return true, false, false
	case "right", "tab":
		dp.clearDigitBuf()
		if dp.segment < segDay {
			dp.segment++
		}
		return true, false, false
	case "up":
		dp.clearDigitBuf()
		dp.incrementSegment(1)
		return true, false, false
	case "down":
		dp.clearDigitBuf()
		dp.incrementSegment(-1)
		return true, false, false
	default:
		if len(key) == 1 && key[0] >= '0' && key[0] <= '9' {
			dp.handleDigit(key[0])
			return true, false, false
		}
	}
	return false, false, false
}

func (dp *datePicker) updatePickerCalendar(msg tea.KeyPressMsg) (bool, bool, bool) {
	key := msg.String()

	switch key {
	case "enter":
		return true, true, true
	case "esc":
		dp.date = dp.original
		return true, true, false
	case "c":
		dp.mode = datePickerSegment
		dp.segment = segYear
		dp.clearDigitBuf()
		return true, false, false
	case "left", "h":
		dp.date = dp.date.AddDate(0, 0, -1)
		return true, false, false
	case "right", "l":
		dp.date = dp.date.AddDate(0, 0, 1)
		return true, false, false
	case "up", "k":
		dp.date = dp.date.AddDate(0, 0, -7)
		return true, false, false
	case "down", "j":
		dp.date = dp.date.AddDate(0, 0, 7)
		return true, false, false
	case "H":
		dp.prevMonth()
		return true, false, false
	case "L":
		dp.nextMonth()
		return true, false, false
	case "t":
		dp.date = localToday()
		return true, false, false
	}
	return false, false, false
}

func (dp *datePicker) incrementSegment(delta int) {
	y, m, d := dp.date.Year(), dp.date.Month(), dp.date.Day()

	switch dp.segment {
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
		dp.date = dp.date.AddDate(0, 0, delta)
		return
	}

	maxDay := daysInMonth(y, m)
	if d > maxDay {
		d = maxDay
	}
	dp.date = time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

func (dp *datePicker) handleDigit(digit byte) {
	dp.digitBuf += string(digit)
	dp.digitCount++

	y, m, d := dp.date.Year(), dp.date.Month(), dp.date.Day()

	switch dp.segment {
	case segYear:
		if dp.digitCount >= 4 {
			val, _ := strconv.Atoi(dp.digitBuf)
			if val < 1 {
				val = 1
			}
			if val > 9999 {
				val = 9999
			}
			y = val
			dp.clearDigitBuf()
			dp.segment = segMonth
		} else {
			return
		}
	case segMonth:
		val, _ := strconv.Atoi(dp.digitBuf)
		if dp.digitCount >= 2 || val > 1 {
			if val < 1 {
				val = 1
			}
			if val > 12 {
				val = 12
			}
			m = time.Month(val)
			dp.clearDigitBuf()
			dp.segment = segDay
		} else {
			return
		}
	case segDay:
		val, _ := strconv.Atoi(dp.digitBuf)
		maxDay := daysInMonth(y, m)
		if dp.digitCount >= 2 || val > 3 {
			if val < 1 {
				val = 1
			}
			if val > maxDay {
				val = maxDay
			}
			d = val
			dp.clearDigitBuf()
		} else {
			return
		}
	}

	maxDay := daysInMonth(y, m)
	if d > maxDay {
		d = maxDay
	}
	dp.date = time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

func (dp *datePicker) prevMonth() {
	y, m, d := dp.date.Year(), dp.date.Month(), dp.date.Day()
	m--
	if m < 1 {
		m = 12
		y--
	}
	maxDay := daysInMonth(y, m)
	if d > maxDay {
		d = maxDay
	}
	dp.date = time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

func (dp *datePicker) nextMonth() {
	y, m, d := dp.date.Year(), dp.date.Month(), dp.date.Day()
	m++
	if m > 12 {
		m = 1
		y++
	}
	maxDay := daysInMonth(y, m)
	if d > maxDay {
		d = maxDay
	}
	dp.date = time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

func (dp *datePicker) clearDigitBuf() {
	dp.digitBuf = ""
	dp.digitCount = 0
}

// View returns the inline segment display (used in both modes).
func (dp *datePicker) View() string {
	return dp.viewSegment()
}

// PickerOverlay returns the calendar grid for overlay rendering.
// Returns empty string when not in calendar mode.
func (dp *datePicker) PickerOverlay() string {
	if dp.mode != datePickerCalendar {
		return ""
	}
	return dp.viewPickerCalendar()
}

func (dp *datePicker) viewSegment() string {
	y, m, d := dp.date.Year(), dp.date.Month(), dp.date.Day()
	dow := dp.date.Weekday().String()[:3]

	yearStr := fmt.Sprintf("%04d", y)
	monthStr := fmt.Sprintf("%02d", int(m))
	dayStr := fmt.Sprintf("%02d", d)

	switch dp.segment {
	case segYear:
		if dp.digitCount > 0 {
			yearStr = dp.digitBuf + strings.Repeat("_", 4-dp.digitCount)
		}
	case segMonth:
		if dp.digitCount > 0 {
			monthStr = dp.digitBuf + strings.Repeat("_", 2-dp.digitCount)
		}
	case segDay:
		if dp.digitCount > 0 {
			dayStr = dp.digitBuf + strings.Repeat("_", 2-dp.digitCount)
		}
	}

	var parts [3]string
	parts[0] = yearStr
	parts[1] = monthStr
	parts[2] = dayStr

	switch dp.segment {
	case segYear:
		parts[0] = dpSegmentFocusStyle.Render(yearStr)
	case segMonth:
		parts[1] = dpSegmentFocusStyle.Render(monthStr)
	case segDay:
		parts[2] = dpSegmentFocusStyle.Render(dayStr)
	}

	return fmt.Sprintf("%s-%s-%s  %s", parts[0], parts[1], parts[2], dow)
}

func (dp *datePicker) viewPickerCalendar() string {
	y, m := dp.date.Year(), dp.date.Month()
	today := localToday()

	var rows []string

	header := fmt.Sprintf("%s %d", m.String()[:3], y)
	pad := (dpGridWidth - len(header)) / 2
	if pad < 0 {
		pad = 0
	}
	rows = append(rows, strings.Repeat(" ", pad)+header)
	rows = append(rows, "Mo Tu We Th Fr Sa Su")

	first := time.Date(y, m, 1, 0, 0, 0, 0, time.Local)
	offset := int(first.Weekday())
	if offset == 0 {
		offset = 7
	}
	offset--

	lastDay := daysInMonth(y, m)

	var row strings.Builder
	col := 0

	for i := 0; i < offset; i++ {
		row.WriteString("   ")
		col++
	}

	for day := 1; day <= lastDay; day++ {
		dayDate := time.Date(y, m, day, 0, 0, 0, 0, time.Local)
		dayStr := fmt.Sprintf("%2d", day)

		isCursor := dp.date.Year() == y && dp.date.Month() == m && dp.date.Day() == day
		isToday := dayDate.Equal(today)

		if isCursor {
			row.WriteString(dpCursorStyle.Render(dayStr))
		} else if isToday {
			row.WriteString(dpTodayStyle.Render(dayStr))
		} else {
			row.WriteString(dayStr)
		}
		col++

		if col == 7 {
			rows = append(rows, row.String())
			row.Reset()
			col = 0
		} else if day < lastDay {
			row.WriteString(" ")
		}
	}
	if col > 0 {
		rows = append(rows, row.String())
	}

	for len(rows) < dpGridHeight {
		rows = append(rows, "")
	}

	return lipgloss.NewStyle().Width(dpGridWidth).Height(dpGridHeight).Render(
		strings.Join(rows, "\n"),
	)
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.Local).Day()
}
