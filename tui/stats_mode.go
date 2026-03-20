package tui

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/typemd/typemd/core"
	"github.com/typemd/typemd/tui/widget"
	tea "charm.land/bubbletea/v2"
)

// statsScreen tracks which screen the stats mode is showing.
type statsScreen int

const (
	statsOverview statsScreen = iota
	statsDetail
)

// statsMode is a sub-model for the fullscreen stats dashboard.
type statsMode struct {
	vault  *core.Vault
	screen statsScreen

	// Vault overview data
	vaultStats *core.VaultStats

	// Type detail data
	typeStats    *core.TypeStats
	detailType   string
	detailCursor int
	detailScroll int

	// Overview navigation
	cursor   int
	scroll   int
	maxNameW int // cached max display name width

	// Layout
	typeLayout string // "fullscreen" or "popup"
	width      int
	height     int
}

// newStatsMode creates a new statsMode and loads vault stats.
func newStatsMode(vault *core.Vault, layout string) *statsMode {
	if layout == "" {
		layout = "fullscreen"
	}
	sm := &statsMode{
		vault:      vault,
		screen:     statsOverview,
		typeLayout: layout,
	}
	sm.loadVaultStats()
	return sm
}

// SetSize updates the available rendering dimensions.
func (sm *statsMode) SetSize(w, h int) {
	sm.width = w
	sm.height = h
}

// loadVaultStats fetches vault-wide statistics and sorts by count descending.
func (sm *statsMode) loadVaultStats() {
	if sm.vault == nil {
		return
	}
	stats, err := sm.vault.VaultStats()
	if err != nil {
		return
	}
	// Sort by count descending
	sort.Slice(stats.Types, func(i, j int) bool {
		return stats.Types[i].Count > stats.Types[j].Count
	})
	sm.vaultStats = stats

	// Cache max display name width for rendering
	sm.maxNameW = 0
	for _, ts := range stats.Types {
		if w := runewidth.StringWidth(ts.DisplayName()); w > sm.maxNameW {
			sm.maxNameW = w
		}
	}
}

// loadTypeStats fetches per-property statistics for a specific type.
func (sm *statsMode) loadTypeStats(typeName string) {
	if sm.vault == nil {
		return
	}
	stats, err := sm.vault.TypeStats(typeName)
	if err != nil {
		return
	}
	sm.typeStats = stats
	sm.detailType = typeName
	sm.detailCursor = 0
	sm.detailScroll = 0
}

// typeCount returns the number of types in the vault stats.
func (sm *statsMode) typeCount() int {
	if sm.vaultStats == nil {
		return 0
	}
	return len(sm.vaultStats.Types)
}

// Update handles key events for the stats mode.
func (sm *statsMode) Update(msg tea.Msg) (*statsMode, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch sm.screen {
		case statsOverview:
			return sm.updateOverview(msg)
		case statsDetail:
			return sm.updateDetail(msg)
		}
	}
	return sm, nil
}

func (sm *statsMode) updateOverview(msg tea.KeyPressMsg) (*statsMode, tea.Cmd) {
	count := sm.typeCount()
	switch msg.String() {
	case "j", "down":
		if sm.cursor < count-1 {
			sm.cursor++
		}
		sm.scroll = widget.AdjustScroll(sm.cursor, sm.scroll, sm.height-6)
	case "k", "up":
		if sm.cursor > 0 {
			sm.cursor--
		}
		sm.scroll = widget.AdjustScroll(sm.cursor, sm.scroll, sm.height-6)
	case "enter":
		if count > 0 && sm.cursor < count {
			typeName := sm.vaultStats.Types[sm.cursor].Name
			sm.loadTypeStats(typeName)
			sm.screen = statsDetail
		}
	case "r":
		sm.loadVaultStats()
		if sm.cursor >= sm.typeCount() {
			sm.cursor = max(0, sm.typeCount()-1)
		}
	}
	return sm, nil
}

