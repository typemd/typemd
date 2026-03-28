package core

import (
	"fmt"
	"net/url"
	"time"
)

// validatePropertyType validates the type field and related constraints of a property.
// Returns validation errors with the given property name prefix.
func validatePropertyType(prop Property, namePrefix string) []error {
	var errs []error

	if prop.Type == "" {
		errs = append(errs, fmt.Errorf("%s: missing required field: type", namePrefix))
		return errs
	}

	if prop.Type == "enum" {
		errs = append(errs, fmt.Errorf("%s: type \"enum\" is no longer supported, use \"select\" with \"options\" instead", namePrefix))
		return errs
	}

	if !validPropertyTypes[prop.Type] {
		errs = append(errs, fmt.Errorf("%s: invalid type %q (valid: string, number, date, datetime, url, checkbox, select, multi_select, relation)", namePrefix, prop.Type))
		return errs
	}

	if (prop.Type == "select" || prop.Type == "multi_select") && len(prop.Options) == 0 {
		errs = append(errs, fmt.Errorf("%s: %s type requires non-empty options", namePrefix, prop.Type))
	}

	if prop.Type == "relation" && prop.Target == "" {
		errs = append(errs, fmt.Errorf("%s: relation type requires target", namePrefix))
	}

	return errs
}

// ValidateSchema validates a type schema itself for correctness.
// If sharedProps is provided, it also validates `use` entries against shared properties.
func ValidateSchema(schema *TypeSchema, sharedProps ...[]Property) []error {
	var errs []error
	if schema.Name == "" {
		errs = append(errs, fmt.Errorf("schema missing required field: name"))
	}
	if schema.Version != "" && schema.Version != DefaultSchemaVersion {
		if err := validateVersion(schema.Version); err != nil {
			errs = append(errs, err)
		}
	}
	if schema.Color != "" {
		if err := validateColor(schema.Color); err != nil {
			errs = append(errs, err)
		}
	}

	// Build shared properties map if provided
	var sharedMap map[string]Property
	if len(sharedProps) > 0 && sharedProps[0] != nil {
		sharedMap = SharedPropertiesMap(sharedProps[0])
	}

	seen := make(map[string]bool)
	seenEmoji := make(map[string]string) // emoji -> property name
	seenPin := make(map[int]string)      // pin -> property name
	for i, prop := range schema.Properties {
		// Handle `use` entries
		if prop.Use != "" {
			if prop.Name != "" {
				errs = append(errs, fmt.Errorf("property[%d]: \"use\" and \"name\" are mutually exclusive", i))
				continue
			}

			// Validate only pin, emoji, and description overrides are present
			if err := validateUseOverrides(i, prop); err != nil {
				errs = append(errs, err)
				continue
			}

			// Validate reference exists in shared properties
			if sharedMap != nil {
				shared, ok := sharedMap[prop.Use]
				if !ok {
					errs = append(errs, fmt.Errorf("property[%d]: shared property %q not found", i, prop.Use))
					continue
				}

				// Use the shared property name for duplicate checking
				propName := shared.Name
				if seen[propName] {
					errs = append(errs, fmt.Errorf("property %q: duplicate property name", propName))
				}
				seen[propName] = true

				// Check emoji uniqueness (use override or shared)
				emoji := prop.Emoji
				if emoji == "" {
					emoji = shared.Emoji
				}
				if emoji != "" {
					if otherProp, ok := seenEmoji[emoji]; ok {
						errs = append(errs, fmt.Errorf("property %q: duplicate property emoji %q (already used by %q)", propName, emoji, otherProp))
					}
					seenEmoji[emoji] = propName
				}

				// Check pin uniqueness (use override or shared)
				pin := prop.Pin
				if pin == 0 {
					pin = shared.Pin
				}
				if pin < 0 {
					errs = append(errs, fmt.Errorf("property %q: pin value must be a positive integer, got %d", propName, pin))
				} else if pin > 0 {
					if otherProp, ok := seenPin[pin]; ok {
						errs = append(errs, fmt.Errorf("property %q: duplicate pin value %d (already used by %q)", propName, pin, otherProp))
					}
					seenPin[pin] = propName
				}
			}
			continue
		}

		if prop.Name == "" {
			errs = append(errs, fmt.Errorf("property[%d]: missing required field: name", i))
			continue
		}
		if prop.Name == NameProperty {
			// Allow name entry with only template set
			onlyTemplate := prop.Template != "" &&
				prop.Type == "" && prop.Emoji == "" && prop.Description == "" && prop.Pin == 0 &&
				len(prop.Options) == 0 && prop.Target == "" && prop.Default == nil &&
				!prop.Multiple && !prop.Bidirectional && prop.Inverse == ""
			if !onlyTemplate {
				errs = append(errs, fmt.Errorf("property %q: only \"template\" is allowed on the name system property entry", prop.Name))
			}
			continue
		}
		if IsSystemProperty(prop.Name) {
			errs = append(errs, fmt.Errorf("property %q: %q is a reserved system property and cannot be defined in type schemas", prop.Name, prop.Name))
			continue
		}
		// Check if local property name conflicts with a shared property name
		if sharedMap != nil {
			if _, ok := sharedMap[prop.Name]; ok {
				errs = append(errs, fmt.Errorf("property %q: conflicts with a shared property name", prop.Name))
			}
		}
		if seen[prop.Name] {
			errs = append(errs, fmt.Errorf("property %q: duplicate property name", prop.Name))
		}
		seen[prop.Name] = true
		if prop.Emoji != "" {
			if otherProp, ok := seenEmoji[prop.Emoji]; ok {
				errs = append(errs, fmt.Errorf("property %q: duplicate property emoji %q (already used by %q)", prop.Name, prop.Emoji, otherProp))
			}
			seenEmoji[prop.Emoji] = prop.Name
		}
		if prop.Pin < 0 {
			errs = append(errs, fmt.Errorf("property %q: pin value must be a positive integer, got %d", prop.Name, prop.Pin))
		} else if prop.Pin > 0 {
			if otherProp, ok := seenPin[prop.Pin]; ok {
				errs = append(errs, fmt.Errorf("property %q: duplicate pin value %d (already used by %q)", prop.Name, prop.Pin, otherProp))
			}
			seenPin[prop.Pin] = prop.Name
		}
		typeErrs := validatePropertyType(prop, fmt.Sprintf("property %q", prop.Name))
		errs = append(errs, typeErrs...)
		if len(typeErrs) > 0 && (prop.Type == "" || prop.Type == "enum" || !validPropertyTypes[prop.Type]) {
			continue
		}
	}
	return errs
}

