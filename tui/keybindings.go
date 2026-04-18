package tui

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/typemd/typemd/core"
	"charm.land/bubbles/v2/key"
)

// Rebindable action names. These are the user-visible snake_case keys used in
// .typemd/config.yaml under tui.keybindings.<action>.
//
// Not rebindable: tab (panel switching), esc (cancel/exit), shift+tab,
// arrow keys inside edit modes, and literal character input inside popups.
// These are reserved because they're used outside updateNormal's switch and
// making them rebindable would risk breaking modal dispatch.
const (
	ActionUp            = "up"
	ActionDown          = "down"
	ActionEnter         = "enter"
	ActionSearch        = "search"
	ActionQuit          = "quit"
	ActionGrowPanel     = "grow_panel"
	ActionShrinkPanel   = "shrink_panel"
	ActionFocusMode     = "focus_mode"
	ActionToggleProps   = "toggle_props"
	ActionToggleWrap    = "toggle_wrap"
	ActionHelp          = "help"
	ActionEnterEdit     = "enter_edit"
	ActionNewObject     = "new_object"
	ActionQuickCreate   = "quick_create"
	ActionRename        = "rename"
	ActionStats         = "stats"
	ActionAIGenerate    = "ai_generate"
	ActionSchemaExplore = "schema_explore"
	ActionSettings      = "settings"
)

// keybindingDefault describes one rebindable action: its default key list, the
// help description shown to the user, and the primary key used for dispatch
// translation.
type keybindingDefault struct {
	// Keys are the default keys the action responds to. First entry is treated
	// as the "primary" key that switch-cases in update.go match against.
	Keys []string
	// HelpKey is the label shown in the help popup (e.g., "↑/k", "ctrl+s").
	HelpKey string
	// HelpDesc is the human-readable description (e.g., "up", "stats").
	HelpDesc string
}

// defaultKeybindings is the single source of truth for the rebindable
// actions. Shared between core/config.go registration and tui/keys.go defaults
// so the two cannot drift. Non-rebindable keys (tab, esc) are held on the
// keyMap directly via hardcoded bindings in buildKeyMap.
var defaultKeybindings = map[string]keybindingDefault{
	ActionUp:            {Keys: []string{"up", "k"}, HelpKey: "↑/k", HelpDesc: "up"},
	ActionDown:          {Keys: []string{"down", "j"}, HelpKey: "↓/j", HelpDesc: "down"},
	ActionEnter:         {Keys: []string{"enter"}, HelpKey: "enter", HelpDesc: "select"},
	ActionSearch:        {Keys: []string{"/"}, HelpKey: "/", HelpDesc: "search"},
	ActionQuit:          {Keys: []string{"q", "ctrl+c"}, HelpKey: "q", HelpDesc: "quit"},
	ActionGrowPanel:     {Keys: []string{"="}, HelpKey: "=", HelpDesc: "grow panel"},
	ActionShrinkPanel:   {Keys: []string{"-"}, HelpKey: "-", HelpDesc: "shrink panel"},
	ActionFocusMode:     {Keys: []string{"."}, HelpKey: ".", HelpDesc: "focus mode"},
	ActionToggleProps:   {Keys: []string{"p"}, HelpKey: "p", HelpDesc: "toggle properties"},
	ActionToggleWrap:    {Keys: []string{"w"}, HelpKey: "w", HelpDesc: "toggle wrap"},
	ActionHelp:          {Keys: []string{"?", "h"}, HelpKey: "?/h", HelpDesc: "help"},
	ActionEnterEdit:     {Keys: []string{"e"}, HelpKey: "e", HelpDesc: "edit"},
	ActionNewObject:     {Keys: []string{"n"}, HelpKey: "n", HelpDesc: "new object"},
	ActionQuickCreate:   {Keys: []string{"N"}, HelpKey: "N", HelpDesc: "quick create (batch)"},
	ActionRename:        {Keys: []string{"r"}, HelpKey: "r", HelpDesc: "rename"},
	ActionStats:         {Keys: []string{"ctrl+s"}, HelpKey: "ctrl+s", HelpDesc: "stats"},
	ActionAIGenerate:    {Keys: []string{"g"}, HelpKey: "g", HelpDesc: "AI generate"},
	ActionSchemaExplore: {Keys: []string{"ctrl+e"}, HelpKey: "ctrl+e", HelpDesc: "AI schema explore"},
	ActionSettings:      {Keys: []string{","}, HelpKey: ",", HelpDesc: "settings"},
}

