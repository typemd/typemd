package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/typemd/typemd/core"
)

// completionVault creates a Vault for shell completion.
// Uses resolveVault (no Open/SQLite) since completion only needs filesystem access.
func completionVault() *core.Vault {
	return resolveVault(vaultPath)
}

// completeObjectID provides two-stage progressive completion for object IDs.
// Stage 1: no "/" in toComplete → complete type names with trailing "/"
// Stage 2: "/" present → complete object names within that type
func completeObjectID(toComplete string) ([]string, cobra.ShellCompDirective) {
	vault := completionVault()

	if !strings.Contains(toComplete, "/") {
		types := vault.ListTypes()
		var completions []string
		for _, t := range types {
			if strings.HasPrefix(t, toComplete) {
				completions = append(completions, t+"/")
			}
		}
		return completions, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	}

	parts := strings.SplitN(toComplete, "/", 2)
	typeName, namePrefix := parts[0], parts[1]

	ids, err := vault.GlobObjectIDs(typeName, namePrefix)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return ids, cobra.ShellCompDirectiveNoFileComp
}

// completeTypeName provides completion for type names.
func completeTypeName(toComplete string) ([]string, cobra.ShellCompDirective) {
	vault := completionVault()
	types := vault.ListTypes()
	var completions []string
	for _, t := range types {
		if strings.HasPrefix(t, toComplete) {
			completions = append(completions, t)
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// completeRelationName provides completion for relation property names
// based on the source object's type schema.
func completeRelationName(fromID string, toComplete string) ([]string, cobra.ShellCompDirective) {
	vault := completionVault()

	parts := strings.SplitN(fromID, "/", 2)
	if len(parts) != 2 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	typeName := parts[0]

	schema, err := vault.LoadType(typeName)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var completions []string
	for _, p := range schema.Properties {
		if p.Type == "relation" && strings.HasPrefix(p.Name, toComplete) {
			completions = append(completions, p.Name)
		}
	}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// completeLinkArgs provides positional completion for link/unlink commands:
// arg 0 = object ID, arg 1 = relation name, arg 2 = object ID.
func completeLinkArgs(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	switch len(args) {
	case 0:
		return completeObjectID(toComplete)
	case 1:
		return completeRelationName(args[0], toComplete)
	case 2:
		return completeObjectID(toComplete)
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}
