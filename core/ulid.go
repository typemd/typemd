package core

import (
	"crypto/rand"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

// GenerateULID returns a new lowercase ULID string.
func GenerateULID() string {
	return strings.ToLower(ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String())
}