// rebindableActionsList holds action names in a stable alphabetical order.
// Precomputed once at package load so hot paths (buildKeyMap, translate) don't
// pay sort+alloc cost per call.
var rebindableActionsList = func() []string {
	names := make([]string, 0, len(defaultKeybindings))
	for name := range defaultKeybindings {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}()

// rebindableActions returns all action names in a stable alphabetical order.
func rebindableActions() []string {
	return rebindableActionsList
}

// KeybindingIssueKind classifies a keybinding validation failure.
type KeybindingIssueKind int

const (
	// IssueUnknownAction: config references an action name that does not exist.
	IssueUnknownAction KeybindingIssueKind = iota
	// IssueInvalidKey: the key string failed validation (typo, bad modifier).
	IssueInvalidKey
	// IssueDuplicateKey: two or more actions resolve to the same key.
	IssueDuplicateKey
)

// KeybindingIssue describes one validation failure produced by buildKeyMap.
// Callers render these as toast warnings at startup.
type KeybindingIssue struct {
	Kind    KeybindingIssueKind
	Action  string   // populated for UnknownAction and InvalidKey
	Value   string   // the offending raw value for UnknownAction/InvalidKey
	Reason  string   // populated for InvalidKey: why validateKeyString rejected it
	Key     string   // populated for DuplicateKey: the duplicated key
	Actions []string // populated for DuplicateKey: actions sharing the key
}

// Message returns a human-readable message for a keybinding issue.
func (i KeybindingIssue) Message() string {
	switch i.Kind {
	case IssueUnknownAction:
		return fmt.Sprintf("Unknown keybinding action %q in tui.keybindings (ignored)", i.Action)
	case IssueInvalidKey:
		if i.Reason != "" {
			return fmt.Sprintf("Invalid key %q for tui.keybindings.%s (%s) — using default", i.Value, i.Action, i.Reason)
		}
		return fmt.Sprintf("Invalid key %q for tui.keybindings.%s — using default", i.Value, i.Action)
	case IssueDuplicateKey:
		return fmt.Sprintf("Duplicate keybinding %q used by: %s", i.Key, strings.Join(i.Actions, ", "))
	}
	return ""
}

// validKeyNames is the set of non-rune named keys accepted by validateKeyString.
// This is a pragmatic subset of what bubbletea/bubbles recognises — enough to
// cover everything a typemd user would reasonably bind.
var validKeyNames = map[string]struct{}{
	"esc": {}, "escape": {}, "enter": {}, "return": {}, "tab": {}, "space": {},
	"backspace": {}, "delete": {}, "del": {},
	"up": {}, "down": {}, "left": {}, "right": {},
	"home": {}, "end": {}, "pgup": {}, "pgdown": {}, "pageup": {}, "pagedown": {},
	"insert": {}, "ins": {},
	"f1": {}, "f2": {}, "f3": {}, "f4": {}, "f5": {}, "f6": {},
	"f7": {}, "f8": {}, "f9": {}, "f10": {}, "f11": {}, "f12": {},
}

// validModifiers accepts ctrl, alt, shift in lower case.
var validModifiers = map[string]struct{}{
	"ctrl":  {},
	"alt":   {},
	"shift": {},
}

// validateKeyString reports whether s is a syntactically valid key descriptor
// for tui.keybindings. It accepts:
//   - single printable runes ("a", "/", "=", "?"),
//   - named keys ("esc", "enter", "tab", "up", "f1"),
//   - any of those prefixed with "ctrl+", "alt+", "shift+" (at most one of each),
//     joined with "+" (e.g. "ctrl+shift+a").
//
// It rejects empty strings, empty segments (any "+" with nothing on one side,
// including "ctrl+", "+a", "ctrl++s"), duplicate modifiers, and unknown named keys.
func validateKeyString(s string) error {
	if s == "" {
		return fmt.Errorf("empty key")
	}

	parts := strings.Split(s, "+")
	if slices.Contains(parts, "") {
		return fmt.Errorf("empty segment")
	}

	keyPart := parts[len(parts)-1]
	modParts := parts[:len(parts)-1]

	seenMods := map[string]bool{}
	for _, m := range modParts {
		lower := strings.ToLower(m)
		if _, ok := validModifiers[lower]; !ok {
			return fmt.Errorf("unknown modifier %q", m)
		}
		if seenMods[lower] {
			return fmt.Errorf("duplicate modifier %q", m)
		}
		seenMods[lower] = true
	}

	// Names are case-insensitive; runes case-sensitive (bubbletea treats "n" and "N" differently).
	if _, ok := validKeyNames[strings.ToLower(keyPart)]; ok {
		return nil
	}
	if len([]rune(keyPart)) == 1 {
		return nil
	}
	return fmt.Errorf("unknown key %q", keyPart)
}

// buildKeyMap merges user overrides from cfg.TUI.Keybindings on top of
// defaultKeybindings and returns the resolved keyMap plus any validation
// issues. Invalid overrides fall back to the default — the returned keyMap is
// always safe to use.
//
// Resolution rules:
//   - unknown action name → IssueUnknownAction (ignored; no effect on keyMap)
//   - value == ""         → treated as "use default" (no issue reported)
//   - invalid key string  → IssueInvalidKey (action keeps its default)
//   - duplicate key       → IssueDuplicateKey (both actions keep their
//     resolved keys; user sees a warning and can fix)
func buildKeyMap(cfg *core.VaultConfig) (keyMap, []KeybindingIssue) {
	var issues []KeybindingIssue

	resolved := make(map[string][]string, len(defaultKeybindings))
	helpKey := make(map[string]string, len(defaultKeybindings))
	for name, def := range defaultKeybindings {
		resolved[name] = append([]string(nil), def.Keys...)
		helpKey[name] = def.HelpKey
	}

	if cfg != nil {
		names := make([]string, 0, len(cfg.TUI.Keybindings))
		for k := range cfg.TUI.Keybindings {
			names = append(names, k)
		}
		sort.Strings(names)

		for _, action := range names {
			raw := strings.TrimSpace(cfg.TUI.Keybindings[action])
			if _, known := defaultKeybindings[action]; !known {
				issues = append(issues, KeybindingIssue{Kind: IssueUnknownAction, Action: action, Value: raw})
				continue
			}
			if raw == "" {
				continue
			}
			if err := validateKeyString(raw); err != nil {
				issues = append(issues, KeybindingIssue{Kind: IssueInvalidKey, Action: action, Value: raw, Reason: err.Error()})
				continue
			}
			resolved[action] = []string{raw}
			helpKey[action] = raw
		}
	}

	// Build keyToPrimary: maps every bound key string → the primary default key
	// for whichever action owns it. Used by translate() in the dispatch hot
	// path so it's an O(1) lookup with zero allocations per keystroke.
	// Also used to detect duplicate-key conflicts.
	keyToPrimary := make(map[string]string, len(resolved))
	keyToActions := map[string][]string{}
	for _, action := range rebindableActions() {
		primary := defaultKeybindings[action].Keys[0]
		for _, k := range resolved[action] {
			keyToPrimary[k] = primary
			keyToActions[k] = append(keyToActions[k], action)
		}
	}

	dupKeys := make([]string, 0)
	for k, acts := range keyToActions {
		if len(acts) > 1 {
			dupKeys = append(dupKeys, k)
		}
	}
	sort.Strings(dupKeys)
	for _, k := range dupKeys {
		acts := append([]string(nil), keyToActions[k]...)
		sort.Strings(acts)
		issues = append(issues, KeybindingIssue{Kind: IssueDuplicateKey, Key: k, Actions: acts})
	}

	mk := func(action string) key.Binding {
		def := defaultKeybindings[action]
		return key.NewBinding(
			key.WithKeys(resolved[action]...),
			key.WithHelp(helpKey[action], def.HelpDesc),
		)
	}

	km := keyMap{
		Up:            mk(ActionUp),
		Down:          mk(ActionDown),
		Enter:         mk(ActionEnter),
		Search:        mk(ActionSearch),
		Quit:          mk(ActionQuit),
		GrowPanel:     mk(ActionGrowPanel),
		ShrinkPanel:   mk(ActionShrinkPanel),
		FocusMode:     mk(ActionFocusMode),
		ToggleProps:   mk(ActionToggleProps),
		ToggleWrap:    mk(ActionToggleWrap),
		Help:          mk(ActionHelp),
		EnterEdit:     mk(ActionEnterEdit),
		NewObject:     mk(ActionNewObject),
		QuickCreate:   mk(ActionQuickCreate),
		Rename:        mk(ActionRename),
		Stats:         mk(ActionStats),
		AIGenerate:    mk(ActionAIGenerate),
		SchemaExplore: mk(ActionSchemaExplore),
		Settings:      mk(ActionSettings),

		// Non-rebindable: reserved for global navigation/cancellation.
		Tab:      key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch panel")),
		ExitEdit: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "exit edit")),
	}
	km.keyToPrimary = keyToPrimary
	return km, issues
}

