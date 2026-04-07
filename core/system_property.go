package core

// System property name constants — stored (frontmatter).
const (
	DescriptionProperty = "description"
	CreatedAtProperty   = "created_at"
	UpdatedAtProperty   = "updated_at"
	TagsProperty        = "tags"
	LockedProperty      = "locked"
	ArchivedProperty    = "archived"
)

// System property name constants — non-stored (runtime, not in frontmatter).
const (
	ObjectTypeProperty = "object_type"
	LinksProperty      = "links"
	BacklinksProperty  = "backlinks"
	CreatedByProperty  = "created_by"
	UpdatedByProperty  = "updated_by"
)

// Built-in type name constants.
const (
	TagTypeName  = "tag"
	PageTypeName = "page"
)

// SystemProperty defines a system-managed property that is automatically
// present on all objects, regardless of type schema.
type SystemProperty struct {
	Name      string
	Type      string
	Target    string // only for relation type
	Multiple  bool   // only for relation type
	Immutable bool
	Derived   bool
	Computed  bool
}

// systemProperties is the authoritative registry of all system properties.
// Stored properties come first (order determines frontmatter output ordering),
// followed by computed properties (never written to frontmatter).
var systemProperties = []SystemProperty{
	// Stored (frontmatter)
	{Name: NameProperty, Type: "text"},
	{Name: DescriptionProperty, Type: "text"},
	{Name: CreatedAtProperty, Type: "datetime", Immutable: true},
	{Name: UpdatedAtProperty, Type: "datetime", Immutable: true},
	{Name: TagsProperty, Type: "relation", Target: TagTypeName, Multiple: true},
	{Name: LockedProperty, Type: "checkbox"},
	{Name: ArchivedProperty, Type: "checkbox"},
	// Derived (not stored in frontmatter)
	{Name: ObjectTypeProperty, Type: "text", Immutable: true, Derived: true},
	{Name: CreatedByProperty, Type: "text", Immutable: true, Derived: true},
	// Computed (not stored in frontmatter)
	{Name: LinksProperty, Type: "text", Immutable: true, Computed: true},
	{Name: BacklinksProperty, Type: "text", Immutable: true, Computed: true},
	{Name: UpdatedByProperty, Type: "text", Immutable: true, Computed: true},
}

// IsSystemProperty returns true if the given name is a system property.
func IsSystemProperty(name string) bool {
	return lookupSystemProperty(name) != nil
}

// IsImmutableSystemProperty returns true if the given name is an immutable
// (auto-managed) system property that cannot be overridden by templates.
func IsImmutableSystemProperty(name string) bool {
	sp := lookupSystemProperty(name)
	return sp != nil && sp.Immutable
}

// lookupSystemProperty returns the SystemProperty for a given name, or nil if not found.
func lookupSystemProperty(name string) *SystemProperty {
	for i, sp := range systemProperties {
		if sp.Name == name {
			return &systemProperties[i]
		}
	}
	return nil
}

// IsNonStoredProperty returns true if the given name is a derived or computed
// system property (not stored in frontmatter).
func IsNonStoredProperty(name string) bool {
	sp := lookupSystemProperty(name)
	return sp != nil && (sp.Derived || sp.Computed)
}

// IsDerivedProperty returns true if the given name is a derived system property.
func IsDerivedProperty(name string) bool {
	sp := lookupSystemProperty(name)
	return sp != nil && sp.Derived
}

// IsComputedProperty returns true if the given name is a computed system property.
func IsComputedProperty(name string) bool {
	sp := lookupSystemProperty(name)
	return sp != nil && sp.Computed
}

// SystemPropertyNames returns all system property names in registry order.
func SystemPropertyNames() []string {
	names := make([]string, len(systemProperties))
	for i, sp := range systemProperties {
		names[i] = sp.Name
	}
	return names
}

// StoredPropertyNames returns only stored (non-derived, non-computed) system property names in registry order.
func StoredPropertyNames() []string {
	names := make([]string, 0, len(systemProperties))
	for _, sp := range systemProperties {
		if !sp.Derived && !sp.Computed {
			names = append(names, sp.Name)
		}
	}
	return names
}
