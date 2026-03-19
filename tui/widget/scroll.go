package widget

// AdjustScroll returns a scroll offset that keeps cursor visible within viewHeight.
func AdjustScroll(cursor, offset, viewHeight int) int {
	if viewHeight <= 0 {
		return 0
	}
	if cursor < offset {
		return cursor
	}
	if cursor >= offset+viewHeight {
		return cursor - viewHeight + 1
	}
	return offset
}
