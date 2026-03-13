package core

import "testing"

func TestPluralName(t *testing.T) {
	tests := []struct {
		name   string
		schema TypeSchema
		want   string
	}{
		{
			name:   "plural set",
			schema: TypeSchema{Name: "book", Plural: "books"},
			want:   "books",
		},
		{
			name:   "plural empty falls back to name",
			schema: TypeSchema{Name: "book"},
			want:   "book",
		},
		{
			name:   "non-English name without plural",
			schema: TypeSchema{Name: "筆記"},
			want:   "筆記",
		},
		{
			name:   "non-English with plural",
			schema: TypeSchema{Name: "person", Plural: "people"},
			want:   "people",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.schema.PluralName()
			if got != tt.want {
				t.Errorf("PluralName() = %q, want %q", got, tt.want)
			}
		})
	}
}
