package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/typemd/typemd/core"
	tea "charm.land/bubbletea/v2"
	"github.com/mattn/go-isatty"
)

// resolveVault creates a Vault with the given path, defaulting to "." if empty.
func resolveVault(path string) *core.Vault {
	if path == "" {
		return core.NewVault(".")
	}
	return core.NewVault(path)
}

// openVault creates and opens a vault.
// The caller must defer vault.Close().
func openVault(path string) (*core.Vault, error) {
	vault := resolveVault(path)
	if err := vault.Open(); err != nil {
		return nil, err
	}
	return vault, nil
}

// printJSON marshals any value as indented JSON and prints it.
func printJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

// printObjects prints objects as JSON or one DisplayID per line.
func printObjects(objects []*core.Object, asJSON bool) error {
	if asJSON {
		return printJSON(objects)
	}
	for _, obj := range objects {
		fmt.Printf("%s/%s\n", obj.Type, obj.GetName())
	}
	return nil
}

// isInteractiveFunc checks whether stdin is an interactive terminal.
// Override in tests to simulate interactive mode.
var isInteractiveFunc = func() bool {
	return isatty.IsTerminal(os.Stdin.Fd()) || isatty.IsCygwinTerminal(os.Stdin.Fd())
}

// disambiguateFunc is the function used to resolve ambiguous matches.
// It takes a list of candidate items and returns the selected ID, or an error.
// Override in tests to bypass the Bubble Tea picker.
var disambiguateFunc = disambiguateWithPicker

// disambiguateWithPicker launches a Bubble Tea picker for interactive selection.
func disambiguateWithPicker(items []disambiguateItem) (string, error) {
	p := tea.NewProgram(newDisambiguatePicker(items))
	finalModel, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("disambiguation picker: %w", err)
	}
	picker, ok := finalModel.(disambiguatePicker)
	if !ok {
		return "", fmt.Errorf("unexpected model type from picker")
	}
	id := picker.selectedID()
	if id == "" {
		return "", nil // cancelled
	}
	return id, nil
}

// resolveIDInteractive resolves an object ID prefix, showing an interactive
// picker if the prefix is ambiguous and stdin is a terminal.
func resolveIDInteractive(vault *core.Vault, prefix string) (string, error) {
	id, err := vault.ResolveID(prefix)
	if err == nil {
		return id, nil
	}

	var ambErr *core.AmbiguousMatchError
	if !errors.As(err, &ambErr) {
		return "", err
	}

	if !isInteractiveFunc() {
		return "", err
	}

	items := buildDisambiguateItems(vault, ambErr.Matches)
	selected, pickErr := disambiguateFunc(items)
	if pickErr != nil {
		return "", pickErr
	}
	if selected == "" {
		return "", err // user cancelled — return original ambiguous error
	}
	return selected, nil
}

// resolveObjectInteractive resolves a prefix and returns the full object,
// showing an interactive picker if the prefix is ambiguous.
func resolveObjectInteractive(vault *core.Vault, prefix string) (*core.Object, error) {
	id, err := resolveIDInteractive(vault, prefix)
	if err != nil {
		return nil, err
	}
	return vault.GetObject(id)
}

