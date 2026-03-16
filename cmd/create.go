package cmd

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	templateFlag string
	typeFlag     string
)

var createCmd = &cobra.Command{
	Use:   "create [type] [name]",
	Short: "Create a new object from a type schema",
	Long: `Create a new object file (Markdown + YAML frontmatter) based on the type schema.
A ULID is automatically appended to the filename for uniqueness.
If the type has a name template, the name argument is optional.
If the type has templates, they are applied automatically (single) or via selection (multiple).

The type argument can be omitted if --type flag is set or cli.default_type is
configured in .typemd/config.yaml. Names are automatically converted to slugs
(e.g., "Some Thought" becomes "some-thought" in the filename).

When a single argument matches a known type name, it is treated as the type.
To use a name that collides with a type name, use the --type flag explicitly.

Examples:
  tmd object create book "Clean Code"
  tmd object create book "Clean Code" -t review
  tmd object create --type idea "Some Thought"
  tmd object create "Some Thought"              # uses default type from config
  tmd object create book`,
	Args: cobra.RangeArgs(0, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath, reindex)
		if err != nil {
			return err
		}
		defer vault.Close()

		typeName, name, err := resolveTypeAndName(args, typeFlag, vault.DefaultType(), vault.ListTypes)
		if err != nil {
			return err
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

// resolveTypeAndName determines the type and name from positional args,
// the --type flag, and the default type from config. listTypes is called
// lazily only when needed to check if an argument is a known type.
func resolveTypeAndName(args []string, typeFlag, defaultType string, listTypes func() []string) (typeName, name string, err error) {
	// If --type flag is set, all positional args are the name
	if typeFlag != "" {
		typeName = typeFlag
		if len(args) > 0 {
			name = strings.Join(args, " ")
		}
		return typeName, name, nil
	}

	switch len(args) {
	case 0:
		// No args: type from config default
		if defaultType == "" {
			return "", "", fmt.Errorf("type is required: provide as argument, use --type flag, or set cli.default_type in .typemd/config.yaml")
		}
		typeName = defaultType
	case 1:
		// 1 arg: check if it's a known type
		if isKnownType(args[0], listTypes()) {
			typeName = args[0]
		} else if defaultType != "" {
			// Not a type → treat as name with default type
			typeName = defaultType
			name = args[0]
		} else {
			return "", "", fmt.Errorf("unknown type %q (no default type configured)", args[0])
		}
	case 2:
		// 2 args: first is type, second is name (backward compatible)
		typeName = args[0]
		name = args[1]
	}

	return typeName, name, nil
}

func isKnownType(name string, knownTypes []string) bool {
	return slices.Contains(knownTypes, name)
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
	createCmd.Flags().StringVar(&typeFlag, "type", "", "object type (overrides config default)")
	objectCmd.AddCommand(createCmd)
}
