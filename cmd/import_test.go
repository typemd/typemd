package cmd

import (
	"testing"
)

func TestImportCmdRegistered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "import" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("import command not registered on rootCmd")
	}
}

func TestImportSubcommands(t *testing.T) {
	subcommands := map[string]bool{
		"scan":    false,
		"plan":    false,
		"execute": false,
	}

	for _, cmd := range importCmd.Commands() {
		name := cmd.Name()
		if _, ok := subcommands[name]; ok {
			subcommands[name] = true
		}
	}

	for name, found := range subcommands {
		if !found {
			t.Errorf("subcommand %q not registered on importCmd", name)
		}
	}
}

func TestImportPlanOutputFlag(t *testing.T) {
	f := importPlanCmd.Flags().Lookup("output")
	if f == nil {
		t.Fatal("--output flag not defined on import plan command")
	}
	if f.Shorthand != "o" {
		t.Errorf("expected shorthand 'o', got %q", f.Shorthand)
	}
}
