package core

import (
	"testing"
)

func TestSortObjects_EmptySlice(t *testing.T) {
	// Should not panic
	SortObjects(nil, []SortRule{{Property: "name", Direction: "asc"}})
	SortObjects([]*Object{}, []SortRule{{Property: "name", Direction: "asc"}})
}

func TestSortObjects_SingleElement(t *testing.T) {
	objs := []*Object{
		{Type: "book", Properties: map[string]any{"name": "Alpha"}},
	}
	SortObjects(objs, []SortRule{{Property: "name", Direction: "asc"}})
	if objs[0].GetName() != "Alpha" {
		t.Errorf("expected Alpha, got %s", objs[0].GetName())
	}
}

func TestSortObjects_NoRules(t *testing.T) {
	objs := []*Object{
		{Type: "book", Properties: map[string]any{"name": "Bravo"}},
		{Type: "book", Properties: map[string]any{"name": "Alpha"}},
	}
	SortObjects(objs, nil)
	// Order should be unchanged
	if objs[0].GetName() != "Bravo" {
		t.Errorf("expected Bravo first (unchanged order), got %s", objs[0].GetName())
	}
}

func TestSortObjects_AllNilValues(t *testing.T) {
	objs := []*Object{
		{Type: "book", Properties: map[string]any{"name": "A"}},
		{Type: "book", Properties: map[string]any{"name": "B"}},
		{Type: "book", Properties: map[string]any{"name": "C"}},
	}
	// Sort by "rating" which none of them have
	SortObjects(objs, []SortRule{{Property: "rating", Direction: "asc"}})
	// All nil -> stable sort preserves order
	if objs[0].GetName() != "A" || objs[1].GetName() != "B" || objs[2].GetName() != "C" {
		t.Errorf("expected A,B,C (stable), got %s,%s,%s",
			objs[0].GetName(), objs[1].GetName(), objs[2].GetName())
	}
}

func TestSortObjects_MixedNumericAndNonNumeric(t *testing.T) {
	objs := []*Object{
		{Type: "book", Properties: map[string]any{"name": "NumStr", "score": "not-a-number"}},
		{Type: "book", Properties: map[string]any{"name": "Num5", "score": 5.0}},
		{Type: "book", Properties: map[string]any{"name": "Num2", "score": 2.0}},
	}
	SortObjects(objs, []SortRule{{Property: "score", Direction: "asc"}})
	// Numeric values (2, 5) sort first; "not-a-number" falls back to string comparison.
	// toFloat64 will fail for "not-a-number", so compareSortValues does string comparison
	// between numeric and string. When one side is numeric and the other is string,
	// both fall back to string comparison via Sprintf.
	// 2 -> "2", 5 -> "5", "not-a-number" -> "not-a-number"
	// "2" < "5" < "n..." so order should be Num2, Num5, NumStr
	if objs[0].GetName() != "Num2" {
		t.Errorf("expected Num2 first, got %s", objs[0].GetName())
	}
	if objs[1].GetName() != "Num5" {
		t.Errorf("expected Num5 second, got %s", objs[1].GetName())
	}
	if objs[2].GetName() != "NumStr" {
		t.Errorf("expected NumStr last, got %s", objs[2].GetName())
	}
}

func TestSortObjects_MultipleSortRules(t *testing.T) {
	objs := []*Object{
		{Type: "book", Properties: map[string]any{"name": "B-High", "status": "active", "rating": 9.0}},
		{Type: "book", Properties: map[string]any{"name": "A-Low", "status": "active", "rating": 3.0}},
		{Type: "book", Properties: map[string]any{"name": "C-Draft", "status": "draft", "rating": 5.0}},
	}
	// Primary: status asc, Secondary: rating asc
	SortObjects(objs, []SortRule{
		{Property: "status", Direction: "asc"},
		{Property: "rating", Direction: "asc"},
	})
	// "active" (A-Low r=3, B-High r=9), then "draft" (C-Draft r=5)
	expected := []string{"A-Low", "B-High", "C-Draft"}
	for i, name := range expected {
		if objs[i].GetName() != name {
			t.Errorf("position %d: expected %s, got %s", i, name, objs[i].GetName())
		}
	}
}

func TestSortObjects_MultipleSortRules_SecondaryBreaksTie(t *testing.T) {
	objs := []*Object{
		{Type: "book", Properties: map[string]any{"name": "Z", "status": "active", "rating": 9.0}},
		{Type: "book", Properties: map[string]any{"name": "A", "status": "active", "rating": 1.0}},
		{Type: "book", Properties: map[string]any{"name": "M", "status": "active", "rating": 5.0}},
	}
	// Primary: status asc (all same), Secondary: rating desc
	SortObjects(objs, []SortRule{
		{Property: "status", Direction: "asc"},
		{Property: "rating", Direction: "desc"},
	})
	// All "active", so sort by rating desc: 9, 5, 1
	expected := []string{"Z", "M", "A"}
	for i, name := range expected {
		if objs[i].GetName() != name {
			t.Errorf("position %d: expected %s, got %s", i, name, objs[i].GetName())
		}
	}
}

func TestSortObjects_NilObjectInSlice(t *testing.T) {
	objs := []*Object{
		{Type: "book", Properties: map[string]any{"name": "Alpha"}},
		nil,
		{Type: "book", Properties: map[string]any{"name": "Bravo"}},
	}
	// Should not panic — nil objects get nil sort values, sorted last
	SortObjects(objs, []SortRule{{Property: "name", Direction: "asc"}})

	// Non-nil objects should come first
	if objs[0] == nil || objs[0].GetName() != "Alpha" {
		t.Errorf("expected Alpha first, got %v", objs[0])
	}
	if objs[1] == nil || objs[1].GetName() != "Bravo" {
		t.Errorf("expected Bravo second, got %v", objs[1])
	}
	if objs[2] != nil {
		t.Errorf("expected nil last, got %v", objs[2])
	}
}

func TestSortObjects_DescendingOrder(t *testing.T) {
	objs := []*Object{
		{Type: "book", Properties: map[string]any{"name": "A", "rating": 1.0}},
		{Type: "book", Properties: map[string]any{"name": "B", "rating": 9.0}},
		{Type: "book", Properties: map[string]any{"name": "C", "rating": 5.0}},
	}
	SortObjects(objs, []SortRule{{Property: "rating", Direction: "desc"}})
	expected := []string{"B", "C", "A"}
	for i, name := range expected {
		if objs[i].GetName() != name {
			t.Errorf("position %d: expected %s, got %s", i, name, objs[i].GetName())
		}
	}
}

func TestSortObjects_NilWithDesc(t *testing.T) {
	objs := []*Object{
		{Type: "book", Properties: map[string]any{"name": "NoRating"}},
		{Type: "book", Properties: map[string]any{"name": "Rated", "rating": 5.0}},
	}
	SortObjects(objs, []SortRule{{Property: "rating", Direction: "desc"}})
	// Nil values sort last even in desc mode
	if objs[0].GetName() != "Rated" {
		t.Errorf("expected Rated first, got %s", objs[0].GetName())
	}
	if objs[1].GetName() != "NoRating" {
		t.Errorf("expected NoRating last, got %s", objs[1].GetName())
	}
}
