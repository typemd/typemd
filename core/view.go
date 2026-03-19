package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ViewLayout defines the presentation layout for a view.
type ViewLayout string

const (
	// ViewLayoutList displays objects as a list.
	ViewLayoutList ViewLayout = "list"
)

// GroupRule defines a single grouping level.
type GroupRule struct {
	Property string `yaml:"property"`
}

// ViewConfig defines a saved view configuration for a type.
type ViewConfig struct {
	Name    string       `yaml:"name"`
	Layout  ViewLayout   `yaml:"layout"`
	Filter  []FilterRule `yaml:"filter,omitempty"`
	Sort    []SortRule   `yaml:"sort,omitempty"`
	GroupBy []GroupRule   `yaml:"group_by,omitempty"`
}

// UnmarshalYAML implements custom YAML unmarshaling for ViewConfig to support
// both legacy string format (group_by: "genre") and new array format
// (group_by: [{property: genre}]).
func (vc *ViewConfig) UnmarshalYAML(value *yaml.Node) error {
	// Use a raw struct to avoid recursion.
	var raw struct {
		Name    string       `yaml:"name"`
		Layout  ViewLayout   `yaml:"layout"`
		Filter  []FilterRule `yaml:"filter,omitempty"`
		Sort    []SortRule   `yaml:"sort,omitempty"`
		GroupBy yaml.Node    `yaml:"group_by"`
	}
	if err := value.Decode(&raw); err != nil {
		return err
	}

	vc.Name = raw.Name
	vc.Layout = raw.Layout
	vc.Filter = raw.Filter
	vc.Sort = raw.Sort

	// Parse group_by based on YAML node kind.
	switch raw.GroupBy.Kind {
	case yaml.ScalarNode:
		// Legacy string format: group_by: "genre"
		s := raw.GroupBy.Value
		if s != "" {
			vc.GroupBy = []GroupRule{{Property: s}}
		}
	case yaml.SequenceNode:
		// New array format: group_by: [{property: genre}]
		var rules []GroupRule
		if err := raw.GroupBy.Decode(&rules); err != nil {
			return fmt.Errorf("decode group_by array: %w", err)
		}
		vc.GroupBy = rules
	case 0:
		// Not present in YAML — leave nil.
	}

	return nil
}

// FilterRule defines a single filter condition.
type FilterRule struct {
	Property string `yaml:"property"`
	Operator string `yaml:"operator"`
	Value    string `yaml:"value,omitempty"`
}

// TypeFilter returns a []FilterRule that matches objects of the given type.
func TypeFilter(typeName string) []FilterRule {
	return []FilterRule{{Property: "type", Operator: "is", Value: typeName}}
}

// viewsDir returns the views directory path for a type.
func (v *Vault) viewsDir(typeName string) string {
	return filepath.Join(v.TypesDir(), typeName, "views")
}

// ListViews returns all saved ViewConfig objects for the given type.
func (v *Vault) ListViews(typeName string) ([]ViewConfig, error) {
	dir := v.viewsDir(typeName)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read views directory: %w", err)
	}

	var views []ViewConfig
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name, ok := strings.CutSuffix(e.Name(), ".yaml")
		if !ok {
			continue
		}
		vc, err := v.LoadView(typeName, name)
		if err != nil {
			continue
		}
		views = append(views, *vc)
	}
	return views, nil
}

// LoadView reads and parses a specific view by name.
func (v *Vault) LoadView(typeName, viewName string) (*ViewConfig, error) {
	path := filepath.Join(v.viewsDir(typeName), viewName+".yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("view %q not found for type %q", viewName, typeName)
		}
		return nil, fmt.Errorf("read view: %w", err)
	}

	var vc ViewConfig
	if err := yaml.Unmarshal(data, &vc); err != nil {
		return nil, fmt.Errorf("parse view: %w", err)
	}
	return &vc, nil
}

// SaveView validates and writes a ViewConfig to disk.
func (v *Vault) SaveView(typeName string, view *ViewConfig) error {
	if view.Name == "" || strings.ContainsAny(view.Name, "/\\") || strings.Contains(view.Name, "..") {
		return fmt.Errorf("invalid view name %q", view.Name)
	}

	dir := v.viewsDir(typeName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create views directory: %w", err)
	}

	data, err := yaml.Marshal(view)
	if err != nil {
		return fmt.Errorf("marshal view: %w", err)
	}

	path := filepath.Join(dir, view.Name+".yaml")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write view: %w", err)
	}
	return nil
}

// DeleteView removes a view file. Cleans up empty views directory.
func (v *Vault) DeleteView(typeName, viewName string) error {
	dir := v.viewsDir(typeName)
	path := filepath.Join(dir, viewName+".yaml")
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("view %q not found for type %q", viewName, typeName)
		}
		return fmt.Errorf("delete view: %w", err)
	}

	// Clean up empty views directory
	entries, err := os.ReadDir(dir)
	if err == nil && len(entries) == 0 {
		os.Remove(dir)
	}

	return nil
}

// DefaultView returns the default view for a type.
// If a saved default.yaml exists, it returns that; otherwise returns an implicit default.
func (v *Vault) DefaultView(typeName string) *ViewConfig {
	vc, err := v.LoadView(typeName, "default")
	if err == nil {
		return vc
	}
	return &ViewConfig{
		Name:   "default",
		Layout: ViewLayoutList,
		Sort:   []SortRule{{Property: "name", Direction: "asc"}},
	}
}
