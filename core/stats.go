package core

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"
)

// VaultStats holds aggregate statistics across all types.
type VaultStats struct {
	Types []TypeSummary `json:"types"`
	Total int           `json:"total"`
}

// TypeSummary holds per-type summary information.
type TypeSummary struct {
	Name        string    `json:"name"`
	Plural      string    `json:"plural,omitempty"`
	Emoji       string    `json:"emoji,omitempty"`
	Count       int       `json:"count"`
	LastUpdated time.Time `json:"last_updated,omitempty"`
}

// DisplayName returns the plural name if set, otherwise the type name.
func (ts TypeSummary) DisplayName() string {
	if ts.Plural != "" {
		return ts.Plural
	}
	return ts.Name
}

// TypeStats holds detailed property statistics for a single type.
type TypeStats struct {
	TypeName   string          `json:"type_name"`
	Emoji      string          `json:"emoji,omitempty"`
	Plural     string          `json:"plural,omitempty"`
	Count      int             `json:"count"`
	Properties []PropertyStats `json:"properties"`
}

// PropertyStats holds aggregated stats for a single property.
type PropertyStats struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Filled int    `json:"filled"`
	Total  int    `json:"total"`
	Stats  any    `json:"stats,omitempty"`
}

// NumberStats holds aggregation for number properties.
type NumberStats struct {
	Sum float64 `json:"sum"`
	Avg float64 `json:"avg"`
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// SelectStats holds frequency distribution for select/multi_select properties.
type SelectStats struct {
	Distribution map[string]int `json:"distribution"`
}

// CheckboxStats holds true/false counts for checkbox properties.
type CheckboxStats struct {
	TrueCount  int `json:"true_count"`
	FalseCount int `json:"false_count"`
}

// DateStats holds earliest/latest for date/datetime properties.
type DateStats struct {
	Earliest string `json:"earliest"`
	Latest   string `json:"latest"`
}

// RelationStats holds link counts for relation properties.
type RelationStats struct {
	Count int `json:"count"`
}

// VaultStats returns aggregate statistics for all types in the vault.
func (s *QueryService) VaultStats() (*VaultStats, error) {
	typeNames, err := s.repo.ListSchemas()
	if err != nil {
		return nil, fmt.Errorf("list types: %w", err)
	}

	stats := &VaultStats{}

	for _, name := range typeNames {
		results, err := s.index.Query("type=" + name)
		if err != nil {
			return nil, fmt.Errorf("query type %s: %w", name, err)
		}
		if len(results) == 0 {
			continue
		}

		schema, _ := s.repo.GetSchema(name)
		summary := TypeSummary{
			Name:  name,
			Count: len(results),
		}
		if schema != nil {
			summary.Emoji = schema.Emoji
			summary.Plural = schema.Plural
		}

		// Find latest updated_at
		for _, r := range results {
			if ua, ok := r.Properties[UpdatedAtProperty]; ok {
				if uaStr, ok := ua.(string); ok {
					if t, err := time.Parse(time.RFC3339, uaStr); err == nil {
						if t.After(summary.LastUpdated) {
							summary.LastUpdated = t
						}
					}
				}
			}
		}

		stats.Types = append(stats.Types, summary)
		stats.Total += len(results)
	}

	sort.Slice(stats.Types, func(i, j int) bool {
		return stats.Types[i].Name < stats.Types[j].Name
	})

	return stats, nil
}

// TypeStats returns per-property aggregate statistics for a single type.
func (s *QueryService) TypeStats(typeName string) (*TypeStats, error) {
	schema, err := s.repo.GetSchema(typeName)
	if err != nil {
		return nil, fmt.Errorf("load type %q: %w", typeName, err)
	}

	results, err := s.index.Query("type=" + typeName)
	if err != nil {
		return nil, fmt.Errorf("query type %s: %w", typeName, err)
	}

	objects := objectResultsToObjects(results)
	ts := &TypeStats{
		TypeName: typeName,
		Emoji:    schema.Emoji,
		Plural:   schema.Plural,
		Count:    len(objects),
	}

	for _, prop := range schema.Properties {
		ps := computePropertyStats(prop, objects)
		ts.Properties = append(ts.Properties, ps)
	}

	return ts, nil
}

// computePropertyStats aggregates stats for a single property across objects.
func computePropertyStats(prop Property, objects []*Object) PropertyStats {
	ps := PropertyStats{
		Name:  prop.Name,
		Type:  prop.Type,
		Total: len(objects),
	}

	switch prop.Type {
	case "number":
		if s := computeNumberStats(prop.Name, objects, &ps.Filled); s != nil {
			ps.Stats = s
		}
	case "select":
		if s := computeSelectStats(prop.Name, objects, &ps.Filled); s != nil {
			ps.Stats = s
		}
	case "multi_select":
		if s := computeMultiSelectStats(prop.Name, objects, &ps.Filled); s != nil {
			ps.Stats = s
		}
	case "checkbox":
		if s := computeCheckboxStats(prop.Name, objects, &ps.Filled); s != nil {
			ps.Stats = s
		}
	case "date", "datetime":
		if s := computeDateStats(prop.Name, objects, &ps.Filled); s != nil {
			ps.Stats = s
		}
	case "relation":
		if s := computeRelationStats(prop.Name, objects, &ps.Filled); s != nil {
			ps.Stats = s
		}
	default:
		// string, url — count filled only
		for _, obj := range objects {
			if v, ok := obj.Properties[prop.Name]; ok && v != nil && v != "" {
				ps.Filled++
			}
		}
	}

	return ps
}

func computeNumberStats(propName string, objects []*Object, filled *int) *NumberStats {
	var values []float64
	for _, obj := range objects {
		v, ok := obj.Properties[propName]
		if !ok || v == nil {
			continue
		}
		f, err := toFloat64(v)
		if err != nil {
			continue
		}
		values = append(values, f)
	}
	*filled = len(values)
	if len(values) == 0 {
		return nil
	}

	stats := &NumberStats{Min: values[0], Max: values[0]}
	for _, v := range values {
		stats.Sum += v
		if v < stats.Min {
			stats.Min = v
		}
		if v > stats.Max {
			stats.Max = v
		}
	}
	stats.Avg = stats.Sum / float64(len(values))
	return stats
}

func computeSelectStats(propName string, objects []*Object, filled *int) *SelectStats {
	dist := make(map[string]int)
	for _, obj := range objects {
		v, ok := obj.Properties[propName]
		if !ok || v == nil {
			continue
		}
		s := fmt.Sprintf("%v", v)
		if s == "" {
			continue
		}
		dist[s]++
		*filled++
	}
	if len(dist) == 0 {
		return nil
	}
	return &SelectStats{Distribution: dist}
}

func computeMultiSelectStats(propName string, objects []*Object, filled *int) *SelectStats {
	dist := make(map[string]int)
	for _, obj := range objects {
		v, ok := obj.Properties[propName]
		if !ok || v == nil {
			continue
		}
		switch val := v.(type) {
		case []any:
			if len(val) == 0 {
				continue
			}
			*filled++
			for _, item := range val {
				dist[fmt.Sprintf("%v", item)]++
			}
		case []string:
			if len(val) == 0 {
				continue
			}
			*filled++
			for _, item := range val {
				dist[item]++
			}
		default:
			s := fmt.Sprintf("%v", v)
			if s != "" {
				dist[s]++
				*filled++
			}
		}
	}
	if len(dist) == 0 {
		return nil
	}
	return &SelectStats{Distribution: dist}
}

func computeCheckboxStats(propName string, objects []*Object, filled *int) *CheckboxStats {
	stats := &CheckboxStats{}
	for _, obj := range objects {
		v, ok := obj.Properties[propName]
		if !ok || v == nil {
			continue
		}
		b, err := toBool(v)
		if err != nil {
			continue
		}
		*filled++
		if b {
			stats.TrueCount++
		} else {
			stats.FalseCount++
		}
	}
	if *filled == 0 {
		return nil
	}
	return stats
}

func computeDateStats(propName string, objects []*Object, filled *int) *DateStats {
	var earliest, latest string
	for _, obj := range objects {
		v, ok := obj.Properties[propName]
		if !ok || v == nil {
			continue
		}
		s := fmt.Sprintf("%v", v)
		if s == "" {
			continue
		}
		*filled++
		if earliest == "" || s < earliest {
			earliest = s
		}
		if latest == "" || s > latest {
			latest = s
		}
	}
	if *filled == 0 {
		return nil
	}
	return &DateStats{Earliest: earliest, Latest: latest}
}

func computeRelationStats(propName string, objects []*Object, filled *int) *RelationStats {
	count := 0
	for _, obj := range objects {
		v, ok := obj.Properties[propName]
		if !ok || v == nil {
			continue
		}
		switch val := v.(type) {
		case []any:
			if len(val) > 0 {
				*filled++
				count += len(val)
			}
		case string:
			if val != "" {
				*filled++
				count++
			}
		default:
			s := fmt.Sprintf("%v", v)
			if s != "" {
				*filled++
				count++
			}
		}
	}
	if count == 0 {
		return nil
	}
	return &RelationStats{Count: count}
}

// toFloat64 converts various numeric representations to float64.
func toFloat64(v any) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	case json.Number:
		return val.Float64()
	default:
		return strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
	}
}

// toBool converts various boolean representations to bool.
func toBool(v any) (bool, error) {
	switch val := v.(type) {
	case bool:
		return val, nil
	case string:
		return strconv.ParseBool(val)
	default:
		return strconv.ParseBool(fmt.Sprintf("%v", v))
	}
}
