package core

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// marshalSchema is the YAML-serializable form of TypeSchema.
// It re-introduces the name property entry when NameTemplate is set,
// since NameTemplate is excluded from standard yaml marshaling (yaml:"-").
type marshalSchema struct {
	Name        string     `yaml:"name"`
	Plural      string     `yaml:"plural,omitempty"`
	Emoji       string     `yaml:"emoji,omitempty"`
	Color       string     `yaml:"color,omitempty"`
	Unique      bool       `yaml:"unique,omitempty"`
	Version     string     `yaml:"version,omitempty"`
	Description string     `yaml:"description,omitempty"`
	Properties  []Property `yaml:"properties"`
}

// MarshalTypeSchema serializes a TypeSchema to YAML bytes suitable for
// writing to .typemd/types/<name>.yaml. It handles the NameTemplate →
// name property entry conversion that yaml:"-" would otherwise drop.
func MarshalTypeSchema(schema *TypeSchema) ([]byte, error) {
	version := schema.Version
	if version == DefaultSchemaVersion {
		version = "" // omit default version from YAML
	}
	ms := marshalSchema{
		Name:        schema.Name,
		Plural:      schema.Plural,
		Emoji:       schema.Emoji,
		Color:       schema.Color,
		Unique:      schema.Unique,
		Version:     version,
		Description: schema.Description,
	}

	// Re-introduce name template as a name property entry if set
	if schema.NameTemplate != "" {
		ms.Properties = append(ms.Properties, Property{
			Name:     NameProperty,
			Template: schema.NameTemplate,
		})
	}

	ms.Properties = append(ms.Properties, schema.Properties...)

	data, err := yaml.Marshal(&ms)
	if err != nil {
		return nil, err
	}
	// yaml.v3 escapes non-ASCII characters (e.g. emoji 👤 → "\U0001F464").
	// Restore UTF-8 literals for human readability.
	return unescapeUnicodeYAML(data), nil
}

// unicodeEscapePattern matches yaml.v3 unicode escape sequences:
// \uXXXX (4 hex digits) and \UXXXXXXXX (8 hex digits).
var unicodeEscapePattern = regexp.MustCompile(`\\[uU][0-9a-fA-F]{4,8}`)

// unescapeUnicodeYAML replaces yaml.v3 unicode escape sequences with
// their UTF-8 literals. yaml.v3 escapes non-ASCII characters like emoji
// (e.g. 👤 → "\U0001F464"), which hurts readability.
func unescapeUnicodeYAML(data []byte) []byte {
	return unicodeEscapePattern.ReplaceAllFunc(data, func(match []byte) []byte {
		// Parse the hex value (skip the \u or \U prefix)
		hexStr := string(match[2:])
		r, err := strconv.ParseInt(hexStr, 16, 32)
		if err != nil {
			return match
		}
		return []byte(string(rune(r)))
	})
}

// parseVersion parses a "major.minor" version string into two integers.
func parseVersion(v string) (int, int, error) {
	parts := strings.SplitN(v, ".", 3)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("version must be in \"major.minor\" format, got %q", v)
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil || major < 0 {
		return 0, 0, fmt.Errorf("version must be in \"major.minor\" format, got %q", v)
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil || minor < 0 {
		return 0, 0, fmt.Errorf("version must be in \"major.minor\" format, got %q", v)
	}
	// Reject leading zeros (e.g. "01.0" or "0.01")
	if parts[0] != strconv.Itoa(major) || parts[1] != strconv.Itoa(minor) {
		return 0, 0, fmt.Errorf("version must be in \"major.minor\" format, got %q", v)
	}
	return major, minor, nil
}

// validateVersion checks that a version string is a valid "major.minor" format.
func validateVersion(v string) error {
	_, _, err := parseVersion(v)
	return err
}

// validColorPresets contains the list of allowed preset color names.
var validColorPresets = []string{
	"red", "blue", "green", "yellow", "purple",
	"orange", "pink", "cyan", "gray", "brown",
}

// validColorPresetSet is a lookup set for validColorPresets.
var validColorPresetSet = func() map[string]bool {
	m := make(map[string]bool, len(validColorPresets))
	for _, c := range validColorPresets {
		m[c] = true
	}
	return m
}()

// hexColorRegexp matches #RGB or #RRGGBB (case-insensitive hex digits).
var hexColorRegexp = regexp.MustCompile(`^#(?:[0-9a-fA-F]{3}|[0-9a-fA-F]{6})$`)

// ValidColorPresets returns the list of valid preset color names.
func ValidColorPresets() []string {
	result := make([]string, len(validColorPresets))
	copy(result, validColorPresets)
	return result
}

// validateColor checks that a color value is either a valid preset name or a valid hex code.
func validateColor(color string) error {
	if validColorPresetSet[color] {
		return nil
	}
	if hexColorRegexp.MatchString(color) {
		return nil
	}
	return fmt.Errorf("invalid color %q (valid: preset name or #RGB/#RRGGBB hex)", color)
}

// CompareVersions compares two "major.minor" version strings.
// Returns -1 if a < b, 0 if a == b, 1 if a > b.
// Both versions must be valid; invalid versions are treated as "0.0".
func CompareVersions(a, b string) int {
	aMajor, aMinor, _ := parseVersion(a)
	bMajor, bMinor, _ := parseVersion(b)
	if aMajor != bMajor {
		if aMajor < bMajor {
			return -1
		}
		return 1
	}
	if aMinor != bMinor {
		if aMinor < bMinor {
			return -1
		}
		return 1
	}
	return 0
}
