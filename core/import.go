package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ── Scan types ──────────────────────────────────────────────────────────────

// ScanResult holds the output of scanning source directories for markdown files.
type ScanResult struct {
	Sources           []SourceInfo       `json:"sources"`
	FileCount         int                `json:"file_count"`
	Directories       []DirInfo          `json:"directories"`
	Patterns          FrontmatterStats   `json:"patterns"`
	ExistingTypes     []ExistingTypeInfo `json:"existing_types"`
	NoFrontmatterCount int              `json:"no_frontmatter_count"`
}

// SourceInfo describes a single scanned markdown file.
type SourceInfo struct {
	Path           string            `json:"path"`
	Size           int64             `json:"size"`
	HasFrontmatter bool             `json:"has_frontmatter"`
	FrontmatterKeys map[string]any   `json:"frontmatter_keys,omitempty"`
}

// DirInfo describes a directory found during scanning.
type DirInfo struct {
	Path      string `json:"path"`
	FileCount int    `json:"file_count"`
}

// FrontmatterStats holds aggregate frontmatter key statistics.
type FrontmatterStats struct {
	Keys map[string]KeyStat `json:"keys"`
}

// KeyStat tracks how often a frontmatter key appears and sample values.
type KeyStat struct {
	Count   int      `json:"count"`
	Samples []string `json:"samples,omitempty"`
}

// ExistingTypeInfo describes a type schema already present in the vault.
type ExistingTypeInfo struct {
	Name       string   `json:"name"`
	Emoji      string   `json:"emoji,omitempty"`
	Properties []string `json:"properties,omitempty"`
}

// ── Plan types ──────────────────────────────────────────────────────────────

// ImportPlan describes a conversion plan: types to create and objects to import.
type ImportPlan struct {
	Types   []TypePlan   `json:"types"`
	Objects []ObjectPlan `json:"objects"`
	Order   []int        `json:"order"`
}

// TypePlan describes a type schema to create during import.
type TypePlan struct {
	Name       string     `json:"name"`
	Emoji      string     `json:"emoji,omitempty"`
	Plural     string     `json:"plural,omitempty"`
	Properties []Property `json:"properties,omitempty"`
}

// ObjectPlan describes how a source file maps to a vault object.
type ObjectPlan struct {
	SourcePath string         `json:"source_path"`
	TypeName   string         `json:"type_name"`
	Name       string         `json:"name"`
	Properties map[string]any `json:"properties,omitempty"`
	Body       string         `json:"body,omitempty"`
	Conflict   string         `json:"conflict"` // "none", "skip", "overwrite"
	DependsOn  []int          `json:"depends_on,omitempty"`
}

// ── Report types ────────────────────────────────────────────────────────────

// ImportReport summarizes the result of executing an import plan.
type ImportReport struct {
	TypesCreated    int              `json:"types_created"`
	ObjectsCreated  int              `json:"objects_created"`
	ObjectsSkipped  int              `json:"objects_skipped"`
	ObjectsFailed   int              `json:"objects_failed"`
	Details         []ImportDetail   `json:"details"`
	UnresolvedRefs  []UnresolvedRef  `json:"unresolved_refs,omitempty"`
	Suggestions     []string         `json:"suggestions,omitempty"`
}

// ImportDetail records the outcome of a single import operation.
type ImportDetail struct {
	SourcePath string `json:"source_path"`
	ObjectID   string `json:"object_id,omitempty"`
	Status     string `json:"status"` // "created", "skipped", "failed"
	Error      string `json:"error,omitempty"`
}

// UnresolvedRef records a reference that could not be resolved after import.
type UnresolvedRef struct {
	SourceObjectID string `json:"source_object_id"`
	Reference      string `json:"reference"`
}

// ── Scan implementation ─────────────────────────────────────────────────────

// ScanSources scans the given paths for markdown files, extracting frontmatter
// patterns and collecting file statistics. Paths are resolved relative to the
// vault root.
func (v *Vault) ScanSources(paths []string) (*ScanResult, error) {
	result := &ScanResult{
		Patterns: FrontmatterStats{Keys: make(map[string]KeyStat)},
	}

	dirCounts := make(map[string]int)

	for _, p := range paths {
		absPath := p
		if !filepath.IsAbs(p) {
			absPath = filepath.Join(v.Root, p)
		}

		info, err := os.Stat(absPath)
		if err != nil {
			return nil, fmt.Errorf("source path %q: %w", p, err)
		}

		if info.IsDir() {
			if err := v.scanDir(absPath, result, dirCounts); err != nil {
				return nil, err
			}
		} else {
			if strings.HasSuffix(strings.ToLower(absPath), ".md") {
				if err := v.scanFile(absPath, result); err != nil {
					return nil, err
				}
				dir := filepath.Dir(absPath)
				dirCounts[dir]++
			}
		}
	}

	// Build directory info
	for dir, count := range dirCounts {
		result.Directories = append(result.Directories, DirInfo{
			Path:      dir,
			FileCount: count,
		})
	}

	// Include existing vault types
	for _, typeName := range v.ListTypes() {
		schema, err := v.LoadType(typeName)
		if err != nil {
			continue
		}
		eti := ExistingTypeInfo{
			Name:  typeName,
			Emoji: schema.Emoji,
		}
		for _, prop := range schema.Properties {
			eti.Properties = append(eti.Properties, prop.Name)
		}
		result.ExistingTypes = append(result.ExistingTypes, eti)
	}

	return result, nil
}

