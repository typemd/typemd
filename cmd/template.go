package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/typemd/typemd/core"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage object templates (list, show, create, delete)",
}

// parseTemplateArg splits a "type/name" argument into type and name parts.
func parseTemplateArg(arg string) (typeName, name string, err error) {
	parts := strings.SplitN(arg, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid template identifier %q: expected type/name format", arg)
	}
	return parts[0], parts[1], nil
}

// templateListEntry represents a template for JSON output.
type templateListEntry struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

var templateListCmd = &cobra.Command{
	Use:   "list [type]",
	Short: "List available templates",
	Long: `List all available templates, optionally filtered by type.

Examples:
  tmd template list
  tmd template list book
  tmd template list --json`,
	Args: cobra.MaximumNArgs(1),
	ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return completeTemplateTypes()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath)
		if err != nil {
			return err
		}
		defer vault.Close()

		jsonOutput, _ := cmd.Flags().GetBool("json")

		var entries []templateListEntry

		if len(args) == 1 {
			// Filter by specific type
			names, err := vault.ListTemplates(args[0])
			if err != nil {
				return err
			}
			for _, n := range names {
				entries = append(entries, templateListEntry{Type: args[0], Name: n})
			}
		} else {
			// All types — scan templates directory for type subdirectories
			types := listTemplateTypes(vault)
			for _, typeName := range types {
				names, err := vault.ListTemplates(typeName)
				if err != nil {
					continue
				}
				for _, n := range names {
					entries = append(entries, templateListEntry{Type: typeName, Name: n})
				}
			}
		}

		if jsonOutput {
			if entries == nil {
				entries = []templateListEntry{}
			}
			return printJSON(entries)
		}

		for _, e := range entries {
			fmt.Printf("%s/%s\n", e.Type, e.Name)
		}
		return nil
	},
}

var templateShowCmd = &cobra.Command{
	Use:   "show <type/name>",
	Short: "Show a template's content",
	Long: `Display a template's frontmatter properties and body content.

Examples:
  tmd template show book/review
  tmd template show note/meeting`,
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return completeTemplateIDs()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		typeName, name, err := parseTemplateArg(args[0])
		if err != nil {
			return err
		}

		vault, err := openVault(vaultPath)
		if err != nil {
			return err
		}
		defer vault.Close()

		tmpl, err := vault.LoadTemplate(typeName, name)
		if err != nil {
			return err
		}

		// Title
		fmt.Printf("%s/%s\n\n", typeName, name)

		// Properties
		fmt.Println("Properties")
		fmt.Println("──────────")
		if len(tmpl.Properties) == 0 {
			fmt.Println("  (none)")
		} else {
			keys := make([]string, 0, len(tmpl.Properties))
			for k := range tmpl.Properties {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Printf("  %s: %v\n", k, tmpl.Properties[k])
			}
		}

		// Body
		fmt.Println()
		fmt.Println("Body")
		fmt.Println("────")
		body := strings.TrimSpace(tmpl.Body)
		if body == "" {
			fmt.Println("  (empty)")
		} else {
			for _, line := range strings.Split(body, "\n") {
				fmt.Printf("  %s\n", line)
			}
		}

		return nil
	},
}

// openEditorFunc is the function used to open a file in the user's editor.
// Override in tests to skip actual editor launch.
var openEditorFunc = openEditor

