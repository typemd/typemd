package core

import (
	"fmt"
	"strings"
)

// ObjectID is a Value Object representing a unique object identifier.
// Format: "type/filename" where filename includes a ULID suffix.
type ObjectID struct {
	Type     string
	Filename string
}

// ParseObjectID parses a raw "type/filename" string into an ObjectID.
func ParseObjectID(raw string) (ObjectID, error) {
	parts := strings.SplitN(raw, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return ObjectID{}, fmt.Errorf("invalid object ID: %q", raw)
	}
	return ObjectID{Type: parts[0], Filename: parts[1]}, nil
}

// MustParseObjectID parses a raw string into an ObjectID, panicking on invalid input.
// Use only in tests or with known-valid IDs.
func MustParseObjectID(raw string) ObjectID {
	id, err := ParseObjectID(raw)
	if err != nil {
		panic(err)
	}
	return id
}

// NewObjectID creates an ObjectID with a generated ULID suffix.
func NewObjectID(typeName, slug string) (ObjectID, error) {
	ulidStr, err := GenerateULID()
	if err != nil {
		return ObjectID{}, err
	}
	return ObjectID{
		Type:     typeName,
		Filename: slug + "-" + ulidStr,
	}, nil
}

// String returns the canonical "type/filename" representation.
func (id ObjectID) String() string {
	return id.Type + "/" + id.Filename
}

// DisplayName returns the filename with ULID suffix stripped.
func (id ObjectID) DisplayName() string {
	return StripULID(id.Filename)
}

// DisplayID returns "type/display-name" with ULID suffix stripped.
func (id ObjectID) DisplayID() string {
	return id.Type + "/" + id.DisplayName()
}

// IsZero returns true if the ObjectID is uninitialized.
func (id ObjectID) IsZero() bool {
	return id.Type == "" && id.Filename == ""
}
