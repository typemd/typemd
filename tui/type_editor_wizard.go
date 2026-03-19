package tui

import (
	"fmt"
	"strings"

	"github.com/typemd/typemd/core"
	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/textinput"
)

// ── Add Property Wizard ─────────────────────────────────────────────────────

func (te *typeEditor) startAddWizard() {
	ni := textinput.New()
	ni.Placeholder = "property name"
	ni.CharLimit = 50
	ni.Focus()

	oi := textinput.New()
	oi.Placeholder = "option1, option2, ..."
	oi.CharLimit = 200

	ii := textinput.New()
	ii.Placeholder = "inverse property name"
	ii.CharLimit = 50

	// Gather target types
	var targets []string
	if te.vault != nil {
		targets = te.vault.ListTypes()
	}

	te.wizard = &addPropWizard{
		step:            wizStepName,
		nameInput:       ni,
		typeList:        propertyTypeList,
		optionsInput:    oi,
		relTargets:      targets,
		relInverseInput: ii,
	}
	te.mode = teModeAddWizard
}

func (te *typeEditor) updateAddWizard(msg tea.KeyPressMsg) (*typeEditor, tea.Cmd) {
	if te.wizard == nil {
		te.mode = teModeView
		return te, nil
	}
	wiz := te.wizard

	switch wiz.step {
	case wizStepName:
		return te.updateWizardName(msg)
	case wizStepType:
		return te.updateWizardType(msg)
	case wizStepOptions:
		return te.updateWizardOptions(msg)
	case wizStepRelation:
		return te.updateWizardRelation(msg)
	}
	return te, nil
}

func (te *typeEditor) updateWizardName(msg tea.KeyPressMsg) (*typeEditor, tea.Cmd) {
	wiz := te.wizard
	switch msg.String() {
	case "enter":
		name := strings.TrimSpace(wiz.nameInput.Value())
		if name == "" {
			return te, nil
		}
		// Check for duplicate
		for _, p := range te.schema.Properties {
			if p.Name == name {
				te.saveErr = fmt.Sprintf("property %q already exists", name)
				return te, nil
			}
		}
		// Check reserved system property names
		if core.IsSystemProperty(name) {
			te.saveErr = fmt.Sprintf("%q is a reserved system property", name)
			return te, nil
		}
		te.saveErr = ""
		wiz.propName = name
		wiz.step = wizStepType
		return te, nil
	case "esc":
		te.cancelWizard()
		return te, nil
	}
	var cmd tea.Cmd
	wiz.nameInput, cmd = wiz.nameInput.Update(msg)
	return te, cmd
}

func (te *typeEditor) updateWizardType(msg tea.KeyPressMsg) (*typeEditor, tea.Cmd) {
	wiz := te.wizard
	switch msg.String() {
	case "up", "k":
		if wiz.typeCursor > 0 {
			wiz.typeCursor--
		}
	case "down", "j":
		if wiz.typeCursor < len(wiz.typeList)-1 {
			wiz.typeCursor++
		}
	case "enter":
		wiz.propType = wiz.typeList[wiz.typeCursor]
		switch wiz.propType {
		case "select", "multi_select":
			wiz.step = wizStepOptions
			wiz.optionsInput.Focus()
		case "relation":
			wiz.step = wizStepRelation
		default:
			te.finishWizard()
		}
	case "esc":
		wiz.step = wizStepName
	}
	return te, nil
}

func (te *typeEditor) updateWizardOptions(msg tea.KeyPressMsg) (*typeEditor, tea.Cmd) {
	wiz := te.wizard
	switch msg.String() {
	case "enter":
		te.finishWizard()
	case "esc":
		wiz.step = wizStepType
	default:
		var cmd tea.Cmd
		wiz.optionsInput, cmd = wiz.optionsInput.Update(msg)
		return te, cmd
	}
	return te, nil
}

func (te *typeEditor) updateWizardRelation(msg tea.KeyPressMsg) (*typeEditor, tea.Cmd) {
	wiz := te.wizard
	switch msg.String() {
	case "up", "k":
		switch wiz.relFieldCursor {
		case 0: // target list
			if wiz.relTargetCursor > 0 {
				wiz.relTargetCursor--
			}
		default:
			if wiz.relFieldCursor > 0 {
				wiz.relFieldCursor--
			}
		}
	case "down", "j":
		switch wiz.relFieldCursor {
		case 0: // target list
			if wiz.relTargetCursor < len(wiz.relTargets)-1 {
				wiz.relTargetCursor++
			}
		default:
			if wiz.relFieldCursor < 3 {
				wiz.relFieldCursor++
			}
		}
	case "tab":
		wiz.relFieldCursor = (wiz.relFieldCursor + 1) % 4
		if wiz.relFieldCursor == 3 && !wiz.relBidir {
			wiz.relFieldCursor = 0 // skip inverse if not bidirectional
		}
		if wiz.relFieldCursor == 3 {
			wiz.relInverseInput.Focus()
		} else {
			wiz.relInverseInput.Blur()
		}
	case "enter":
		if wiz.relFieldCursor == 1 { // toggle multiple
			wiz.relMultiple = !wiz.relMultiple
		} else if wiz.relFieldCursor == 2 { // toggle bidirectional
			wiz.relBidir = !wiz.relBidir
		} else if wiz.relFieldCursor == 0 || wiz.relFieldCursor == 3 {
			te.finishWizard()
		}
	case "esc":
		wiz.step = wizStepType
		wiz.relInverseInput.Blur()
	case " ":
		if wiz.relFieldCursor == 1 {
			wiz.relMultiple = !wiz.relMultiple
		} else if wiz.relFieldCursor == 2 {
			wiz.relBidir = !wiz.relBidir
		}
	default:
		if wiz.relFieldCursor == 3 {
			var cmd tea.Cmd
			wiz.relInverseInput, cmd = wiz.relInverseInput.Update(msg)
			return te, cmd
		}
	}
	return te, nil
}

func (te *typeEditor) finishWizard() {
	wiz := te.wizard
	prop := core.Property{
		Name: wiz.propName,
		Type: wiz.propType,
	}

	switch wiz.propType {
	case "select", "multi_select":
		raw := strings.TrimSpace(wiz.optionsInput.Value())
		if raw != "" {
			for _, v := range strings.Split(raw, ",") {
				v = strings.TrimSpace(v)
				if v != "" {
					prop.Options = append(prop.Options, core.Option{Value: v})
				}
			}
		}
	case "relation":
		if len(wiz.relTargets) > 0 && wiz.relTargetCursor < len(wiz.relTargets) {
			prop.Target = wiz.relTargets[wiz.relTargetCursor]
		}
		prop.Multiple = wiz.relMultiple
		prop.Bidirectional = wiz.relBidir
		if wiz.relBidir {
			prop.Inverse = strings.TrimSpace(wiz.relInverseInput.Value())
		}
	}

	te.schema.Properties = append(te.schema.Properties, prop)
	te.save()
	te.cancelWizard()
}

func (te *typeEditor) cancelWizard() {
	te.wizard = nil
	te.mode = teModeView
	te.saveErr = ""
}