func (sm *statsMode) updateDetail(msg tea.KeyPressMsg) (*statsMode, tea.Cmd) {
	switch msg.String() {
	case "esc":
		sm.screen = statsOverview
		sm.typeStats = nil
		sm.detailType = ""
	case "j", "down":
		if sm.typeStats != nil {
			if sm.detailCursor < len(sm.typeStats.Properties)-1 {
				sm.detailCursor++
			}
			sm.detailScroll = widget.AdjustScroll(sm.detailCursor, sm.detailScroll, sm.height-6)
		}
	case "k", "up":
		if sm.detailCursor > 0 {
			sm.detailCursor--
		}
		sm.detailScroll = widget.AdjustScroll(sm.detailCursor, sm.detailScroll, sm.height-6)
	case "t":
		if sm.typeLayout == "fullscreen" {
			sm.typeLayout = "popup"
		} else {
			sm.typeLayout = "fullscreen"
		}
	case "r":
		if sm.detailType != "" {
			sm.loadTypeStats(sm.detailType)
		}
	}
	return sm, nil
}

// View renders the stats mode content.
func (sm *statsMode) View() string {
	switch sm.screen {
	case statsDetail:
		return sm.viewDetail()
	default:
		return sm.viewOverview()
	}
}

// viewOverview renders the vault overview screen.
func (sm *statsMode) viewOverview() string {
	if sm.vaultStats == nil || len(sm.vaultStats.Types) == 0 {
		return "  No objects in vault."
	}

	var lines []string

	// Header
	typeCount := len(sm.vaultStats.Types)
	lines = append(lines, fmt.Sprintf("  %d objects across %d types", sm.vaultStats.Total, typeCount))
	lines = append(lines, "")

	// Type list
	visibleH := sm.height - 4 // header + spacing
	if visibleH < 1 {
		visibleH = 1
	}

	for i, ts := range sm.vaultStats.Types {
		if i < sm.scroll {
			continue
		}
		if i >= sm.scroll+visibleH {
			break
		}

		emoji := ts.Emoji
		if emoji == "" {
			emoji = " "
		}
		name := ts.DisplayName()

		lastUpdated := ""
		if !ts.LastUpdated.IsZero() {
			lastUpdated = "  " + relativeTime(ts.LastUpdated)
		}

		line := fmt.Sprintf("  %s %-*s  %3d%s", padEmoji(emoji), sm.maxNameW, name, ts.Count, lastUpdated)

		if i == sm.cursor {
			line = highlightStyle.Render(line)
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// viewDetail renders the type detail screen.
func (sm *statsMode) viewDetail() string {
	if sm.typeStats == nil {
		return "  Loading..."
	}

	var lines []string

	// Header
	emoji := ""
	if sm.typeStats.Emoji != "" {
		emoji = padEmoji(sm.typeStats.Emoji) + " "
	}
	lines = append(lines, fmt.Sprintf("  %s%s (%d objects)", emoji, sm.typeStats.TypeName, sm.typeStats.Count))
	lines = append(lines, "  "+strings.Repeat("─", min(40, sm.width-4)))
	lines = append(lines, "")

	if len(sm.typeStats.Properties) == 0 {
		lines = append(lines, "  No properties defined.")
		return strings.Join(lines, "\n")
	}

	// Calculate bar chart width
	barMaxW := sm.width - 30 // leave space for labels and counts
	if barMaxW < 10 {
		barMaxW = 10
	}
	if barMaxW > 40 {
		barMaxW = 40
	}

	visibleH := sm.height - 6
	if visibleH < 1 {
		visibleH = 1
	}

	lineIdx := 0
	for _, ps := range sm.typeStats.Properties {
		propLines := sm.renderPropertyStats(ps, barMaxW)

		for _, pl := range propLines {
			if lineIdx >= sm.detailScroll && lineIdx < sm.detailScroll+visibleH {
				lines = append(lines, pl)
			}
			lineIdx++
		}
	}

	return strings.Join(lines, "\n")
}

// renderPropertyStats renders statistics for a single property.
func (sm *statsMode) renderPropertyStats(ps core.PropertyStats, barMaxW int) []string {
	var lines []string

	fillRate := ""
	if ps.Total > 0 {
		pct := float64(ps.Filled) / float64(ps.Total) * 100
		fillRate = fmt.Sprintf("  filled: %d/%d (%.0f%%)", ps.Filled, ps.Total, pct)
	}

	lines = append(lines, fmt.Sprintf("  %s (%s)%s", ps.Name, ps.Type, fillRate))

	if ps.Stats == nil {
		lines = append(lines, "")
		return lines
	}

	switch s := ps.Stats.(type) {
	case *core.NumberStats:
		lines = append(lines, fmt.Sprintf("    min: %g  max: %g  avg: %g  sum: %g", s.Min, s.Max, s.Avg, s.Sum))
	case *core.SelectStats:
		lines = append(lines, sm.renderDistribution(s.Distribution, barMaxW)...)
	case *core.CheckboxStats:
		total := s.TrueCount + s.FalseCount
		if total > 0 {
			truePct := float64(s.TrueCount) / float64(total) * 100
			lines = append(lines, fmt.Sprintf("    true: %d (%.0f%%)  false: %d (%.0f%%)", s.TrueCount, truePct, s.FalseCount, 100-truePct))
		}
	case *core.DateStats:
		lines = append(lines, fmt.Sprintf("    earliest: %s  latest: %s", s.Earliest, s.Latest))
	case *core.RelationStats:
		lines = append(lines, fmt.Sprintf("    total links: %d", s.Count))
	}

	lines = append(lines, "")
	return lines
}

// renderDistribution renders a select/multi_select frequency distribution as bar charts.
func (sm *statsMode) renderDistribution(dist map[string]int, barMaxW int) []string {
	if len(dist) == 0 {
		return nil
	}

	// Sort by count descending
	type entry struct {
		name  string
		count int
	}
	var entries []entry
	maxCount := 0
	maxLabelW := 0
	for name, count := range dist {
		entries = append(entries, entry{name, count})
		if count > maxCount {
			maxCount = count
		}
		if w := runewidth.StringWidth(name); w > maxLabelW {
			maxLabelW = w
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].count > entries[j].count
	})

	if maxLabelW > 15 {
		maxLabelW = 15
	}

	var lines []string
	for _, e := range entries {
		bar := renderBar(e.count, maxCount, barMaxW)
		name := e.name
		if runewidth.StringWidth(name) > maxLabelW {
			name = runewidth.Truncate(name, maxLabelW, "…")
		}
		lines = append(lines, fmt.Sprintf("    %-*s %s %d", maxLabelW, name, bar, e.count))
	}
	return lines
}

// barBlocks are Unicode block elements for sub-character precision in bar charts.
var barBlocks = []rune{' ', '▏', '▎', '▍', '▌', '▋', '▊', '▉', '█'}

// renderBar renders a proportional horizontal bar using block characters.
func renderBar(value, maxValue, maxWidth int) string {
	if maxValue == 0 || maxWidth == 0 {
		return ""
	}

	// Calculate the fractional width
	ratio := float64(value) / float64(maxValue)
	exactWidth := ratio * float64(maxWidth)

	fullBlocks := int(exactWidth)
	fraction := exactWidth - float64(fullBlocks)

	partialIdx := int(math.Round(fraction * 8))
	if partialIdx >= len(barBlocks) {
		partialIdx = len(barBlocks) - 1
	}

	bar := strings.Repeat("█", fullBlocks)
	if partialIdx > 0 && fullBlocks < maxWidth {
		bar += string(barBlocks[partialIdx])
	}

	return bar
}

// relativeTime formats a time as a human-readable relative string.
func relativeTime(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(d.Minutes())
		if m == 1 {
			return "1 min ago"
		}
		return fmt.Sprintf("%d mins ago", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		if h == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", h)
	case d < 30*24*time.Hour:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	default:
		return t.Format("2006-01-02")
	}
}

// HelpBar returns the help bar text for stats mode.
func (sm *statsMode) HelpBar() string {
	switch sm.screen {
	case statsDetail:
		return "  [STATS]  ↑↓: scroll  esc: back  t: toggle layout  r: refresh  q: quit"
	default:
		return "  [STATS]  ↑↓: navigate  enter: type detail  r: refresh  esc: exit  q: quit"
	}
}

// titleContent returns the title bar text for stats mode.
func (sm *statsMode) titleContent() string {
	switch sm.screen {
	case statsDetail:
		if sm.typeStats != nil {
			emoji := ""
			if sm.typeStats.Emoji != "" {
				emoji = padEmoji(sm.typeStats.Emoji) + " "
			}
			return fmt.Sprintf(" %s%s — Statistics", emoji, sm.typeStats.TypeName)
		}
		return " Type Statistics"
	default:
		return " Vault Statistics"
	}
}