func (v *Vault) scanDir(dir string, result *ScanResult, dirCounts map[string]int) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading directory %q: %w", dir, err)
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			if err := v.scanDir(fullPath, result, dirCounts); err != nil {
				return err
			}
			continue
		}
		if !strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
			continue
		}
		if err := v.scanFile(fullPath, result); err != nil {
			return err
		}
		dirCounts[dir]++
	}
	return nil
}

func (v *Vault) scanFile(path string, result *ScanResult) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	src := SourceInfo{
		Path: path,
		Size: info.Size(),
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	fm, _, err := parseFrontmatter(data)
	hasFM := err == nil && len(fm) > 0
	src.HasFrontmatter = hasFM
	if hasFM {
		src.FrontmatterKeys = fm
		for key, val := range fm {
			stat := result.Patterns.Keys[key]
			stat.Count++
			// Add sample values (up to 3)
			if len(stat.Samples) < 3 {
				s := fmt.Sprintf("%v", val)
				if s != "" {
					stat.Samples = append(stat.Samples, s)
				}
			}
			result.Patterns.Keys[key] = stat
		}
	} else {
		result.NoFrontmatterCount++
	}

	result.Sources = append(result.Sources, src)
	result.FileCount++
	return nil
}

// ── Plan implementation ─────────────────────────────────────────────────────

// GeneratePlan creates an import plan from scan results and AI-provided
// classifications. It determines import order by dependency and detects
// existing objects to set conflict status.
func (v *Vault) GeneratePlan(classifications []ObjectPlan) (*ImportPlan, error) {
	plan := &ImportPlan{}

	// Detect which types need to be created
	existingTypes := make(map[string]bool)
	for _, t := range v.ListTypes() {
		existingTypes[t] = true
	}

	neededTypes := make(map[string]bool)
	for _, obj := range classifications {
		if !existingTypes[obj.TypeName] && !neededTypes[obj.TypeName] {
			neededTypes[obj.TypeName] = true
			plan.Types = append(plan.Types, TypePlan{Name: obj.TypeName})
		}
	}

	// Detect conflicts with existing objects — batch queries by type
	existingNames := make(map[string]map[string]bool) // type → set of names
	if v.Queries != nil {
		for typeName := range existingTypes {
			results, err := v.index.Query(TypeFilter(typeName))
			if err != nil {
				continue
			}
			names := make(map[string]bool)
			for _, r := range results {
				if name, ok := r.Properties["name"].(string); ok {
					names[name] = true
				}
			}
			existingNames[typeName] = names
		}
	}

	for i := range classifications {
		if classifications[i].Conflict == "" {
			classifications[i].Conflict = "none"
		}
		if names, ok := existingNames[classifications[i].TypeName]; ok {
			if names[classifications[i].Name] {
				classifications[i].Conflict = "skip"
			}
		}
	}

	plan.Objects = classifications

	// Compute dependency order using topological sort
	plan.Order = v.computeImportOrder(classifications)

	return plan, nil
}

// computeImportOrder returns indices in dependency order.
// Tags come first, then objects without dependencies, then the rest.
func (v *Vault) computeImportOrder(objects []ObjectPlan) []int {
	n := len(objects)
	if n == 0 {
		return nil
	}

	// Build dependency graph
	inDegree := make([]int, n)
	for i, obj := range objects {
		inDegree[i] = len(obj.DependsOn)
	}

	// Separate tags from other objects for priority
	var tagIndices, otherReady, rest []int
	for i, obj := range objects {
		if inDegree[i] == 0 {
			if obj.TypeName == "tag" {
				tagIndices = append(tagIndices, i)
			} else {
				otherReady = append(otherReady, i)
			}
		} else {
			rest = append(rest, i)
		}
	}

	// Topological sort: tags first, then ready, then resolve dependencies
	var order []int
	order = append(order, tagIndices...)
	order = append(order, otherReady...)

	done := make(map[int]bool)
	for _, i := range order {
		done[i] = true
	}

	// Process remaining with dependencies
	changed := true
	for changed {
		changed = false
		for _, i := range rest {
			if done[i] {
				continue
			}
			allDone := true
			for _, dep := range objects[i].DependsOn {
				if !done[dep] {
					allDone = false
					break
				}
			}
			if allDone {
				order = append(order, i)
				done[i] = true
				changed = true
			}
		}
	}

	// Append any remaining (circular dependencies broken arbitrarily)
	for i := 0; i < n; i++ {
		if !done[i] {
			order = append(order, i)
		}
	}

	return order
}