var templateCreateCmd = &cobra.Command{
	Use:   "create <type/name>",
	Short: "Create a new template file",
	Long: `Create a new template file and open it in your editor.

The editor is resolved from $EDITOR, $VISUAL, or defaults to vi.

Examples:
  tmd template create book/review
  tmd template create note/meeting`,
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return completeTemplateTypes()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		typeName, name, err := parseTemplateArg(args[0])
		if err != nil {
			return err
		}

		vault, err := openVault(vaultPath)
		if err != nil {
			return err
		}
		defer vault.Close()

		// Check if template already exists
		if _, statErr := os.Stat(vault.TemplatePath(typeName, name)); statErr == nil {
			return fmt.Errorf("template %s/%s already exists", typeName, name)
		}

		tmpl := &core.Template{
			Name:       name,
			Properties: map[string]any{},
			Body:       "",
		}
		if err := vault.SaveTemplate(typeName, name, tmpl); err != nil {
			return err
		}

		fmt.Printf("Created template %s/%s\n", typeName, name)

		filePath := vault.TemplatePath(typeName, name)
		return openEditorFunc(filePath)
	},
}

var templateDeleteCmd = &cobra.Command{
	Use:   "delete <type/name>",
	Short: "Delete a template file",
	Long: `Delete a template file from the vault.

In interactive terminals, prompts for confirmation unless --force is used.

Examples:
  tmd template delete book/review
  tmd template delete book/review --force`,
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return completeTemplateIDs()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		typeName, name, err := parseTemplateArg(args[0])
		if err != nil {
			return err
		}

		vault, err := openVault(vaultPath)
		if err != nil {
			return err
		}
		defer vault.Close()

		force, _ := cmd.Flags().GetBool("force")
		if !force && isInteractiveFunc() {
			fmt.Printf("Delete template %s/%s? [y/N] ", typeName, name)
			scanner := bufio.NewScanner(os.Stdin)
			if !scanner.Scan() {
				return fmt.Errorf("no input received")
			}
			answer := strings.ToLower(strings.TrimSpace(scanner.Text()))
			if answer != "y" && answer != "yes" {
				fmt.Println("Cancelled.")
				return nil
			}
		}

		if err := vault.DeleteTemplate(typeName, name); err != nil {
			return err
		}

		fmt.Printf("Deleted template %s/%s\n", typeName, name)
		return nil
	},
}

// resolveEditor returns the editor command from environment variables.
func resolveEditor() (string, error) {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor, nil
	}
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor, nil
	}
	if _, err := exec.LookPath("vi"); err == nil {
		return "vi", nil
	}
	return "", fmt.Errorf("no editor found: set $EDITOR or $VISUAL")
}

// openEditor opens a file in the user's editor.
func openEditor(filePath string) error {
	editor, err := resolveEditor()
	if err != nil {
		return err
	}
	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// listTemplateTypes returns type names that have template directories,
// by scanning the templates/ directory. This discovers templates for types
// that may not have schemas yet.
func listTemplateTypes(vault *core.Vault) []string {
	entries, err := os.ReadDir(vault.TemplatesDir())
	if err != nil {
		return nil
	}
	var types []string
	for _, e := range entries {
		if e.IsDir() {
			types = append(types, e.Name())
		}
	}
	sort.Strings(types)
	return types
}

// completeTemplateTypes returns type names that have templates, for shell completion.
func completeTemplateTypes() ([]string, cobra.ShellCompDirective) {
	vault, err := openVault(vaultPath)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	defer vault.Close()

	return listTemplateTypes(vault), cobra.ShellCompDirectiveNoFileComp
}

// completeTemplateIDs returns type/name pairs for shell completion.
func completeTemplateIDs() ([]string, cobra.ShellCompDirective) {
	vault, err := openVault(vaultPath)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	defer vault.Close()

	var ids []string
	for _, typeName := range listTemplateTypes(vault) {
		names, err := vault.ListTemplates(typeName)
		if err != nil {
			continue
		}
		for _, n := range names {
			ids = append(ids, typeName+"/"+n)
		}
	}
	return ids, cobra.ShellCompDirectiveNoFileComp
}

func init() {
	templateListCmd.Flags().Bool("json", false, "Output results as JSON")
	templateDeleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")

	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateShowCmd)
	templateCmd.AddCommand(templateCreateCmd)
	templateCmd.AddCommand(templateDeleteCmd)
	rootCmd.AddCommand(templateCmd)
}
