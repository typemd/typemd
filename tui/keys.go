package tui

import "charm.land/bubbles/v2/key"

// keyMap holds the resolved set of global TUI keybindings. Each field is a
// bubbles key.Binding; the resolved map mirrors the same information keyed by
// action name and is used by keyMap.translate to route user-customised keys
// to the correct dispatch case in update.go.
//
// See defaultKeybindings in keybindings.go for the single source of truth, and
// buildKeyMap to construct a keyMap with user overrides applied.
type keyMap struct {
	Up            key.Binding
	Down          key.Binding
	Enter         key.Binding
	Tab           key.Binding
	Search        key.Binding
	Quit          key.Binding
	GrowPanel     key.Binding
	ShrinkPanel   key.Binding
	FocusMode     key.Binding
	ToggleProps   key.Binding
	ToggleWrap    key.Binding
	Help          key.Binding
	EnterEdit     key.Binding
	ExitEdit      key.Binding
	NewObject     key.Binding
	QuickCreate   key.Binding
	Rename        key.Binding
	Stats         key.Binding
	AIGenerate    key.Binding
	SchemaExplore key.Binding
	Settings      key.Binding

	// keyToPrimary maps every bound key string to the primary default key for
	// whichever action currently owns it. Populated once by buildKeyMap so
	// translate() is an O(1) lookup on the dispatch hot path.
	keyToPrimary map[string]string
}

// keys is the package-level keyMap used by tests that pre-date the
// configurable-keybindings refactor. It always reflects the compile-time
// defaults. Production code paths read from the model's resolved keyMap
// (model.keys), not this variable.
var keys = defaultKeyMap()
