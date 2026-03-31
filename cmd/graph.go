package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/typemd/typemd/core"
)

var graphTypes []string
var graphNoRelations bool
var graphNoWikiLinks bool

var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Export the object relation graph in DOT format",
	Long: `Export the object relation graph in DOT format.

Outputs a Graphviz DOT digraph to stdout. Nodes represent objects,
edges represent relations and wiki-links. Pipe to dot for rendering.

Examples:
  tmd graph
  tmd graph --type book --type person
  tmd graph --no-wikilinks
  tmd graph | dot -Tpng -o graph.png`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath)
		if err != nil {
			return err
		}
		defer vault.Close()

		opts := core.GraphOptions{
			Types:       graphTypes,
			NoRelations: graphNoRelations,
			NoWikiLinks: graphNoWikiLinks,
		}

		return vault.ExportDOT(os.Stdout, opts)
	},
}

func init() {
	graphCmd.Flags().StringArrayVar(&graphTypes, "type", nil, "Filter to specific object types (can be repeated)")
	graphCmd.Flags().BoolVar(&graphNoRelations, "no-relations", false, "Exclude relation edges")
	graphCmd.Flags().BoolVar(&graphNoWikiLinks, "no-wikilinks", false, "Exclude wiki-link edges")
	graphCmd.RegisterFlagCompletionFunc("type", func(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completeTypeName(toComplete)
	})
	rootCmd.AddCommand(graphCmd)
}
