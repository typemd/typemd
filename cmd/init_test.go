package cmd

import (
	"testing"
)

func TestInitCmd_NoStartersFlag(t *testing.T) {
	f := initCmd.Flags().Lookup("no-starters")
	if f == nil {
		t.Fatal("expected --no-starters flag to exist")
	}
	if f.DefValue != "false" {
		t.Errorf("expected default false, got %q", f.DefValue)
	}
}

