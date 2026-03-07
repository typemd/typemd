package core

import (
	"crypto/rand"
	"fmt"
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
