package widget

import "testing"

func TestAdjustScroll(t *testing.T) {
	tests := []struct {
		name       string
		cursor     int
		offset     int
		viewHeight int
		want       int
	}{
		{"cursor above viewport", 2, 5, 10, 2},
		{"cursor below viewport", 15, 5, 10, 6},
		{"cursor within viewport", 7, 5, 10, 5},
		{"zero view height", 5, 0, 0, 0},
		{"cursor at top edge", 5, 5, 10, 5},
		{"cursor at bottom edge", 14, 5, 10, 5},
		{"cursor one past bottom", 15, 5, 10, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AdjustScroll(tt.cursor, tt.offset, tt.viewHeight)
			if got != tt.want {
				t.Errorf("AdjustScroll(%d, %d, %d) = %d, want %d",
					tt.cursor, tt.offset, tt.viewHeight, got, tt.want)
			}
		})
	}
}
