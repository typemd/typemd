package core

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// GraphOptions configures DOT graph export.
type GraphOptions struct {
	Types       []string // Filter to these object types (empty = all)
	NoRelations bool     // Exclude relation edges
	NoWikiLinks bool     // Exclude wiki-link edges
}

// ExportDOT writes the vault's object graph in DOT format to the given writer.
// Nodes represent objects; edges represent relations and wiki-links.
func (v *Vault) ExportDOT(w io.Writer, opts GraphOptions) error {
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	// Build type filter.
	var filter []FilterRule
	if len(opts.Types) == 1 {
		filter = TypeFilter(opts.Types[0])
	} else if len(opts.Types) > 1 {
		for _, t := range opts.Types {
			filter = append(filter, FilterRule{Property: "type", Operator: "is", Value: t})
		}
	}

	objects, err := v.QueryObjects(filter)
	if err != nil {
		return fmt.Errorf("query objects: %w", err)
	}

	// Build a set of object IDs in the graph for edge filtering.
	idSet := make(map[string]bool, len(objects))
	for _, obj := range objects {
		idSet[obj.ID] = true
	}

	fmt.Fprintln(bw, "digraph vault {")
	fmt.Fprintln(bw, "  rankdir=LR;")
	fmt.Fprintln(bw, "  node [shape=box];")

	// Write nodes.
	for _, obj := range objects {
		label := nodeLabel(v, obj)
		fmt.Fprintf(bw, "  %s [label=%s];\n", dotID(obj.ID), dotQuote(label))
	}

	// Track seen edges to deduplicate bidirectional relations.
	type edgeKey struct{ from, to, label string }
	seen := make(map[edgeKey]bool)

	// Write relation edges.
	if !opts.NoRelations {
		for _, obj := range objects {
			relations, err := v.ListRelations(obj.ID)
			if err != nil {
				return fmt.Errorf("list relations for %s: %w", obj.ID, err)
			}
			for _, rel := range relations {
				// Only emit edges where this object is the source.
				if rel.FromID != obj.ID {
					continue
				}
				// Skip edges to objects outside the filtered set.
				if !idSet[rel.ToID] {
					continue
				}
				key := edgeKey{rel.FromID, rel.ToID, rel.Name}
				reverseKey := edgeKey{rel.ToID, rel.FromID, rel.Name}
				if seen[key] || seen[reverseKey] {
					continue
				}
				seen[key] = true
				fmt.Fprintf(bw, "  %s -> %s [label=%s];\n",
					dotID(rel.FromID), dotID(rel.ToID), dotQuote(rel.Name))
			}
		}
	}

	// Write wiki-link edges.
	if !opts.NoWikiLinks {
		for _, obj := range objects {
			links, err := v.ListWikiLinks(obj.ID)
			if err != nil {
				return fmt.Errorf("list wikilinks for %s: %w", obj.ID, err)
			}
			for _, link := range links {
				// Skip unresolved links.
				if link.ToID == "" {
					continue
				}
				// Skip edges to objects outside the filtered set.
				if !idSet[link.ToID] {
					continue
				}
				fmt.Fprintf(bw, "  %s -> %s [label=%s style=dashed];\n",
					dotID(link.FromID), dotID(link.ToID), dotQuote("wikilink"))
			}
		}
	}

	fmt.Fprintln(bw, "}")
	return nil
}

// nodeLabel returns the display label for a node: "emoji name" or just "name".
func nodeLabel(v *Vault, obj *Object) string {
	name := obj.GetName()
	schema, err := v.LoadType(obj.Type)
	if err != nil || schema.Emoji == "" {
		return name
	}
	return schema.Emoji + " " + name
}

// dotID wraps an object ID as a DOT identifier.
func dotID(id string) string {
	return `"` + id + `"`
}

// dotQuote escapes a string for use as a DOT label value.
func dotQuote(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return `"` + s + `"`
}
