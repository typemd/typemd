package tui

import "charm.land/bubbles/v2/key"

type keyMap struct {
	Up          key.Binding
	Down        key.Binding
	Enter       key.Binding
	Tab         key.Binding
	Search      key.Binding
	Quit        key.Binding
	GrowPanel   key.Binding
	ShrinkPanel key.Binding
	ToggleProps key.Binding
	ToggleWrap  key.Binding
	Help        key.Binding
	EnterEdit   key.Binding
	ExitEdit    key.Binding
	NewObject   key.Binding
	QuickCreate key.Binding
}

var keys = keyMap{
	Up:          key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
	Down:        key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
	Enter:       key.NewBinding(key.WithKeys("enter", " "), key.WithHelp("enter", "select")),
	Tab:         key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "switch panel")),
	Search:      key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "search")),
	Quit:        key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	GrowPanel:   key.NewBinding(key.WithKeys("]"), key.WithHelp("]", "grow panel")),
	ShrinkPanel: key.NewBinding(key.WithKeys("["), key.WithHelp("[", "shrink panel")),
	ToggleProps: key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "toggle properties")),
	ToggleWrap:  key.NewBinding(key.WithKeys("w"), key.WithHelp("w", "toggle wrap")),
	Help:        key.NewBinding(key.WithKeys("?", "h"), key.WithHelp("?/h", "help")),
	EnterEdit:   key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit")),
	ExitEdit:    key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "exit edit")),
	NewObject:   key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "new object")),
	QuickCreate: key.NewBinding(key.WithKeys("N"), key.WithHelp("N", "quick create (batch)")),
}