// defaultKeyMap returns a keyMap with every action bound to its compile-time
// default. Equivalent to buildKeyMap(nil) without the issue slice.
func defaultKeyMap() keyMap {
	km, _ := buildKeyMap(nil)
	return km
}

// init registers every rebindable action under "tui.keybindings.<action>" so
// `tmd config get/set` and the TUI Settings page work for them out of the box.
// Registered once at package load — matching the other static entries in
// core/config.go's configKeyRegistry.
func init() {
	for action, def := range defaultKeybindings {
		core.RegisterConfigKey(
			"tui.keybindings."+action,
			func(cfg *core.VaultConfig) string {
				if cfg == nil {
					return ""
				}
				return cfg.TUI.Keybindings[action]
			},
			func(cfg *core.VaultConfig, value string) {
				if cfg.TUI.Keybindings == nil {
					cfg.TUI.Keybindings = map[string]string{}
				}
				cfg.TUI.Keybindings[action] = value
			},
			def.HelpKey,
			"TUI keybinding for "+def.HelpDesc,
		)
	}
}

// translate maps an incoming key string (e.g. msg.String()) to the primary
// default key string for whichever action currently owns that key. If the key
// is not bound to any rebindable action, translate returns the input unchanged.
//
// This is the dispatch bridge: update.go keeps writing `case "ctrl+s"` (stats)
// and translate makes a user-configured `ctrl+d` route through that same case.
// O(1) lookup against keyToPrimary, populated once by buildKeyMap.
func (k keyMap) translate(s string) string {
	if p, ok := k.keyToPrimary[s]; ok {
		return p
	}
	return s
}