// validateUseOverrides checks that a `use` property entry only has allowed override fields.
func validateUseOverrides(index int, prop Property) error {
	// Only pin, emoji, and description overrides are allowed on use entries.
	disallowed := []struct {
		fieldName string
		isSet     bool
	}{
		{"type", prop.Type != ""},
		{"options", len(prop.Options) > 0},
		{"target", prop.Target != ""},
		{"default", prop.Default != nil},
		{"multiple", prop.Multiple},
		{"bidirectional", prop.Bidirectional},
		{"inverse", prop.Inverse != ""},
		{"template", prop.Template != ""},
	}

	for _, f := range disallowed {
		if f.isSet {
			return fmt.Errorf("property[%d] use %q: only \"pin\", \"emoji\", and \"description\" overrides are allowed on \"use\" entries, got %q", index, prop.Use, f.fieldName)
		}
	}
	return nil
}

// ValidateObject validates object properties against a type schema.
// Lenient mode: only validates properties defined in schema, ignores extra properties.
// Properties defined in schema but missing from props are also ignored.
func ValidateObject(props map[string]any, schema *TypeSchema) []error {
	var errs []error

	for _, prop := range schema.Properties {
		val, ok := props[prop.Name]
		if !ok || val == nil {
			continue
		}

		switch prop.Type {
		case "string":
			if _, ok := val.(string); !ok {
				errs = append(errs, fmt.Errorf("property %q: expected string, got %T", prop.Name, val))
			}
		case "number":
			switch val.(type) {
			case int, int64, float64:
				// valid
			default:
				errs = append(errs, fmt.Errorf("property %q: expected number, got %T", prop.Name, val))
			}
		case "date":
			errs = append(errs, validateDate(prop.Name, val)...)
		case "datetime":
			errs = append(errs, validateDatetime(prop.Name, val)...)
		case "url":
			errs = append(errs, validateURL(prop.Name, val)...)
		case "checkbox":
			if _, ok := val.(bool); !ok {
				errs = append(errs, fmt.Errorf("property %q: expected boolean, got %T", prop.Name, val))
			}
		case "select":
			errs = append(errs, validateSelect(prop, val)...)
		case "multi_select":
			errs = append(errs, validateMultiSelect(prop, val)...)
		case "relation":
			if prop.Multiple {
				arr, ok := val.([]any)
				if !ok {
					errs = append(errs, fmt.Errorf("property %q: expected array, got %T", prop.Name, val))
					continue
				}
				for i, item := range arr {
					if _, ok := item.(string); !ok {
						errs = append(errs, fmt.Errorf("property %q[%d]: expected string, got %T", prop.Name, i, item))
					}
				}
			} else {
				if _, ok := val.(string); !ok {
					errs = append(errs, fmt.Errorf("property %q: expected string, got %T", prop.Name, val))
				}
			}
		}
	}

	return errs
}

// datetimeFormats are the accepted string formats for datetime values.
var datetimeFormats = []string{
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05Z",
}

// validateDateInput validates a date string in YYYY-MM-DD format.
func validateDateInput(input string) error {
	if !dateRegexp.MatchString(input) {
		return fmt.Errorf("expected date in YYYY-MM-DD format, got %q", input)
	}
	if _, err := time.Parse("2006-01-02", input); err != nil {
		return fmt.Errorf("invalid date %q: %v", input, err)
	}
	return nil
}

