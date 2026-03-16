package core

import (
	"regexp"
	"strings"
	"unicode"
)

var multiHyphenPattern = regexp.MustCompile(`-{2,}`)

// Slugify converts a natural-language name to a valid slug.
// It lowercases, replaces spaces/underscores with hyphens, removes non-alphanumeric
// characters (except hyphens), collapses consecutive hyphens, and trims.
func Slugify(name string) string {
	// Lowercase
	s := strings.ToLower(name)

	// Replace spaces and underscores with hyphens
	s = strings.Map(func(r rune) rune {
		switch {
		case r == ' ' || r == '_':
			return '-'
		case unicode.IsLetter(r):
			return r
		case unicode.IsDigit(r):
			return r
		case r == '-':
			return r
		default:
			return -1 // remove
		}
	}, s)

	// Collapse consecutive hyphens
	s = multiHyphenPattern.ReplaceAllString(s, "-")

	// Trim leading/trailing hyphens
	s = strings.Trim(s, "-")

	return s
}