// ── Execute implementation ──────────────────────────────────────────────────

// ExecutePlan executes an import plan: creates types, creates objects in
// dependency order, and runs reconciliation.
func (v *Vault) ExecutePlan(plan *ImportPlan) (*ImportReport, error) {
	report := &ImportReport{}

	// Phase 1: Create type schemas
	for _, tp := range plan.Types {
		// Skip if type already exists
		if _, err := v.LoadType(tp.Name); err == nil {
			continue
		}

		schema := &TypeSchema{
			Name:       tp.Name,
			Emoji:      tp.Emoji,
			Plural:     tp.Plural,
			Properties: tp.Properties,
		}
		if err := v.SaveType(schema); err != nil {
			return report, fmt.Errorf("creating type %q: %w", tp.Name, err)
		}
		report.TypesCreated++
	}

	// Phase 2: Create objects in dependency order
	for _, idx := range plan.Order {
		if idx < 0 || idx >= len(plan.Objects) {
			continue
		}
		obj := plan.Objects[idx]
		detail := ImportDetail{SourcePath: obj.SourcePath}

		switch obj.Conflict {
		case "skip":
			detail.Status = "skipped"
			report.ObjectsSkipped++
			report.Details = append(report.Details, detail)
			continue
		}

		created, err := v.NewObject(obj.TypeName, obj.Name, "")
		if err != nil {
			detail.Status = "failed"
			detail.Error = err.Error()
			report.ObjectsFailed++
			report.Details = append(report.Details, detail)
			continue
		}

		// Set properties — collect errors but don't fail the whole import
		var propErrors []string
		for key, val := range obj.Properties {
			if err := v.SetProperty(created.ID, key, val); err != nil {
				propErrors = append(propErrors, fmt.Sprintf("%s: %v", key, err))
			}
		}

		// Set body content
		if obj.Body != "" {
			created.Body = obj.Body
			if err := v.SaveObject(created); err != nil {
				propErrors = append(propErrors, fmt.Sprintf("body: %v", err))
			}
		}

		detail.ObjectID = created.ID
		detail.Status = "created"
		if len(propErrors) > 0 {
			detail.Error = strings.Join(propErrors, "; ")
		}
		report.ObjectsCreated++
		report.Details = append(report.Details, detail)
	}

	// Phase 3: Reconcile to resolve wiki-links
	if v.reconciler != nil {
		_, _, _ = v.reconciler.Reconcile()
	}

	// Phase 4: Check for unresolved references
	v.detectUnresolvedRefs(report)

	// Phase 5: Generate suggestions
	v.generateSuggestions(report)

	return report, nil
}

// detectUnresolvedRefs checks imported objects for unresolved wiki-links
// by looking at stored wiki-link records in the index.
func (v *Vault) detectUnresolvedRefs(report *ImportReport) {
	if v.index == nil {
		return
	}
	for _, detail := range report.Details {
		if detail.Status != "created" || detail.ObjectID == "" {
			continue
		}
		storedLinks, err := v.index.ListWikiLinks(detail.ObjectID)
		if err != nil {
			continue
		}
		for _, link := range storedLinks {
			if link.ToID == "" {
				report.UnresolvedRefs = append(report.UnresolvedRefs, UnresolvedRef{
					SourceObjectID: detail.ObjectID,
					Reference:      link.Target,
				})
			}
		}
	}
}

// generateSuggestions adds follow-up suggestions based on import results.
func (v *Vault) generateSuggestions(report *ImportReport) {
	if len(report.UnresolvedRefs) > 0 {
		report.Suggestions = append(report.Suggestions,
			fmt.Sprintf("Create objects for %d unresolved references", len(report.UnresolvedRefs)))
	}
	if report.ObjectsFailed > 0 {
		report.Suggestions = append(report.Suggestions,
			fmt.Sprintf("Review and re-import %d failed files", report.ObjectsFailed))
	}
}