// validateDatetimeInput validates a datetime string in ISO 8601 format.
func validateDatetimeInput(input string) error {
	for _, f := range datetimeFormats {
		if _, err := time.Parse(f, input); err == nil {
			return nil
		}
	}
	return fmt.Errorf("expected datetime in ISO 8601 format (e.g. 2006-01-02T15:04:05), got %q", input)
}

// validateURLInput validates a URL string (must have http:// or https:// scheme).
func validateURLInput(input string) error {
	u, err := url.Parse(input)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return fmt.Errorf("url must start with http:// or https://, got %q", input)
	}
	return nil
}

// validateSelectInput validates a string against allowed select options.
func validateSelectInput(options []Option, input string) error {
	for _, opt := range options {
		if opt.Value == input {
			return nil
		}
	}
	vals := make([]string, len(options))
	for i, o := range options {
		vals[i] = o.Value
	}
	return fmt.Errorf("value %q not in allowed options %v", input, vals)
}

// validateTimeProperty validates a time-based property value (date or datetime).
// Handles time.Time values from YAML auto-parsing and string values.
func validateTimeProperty(name string, val any, typeName string, validate func(string) error) []error {
	switch v := val.(type) {
	case time.Time:
		return nil
	case string:
		if err := validate(v); err != nil {
			return []error{fmt.Errorf("property %q: %s", name, err)}
		}
		return nil
	default:
		return []error{fmt.Errorf("property %q: expected %s string or time.Time, got %T", name, typeName, val)}
	}
}

// validateDate validates a date property value (YYYY-MM-DD format).
func validateDate(name string, val any) []error {
	return validateTimeProperty(name, val, "date", validateDateInput)
}

// validateDatetime validates a datetime property value (ISO 8601 with time).
func validateDatetime(name string, val any) []error {
	return validateTimeProperty(name, val, "datetime", validateDatetimeInput)
}

// validateURL validates a url property value (must have http:// or https:// scheme).
func validateURL(name string, val any) []error {
	s, ok := val.(string)
	if !ok {
		return []error{fmt.Errorf("property %q: expected string for url, got %T", name, val)}
	}
	if err := validateURLInput(s); err != nil {
		return []error{fmt.Errorf("property %q: %s", name, err)}
	}
	return nil
}

// validateSelect validates a select property value against options.
func validateSelect(prop Property, val any) []error {
	s, ok := val.(string)
	if !ok {
		return []error{fmt.Errorf("property %q: expected string for select, got %T", prop.Name, val)}
	}
	if err := validateSelectInput(prop.Options, s); err != nil {
		return []error{fmt.Errorf("property %q: %s", prop.Name, err)}
	}
	return nil
}

// validateMultiSelect validates a multi_select property value.
// Accepts a list of values (each must be in options). Coerces a single string to a list.
func validateMultiSelect(prop Property, val any) []error {
	var items []string

	switch v := val.(type) {
	case string:
		items = []string{v}
	case []any:
		for i, item := range v {
			s, ok := item.(string)
			if !ok {
				return []error{fmt.Errorf("property %q[%d]: expected string, got %T", prop.Name, i, item)}
			}
			items = append(items, s)
		}
	case []string:
		items = v
	default:
		return []error{fmt.Errorf("property %q: expected string or array for multi_select, got %T", prop.Name, val)}
	}

	allowed := make(map[string]bool, len(prop.Options))
	for _, opt := range prop.Options {
		allowed[opt.Value] = true
	}

	var errs []error
	optionVals := prop.OptionValues()
	for _, item := range items {
		if !allowed[item] {
			errs = append(errs, fmt.Errorf("property %q: value %q not in allowed options %v", prop.Name, item, optionVals))
		}
	}
	return errs
}

// ValidatePropertyValue validates a string input against a property type.
// This is used by the TUI for inline editing validation before accepting user input.
// For select/multi_select types, pass the property's Options; otherwise options can be nil.
// Returns nil if valid, or an error describing the validation failure.
func ValidatePropertyValue(propType string, options []Option, input string) error {
	switch propType {
	case "string":
		return nil
	case "number":
		return validateNumberInput(input)
	case "date":
		return validateDateInput(input)
	case "datetime":
		return validateDatetimeInput(input)
	case "url":
		return validateURLInput(input)
	case "select":
		return validateSelectInput(options, input)
	default:
		return nil
	}
}

// validateNumberInput checks if a string is a valid number (int or float).
func validateNumberInput(input string) error {
	if input == "" {
		return fmt.Errorf("expected a number, got empty string")
	}
	// Allow optional leading minus, digits, optional decimal point with digits
	valid := true
	dotSeen := false
	for i, c := range input {
		if c == '-' && i == 0 {
			continue
		}
		if c == '.' && !dotSeen {
			dotSeen = true
			continue
		}
		if c < '0' || c > '9' {
			valid = false
			break
		}
	}
	// Don't allow just "-" or "." or "-."
	if input == "-" || input == "." || input == "-." {
		valid = false
	}
	// Don't allow trailing dot like "123."
	if len(input) > 0 && input[len(input)-1] == '.' {
		valid = false
	}
	if !valid {
		return fmt.Errorf("expected a number, got %q", input)
	}
	return nil
}
