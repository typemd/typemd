package core

import (
	"regexp"
	"strings"
	"time"
)

var templatePlaceholderRegexp = regexp.MustCompile(`\{\{\s*(\w+):([^}]+)\}\}`)

var dateFormatReplacer = strings.NewReplacer(
	"YYYY", "2006",
	"MM", "01",
	"DD", "02",
	"HH", "15",
	"mm", "04",
	"ss", "05",
)

// EvaluateNameTemplate evaluates a name template string, replacing placeholders.
// Currently supports {{ date:FORMAT }} with user-friendly format tokens.
func EvaluateNameTemplate(template string, now time.Time) string {
	return templatePlaceholderRegexp.ReplaceAllStringFunc(template, func(match string) string {
		parts := templatePlaceholderRegexp.FindStringSubmatch(match)
		if len(parts) != 3 {
			return match
		}
		kind := strings.TrimSpace(parts[1])
		arg := strings.TrimSpace(parts[2])

		switch kind {
		case "date":
			goFormat := ConvertDateFormat(arg)
			return now.Format(goFormat)
		default:
			return match
		}
	})
}

// ConvertDateFormat converts user-friendly date format tokens to Go reference time.
// Tokens: YYYY (year), MM (month), DD (day), HH (24-hour), mm (minute), ss (second).
// Unrecognized tokens pass through as literal text.
func ConvertDateFormat(format string) string {
	return dateFormatReplacer.Replace(format)
}
