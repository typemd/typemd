package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var templateFlag string

var createCmd = &cobra.Command{
	Use:   "create <type> [name]",
	Short: "Create a new object from a type schema",
	Long: `Create a new object file (Markdown + YAML frontmatter) based on the type schema.
A ULID is automatically appended to the filename for uniqueness.
If the type has a name template, the name argument is optional.
If the type has templates, they are applied automatically (single) or via selection (multiple).

Examples:
  tmd object create book clean-code
  tmd object create book clean-code -t review
  tmd object create person robert-martin
  tmd object create journal`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath, reindex)
		if err != nil {
			return err
		}
		defer vault.Close()

		typeName := args[0]
		name := ""
		if len(args) > 1 {
			name = args[1]
		}

		// Resolve template
		templateName := templateFlag
		if templateName == "" {
			templates, err := vault.ListTemplates(typeName)
			if err != nil {
				return err
			}
			switch len(templates) {
			case 0:
				// No templates, proceed without
			case 1:
				templateName = templates[0]
			default:
				templateName, err = promptTemplateSelection(templates)
				if err != nil {
					return err
				}
			}
		}

		obj, err := vault.NewObject(typeName, name, templateName)
		if err != nil {
			return err
		}

		fmt.Printf("Created %s\n", obj.ID)
		return nil
	},
}

func promptTemplateSelection(templates []string) (string, error) {
	fmt.Println("Multiple templates available:")
	for i, name := range templates {
		fmt.Printf("  %d. %s\n", i+1, name)
	}
	fmt.Print("Select template (number): ")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return "", fmt.Errorf("no input received")
	}
	input := strings.TrimSpace(scanner.Text())
	n, err := strconv.Atoi(input)
	if err != nil || n < 1 || n > len(templates) {
		return "", fmt.Errorf("invalid selection: %s", input)
	}
	return templates[n-1], nil
}

func init() {
	createCmd.Flags().StringVarP(&templateFlag, "template", "t", "", "template name to use")
	objectCmd.AddCommand(createCmd)
}
