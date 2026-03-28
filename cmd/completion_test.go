package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/typemd/typemd/core"
)

func setupTestVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	v := core.NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("init vault: %v", err)
	}
	// Write a book type schema (directory format)
	os.MkdirAll(v.TypesDir()+"/book", 0755)
	os.WriteFile(v.TypesDir()+"/book/schema.yaml", []byte(`name: book
emoji: "\U0001F4DA"
plural: books
properties:
  - name: author
    type: string
`), 0644)
	// Create an object file
	objDir := v.ObjectDir("book")
	os.MkdirAll(objDir, 0755)
	os.WriteFile(objDir+"/clean-code-01abc123.md", []byte("---\nname: Clean Code\n---\n"), 0644)
	return dir
}

func TestCompleteObjectID_EmptyVault(t *testing.T) {
	dir := t.TempDir()
	v := core.NewVault(dir)
	v.Init()

	vaultPath = dir
	completions, directive := completeObjectID("")
	if len(completions) == 0 {
		// Only built-in types (tag, page) should appear
		// but they don't have objects dirs, so type completion still works
	}
	if directive&cobra.ShellCompDirectiveNoFileComp == 0 {
		t.Error("expected NoFileComp directive")
	}
}

func TestCompleteObjectID_TypeStage(t *testing.T) {
	dir := setupTestVault(t)
	vaultPath = dir

	completions, directive := completeObjectID("bo")
	if len(completions) != 1 || completions[0] != "book/" {
		t.Errorf("expected [book/], got %v", completions)
	}
	if directive&cobra.ShellCompDirectiveNoSpace == 0 {
		t.Error("expected NoSpace directive for type completion")
	}
	if directive&cobra.ShellCompDirectiveNoFileComp == 0 {
		t.Error("expected NoFileComp directive")
	}
}

func TestCompleteObjectID_ObjectStage(t *testing.T) {
	dir := setupTestVault(t)
	vaultPath = dir

	completions, directive := completeObjectID("book/clean")
	if len(completions) != 1 || completions[0] != "book/clean-code-01abc123" {
		t.Errorf("expected [book/clean-code-01abc123], got %v", completions)
	}
	if directive&cobra.ShellCompDirectiveNoFileComp == 0 {
		t.Error("expected NoFileComp directive")
	}
}

func TestCompleteObjectID_NoMatch(t *testing.T) {
	dir := setupTestVault(t)
	vaultPath = dir

	completions, _ := completeObjectID("zzz")
	if len(completions) != 0 {
		t.Errorf("expected no completions, got %v", completions)
	}
}

func TestCompleteTypeName_Prefix(t *testing.T) {
	dir := setupTestVault(t)
	vaultPath = dir

	completions, directive := completeTypeName("bo")
	if len(completions) != 1 || completions[0] != "book" {
		t.Errorf("expected [book], got %v", completions)
	}
	if directive&cobra.ShellCompDirectiveNoFileComp == 0 {
		t.Error("expected NoFileComp directive")
	}
}

func TestCompleteTypeName_Empty(t *testing.T) {
	dir := setupTestVault(t)
	vaultPath = dir

	completions, _ := completeTypeName("")
	// Should include at least book, tag, page (built-in types)
	found := make(map[string]bool)
	for _, c := range completions {
		found[c] = true
	}
	for _, expected := range []string{"book", "tag", "page"} {
		if !found[expected] {
			t.Errorf("expected %q in completions, got %v", expected, completions)
		}
	}
}

func TestCompleteRelationName_NoRelations(t *testing.T) {
	dir := setupTestVault(t)
	vaultPath = dir

	completions, _ := completeRelationName("book/clean-code-01abc123", "")
	if len(completions) != 0 {
		t.Errorf("expected no completions (book has no relations), got %v", completions)
	}
}

func TestCompleteRelationName_InvalidID(t *testing.T) {
	dir := setupTestVault(t)
	vaultPath = dir

	completions, directive := completeRelationName("invalid", "")
	if len(completions) != 0 {
		t.Errorf("expected no completions, got %v", completions)
	}
	if directive&cobra.ShellCompDirectiveNoFileComp == 0 {
		t.Error("expected NoFileComp directive")
	}
}

func TestCompleteObjectID_VaultPathFlag(t *testing.T) {
	dir := setupTestVault(t)

	// Simulate --vault flag being set
	vaultPath = dir
	completions, _ := completeObjectID("book/")
	if len(completions) != 1 {
		t.Errorf("expected 1 completion with vault flag, got %v", completions)
	}
}

func TestCompletionSubcommand(t *testing.T) {
	t.Run("bash", func(t *testing.T) {
		var buf bytes.Buffer
		if err := rootCmd.GenBashCompletionV2(&buf, true); err != nil {
			t.Fatalf("GenBashCompletionV2 failed: %v", err)
		}
		if buf.Len() == 0 {
			t.Error("bash completion produced no output")
		}
	})

	t.Run("zsh", func(t *testing.T) {
		var buf bytes.Buffer
		if err := rootCmd.GenZshCompletion(&buf); err != nil {
			t.Fatalf("GenZshCompletion failed: %v", err)
		}
		if buf.Len() == 0 {
			t.Error("zsh completion produced no output")
		}
	})

	t.Run("fish", func(t *testing.T) {
		var buf bytes.Buffer
		if err := rootCmd.GenFishCompletion(&buf, true); err != nil {
			t.Fatalf("GenFishCompletion failed: %v", err)
		}
		if buf.Len() == 0 {
			t.Error("fish completion produced no output")
		}
	})
}

