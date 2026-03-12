package core

// System property name constants.
const (
	DescriptionProperty = "description"
	CreatedAtProperty   = "created_at"
	UpdatedAtProperty   = "updated_at"
)

// SystemProperty defines a system-managed property that is automatically
// present on all objects, regardless of type schema.
type SystemProperty struct {
	Name string
	Type string
}

// systemProperties is the authoritative registry of all system properties.
// Order matters: it determines frontmatter output ordering.
var systemProperties = []SystemProperty{
	{Name: NameProperty, Type: "text"},
	{Name: DescriptionProperty, Type: "text"},
	{Name: CreatedAtProperty, Type: "datetime"},
	{Name: UpdatedAtProperty, Type: "datetime"},
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

// SystemPropertyNames returns all system property names in registry order.
func SystemPropertyNames() []string {
	names := make([]string, len(systemProperties))
	for i, sp := range systemProperties {
		names[i] = sp.Name
	}
	return names
}
