package core

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

// GenerateULID returns a new lowercase ULID string.
func GenerateULID() (string, error) {
	id, err := ulid.New(ulid.Timestamp(time.Now()), rand.Reader)
	if err != nil {
		return "", fmt.Errorf("generate ulid: %w", err)
	}
	return strings.ToLower(id.String()), nil
}

var ulidSuffixPattern = regexp.MustCompile(`-[0-9a-hjkmnp-tv-z]{26}$`)

// StripULID removes the ULID suffix from a filename if present.
// ULID suffix pattern: hyphen followed by exactly 26 Crockford's Base32 characters.
func StripULID(filename string) string {
	return ulidSuffixPattern.ReplaceAllString(filename, "")
}
