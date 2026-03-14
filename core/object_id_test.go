package core

import "testing"

func TestParseObjectID(t *testing.T) {
	id, err := ParseObjectID("book/clean-code-01abc")
	if err != nil {
		t.Fatalf("ParseObjectID: %v", err)
	}
	if id.Type != "book" {
		t.Errorf("Type = %q, want %q", id.Type, "book")
	}
	if id.Filename != "clean-code-01abc" {
		t.Errorf("Filename = %q, want %q", id.Filename, "clean-code-01abc")
	}
}

func TestParseObjectID_Invalid(t *testing.T) {
	cases := []string{"", "noslash", "/empty-type", "empty-filename/"}
	for _, raw := range cases {
		_, err := ParseObjectID(raw)
		if err == nil {
			t.Errorf("ParseObjectID(%q) should fail", raw)
		}
	}
}

func TestObjectID_String(t *testing.T) {
	id := ObjectID{Type: "book", Filename: "clean-code-01abc"}
	if s := id.String(); s != "book/clean-code-01abc" {
		t.Errorf("String() = %q, want %q", s, "book/clean-code-01abc")
	}
}

func TestObjectID_DisplayName(t *testing.T) {
	id := ObjectID{Type: "book", Filename: "clean-code-01jqr3k5mpbvn8e0f2g7h9txyz"}
	if dn := id.DisplayName(); dn != "clean-code" {
		t.Errorf("DisplayName() = %q, want %q", dn, "clean-code")
	}
}

func TestObjectID_DisplayID(t *testing.T) {
	id := ObjectID{Type: "book", Filename: "clean-code-01jqr3k5mpbvn8e0f2g7h9txyz"}
	if di := id.DisplayID(); di != "book/clean-code" {
		t.Errorf("DisplayID() = %q, want %q", di, "book/clean-code")
	}
}

func TestObjectID_IsZero(t *testing.T) {
	var zero ObjectID
	if !zero.IsZero() {
		t.Error("zero ObjectID should be zero")
	}
	id := ObjectID{Type: "book", Filename: "test"}
	if id.IsZero() {
		t.Error("non-zero ObjectID should not be zero")
	}
}

func TestNewObjectID(t *testing.T) {
	id, err := NewObjectID("book", "clean-code")
	if err != nil {
		t.Fatalf("NewObjectID: %v", err)
	}
	if id.Type != "book" {
		t.Errorf("Type = %q, want %q", id.Type, "book")
	}
	if id.DisplayName() != "clean-code" {
		t.Errorf("DisplayName() = %q, want %q", id.DisplayName(), "clean-code")
	}
}
