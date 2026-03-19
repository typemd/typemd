package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
)

// viewPicker wraps a huh.Form with a single Select field for choosing a view.
type viewPicker struct {
	typeName string
	form     *huh.Form
	selected string
}

// newViewPicker creates a picker for the given view names.
// The "default" view is always included as the first option.
func newViewPicker(typeName string, viewNames []string) *viewPicker {
	vp := &viewPicker{
		typeName: typeName,
		selected: "default",
	}

	opts := []huh.Option[string]{
		huh.NewOption("default", "default"),
	}
	for _, name := range viewNames {
		if name == "default" {
			continue // already included
		}
		opts = append(opts, huh.NewOption(name, name))
	}

	vp.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a view").
				Options(opts...).
				Value(&vp.selected).
				Filtering(true).
				Height(min(len(opts)+2, 10)),
		),
	).WithWidth(40).WithShowHelp(false)

	return vp
}

// Init returns the form's init command.
func (vp *viewPicker) Init() tea.Cmd {
	return vp.form.Init()
}

// Update handles messages for the picker.
func (vp *viewPicker) Update(msg tea.Msg) (*viewPicker, tea.Cmd) {
	// Esc cancels the picker
	if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == "esc" {
		return nil, nil
	}

	form, cmd := vp.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		vp.form = f
	}

	// Check if form is completed or aborted
	switch vp.form.State {
	case huh.StateCompleted:
		typeName := vp.typeName
		selected := vp.selected
		return nil, func() tea.Msg {
			return openViewMsg{TypeName: typeName, ViewName: selected}
		}
	case huh.StateAborted:
		return nil, nil
	}

	return vp, cmd
}

// View renders the picker as a centered overlay.
func (vp *viewPicker) View(width, height int) string {
	content := vp.form.View()

	return renderPopup(content, width, height, 0)
}
