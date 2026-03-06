package cmd

import (
	"testing"
)

func TestResolveVault(t *testing.T) {
	t.Run("empty path defaults to current directory", func(t *testing.T) {
		v := resolveVault("")
		if v.Root != "." {
			t.Errorf("expected root \".\", got %q", v.Root)
		}
	})

	t.Run("explicit path is used", func(t *testing.T) {
		v := resolveVault("/tmp/testvault")
		if v.Root != "/tmp/testvault" {
			t.Errorf("expected root \"/tmp/testvault\", got %q", v.Root)
		}
	})
}
