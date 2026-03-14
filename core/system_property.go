package core

// System property name constants.
const (
	DescriptionProperty = "description"
	CreatedAtProperty   = "created_at"
	UpdatedAtProperty   = "updated_at"
	TagsProperty        = "tags"
)

// Built-in type name constants.
const (
	TagTypeName = "tag"
)

// SystemProperty defines a system-managed property that is automatically
// present on all objects, regardless of type schema.
type SystemProperty struct {
	Name      string
	Type      string
	Target    string // only for relation type
	Multiple  bool   // only for relation type
	Immutable bool   // true for auto-managed properties (created_at, updated_at)
}

// systemProperties is the authoritative registry of all system properties.
// Order matters: it determines frontmatter output ordering.
var systemProperties = []SystemProperty{
	{Name: NameProperty, Type: "text"},
	{Name: DescriptionProperty, Type: "text"},
	{Name: CreatedAtProperty, Type: "datetime", Immutable: true},
	{Name: UpdatedAtProperty, Type: "datetime", Immutable: true},
	{Name: TagsProperty, Type: "relation", Target: TagTypeName, Multiple: true},
}

// IsSystemProperty returns true if the given name is a system property.
func IsSystemProperty(name string) bool {
	for _, sp := range systemProperties {
		if sp.Name == name {
			return true
		}
	}
	return false
}

// IsImmutableSystemProperty returns true if the given name is an immutable
// (auto-managed) system property that cannot be overridden by templates.
func IsImmutableSystemProperty(name string) bool {
	for _, sp := range systemProperties {
		if sp.Name == name {
			return sp.Immutable
		}
	}
	return false
}

// SystemPropertyNames returns all system property names in registry order.
func SystemPropertyNames() []string {
	names := make([]string, len(systemProperties))
	for i, sp := range systemProperties {
		names[i] = sp.Name
	}
	return names
}
