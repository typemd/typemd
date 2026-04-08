package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
)

func (dc *domainContext) aSourceDirectoryWithMarkdownFiles(dir string, table *godog.Table) {
	srcDir := filepath.Join(dc.vault.Root, dir)
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		panic(fmt.Sprintf("mkdir %s: %v", dir, err))
	}

	headers := table.Rows[0].Cells
	for _, row := range table.Rows[1:] {
		vals := make(map[string]string)
		for i, cell := range row.Cells {
			vals[headers[i].Value] = cell.Value
		}

		filename := vals["filename"]
		var content strings.Builder
		// Build frontmatter if any non-filename fields have values
		hasFM := false
		for k, v := range vals {
			if k != "filename" && k != "body" && v != "" {
				hasFM = true
				break
			}
		}
		if hasFM {
			content.WriteString("---\n")
			for k, v := range vals {
				if k != "filename" && k != "body" && v != "" {
					content.WriteString(fmt.Sprintf("%s: %s\n", k, v))
				}
			}
			content.WriteString("---\n")
		}
		if body, ok := vals["body"]; ok && body != "" {
			content.WriteString("\n" + body + "\n")
		}

		if err := os.WriteFile(filepath.Join(srcDir, filename), []byte(content.String()), 0644); err != nil {
			panic(fmt.Sprintf("write %s/%s: %v", dir, filename, err))
		}
	}
}

func (dc *domainContext) aSourceDirectoryWithNoMarkdownFiles(dir string) {
	srcDir := filepath.Join(dc.vault.Root, dir)
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		panic(fmt.Sprintf("mkdir %s: %v", dir, err))
	}
	// Write a non-markdown file
	if err := os.WriteFile(filepath.Join(srcDir, "image.png"), []byte("fake png"), 0644); err != nil {
		panic(fmt.Sprintf("write %s/image.png: %v", dir, err))
	}
}

func (dc *domainContext) aSourceDirectoryWithPlainMarkdownFiles(dir string, table *godog.Table) {
	srcDir := filepath.Join(dc.vault.Root, dir)
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		panic(fmt.Sprintf("mkdir %s: %v", dir, err))
	}

	headers := table.Rows[0].Cells
	for _, row := range table.Rows[1:] {
		vals := make(map[string]string)
		for i, cell := range row.Cells {
			vals[headers[i].Value] = cell.Value
		}

		filename := vals["filename"]
		body := vals["body"]
		// No frontmatter — just plain markdown
		if err := os.WriteFile(filepath.Join(srcDir, filename), []byte(body+"\n"), 0644); err != nil {
			panic(fmt.Sprintf("write %s/%s: %v", dir, filename, err))
		}
	}
}

func (dc *domainContext) iScanSources(paths string) {
	pathList := strings.Split(paths, ",")
	dc.scanResult, dc.lastErr = dc.vault.ScanSources(pathList)
}

func (dc *domainContext) theScanResultShouldHaveNFiles(count int) error {
	if dc.scanResult == nil {
		return fmt.Errorf("scan result is nil")
	}
	if dc.scanResult.FileCount != count {
		return fmt.Errorf("expected %d files, got %d", count, dc.scanResult.FileCount)
	}
	return nil
}

func (dc *domainContext) theScanFrontmatterShouldShowKeyAppearingNTimes(key string, count int) error {
	if dc.scanResult == nil {
		return fmt.Errorf("scan result is nil")
	}
	stat, ok := dc.scanResult.Patterns.Keys[key]
	if !ok {
		return fmt.Errorf("key %q not found in frontmatter patterns", key)
	}
	if stat.Count != count {
		return fmt.Errorf("key %q: expected count %d, got %d", key, count, stat.Count)
	}
	return nil
}

func (dc *domainContext) theScanResultShouldHaveNFilesWithoutFrontmatter(count int) error {
	if dc.scanResult == nil {
		return fmt.Errorf("scan result is nil")
	}
	if dc.scanResult.NoFrontmatterCount != count {
		return fmt.Errorf("expected %d files without frontmatter, got %d", count, dc.scanResult.NoFrontmatterCount)
	}
	return nil
}

func (dc *domainContext) theScanResultShouldIncludeExistingType(typeName string) error {
	if dc.scanResult == nil {
		return fmt.Errorf("scan result is nil")
	}
	for _, t := range dc.scanResult.ExistingTypes {
		if t.Name == typeName {
			return nil
		}
	}
	return fmt.Errorf("existing type %q not found in scan result", typeName)
}

func initImportScanSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^a source directory "([^"]*)" with markdown files:$`, dc.aSourceDirectoryWithMarkdownFiles)
	ctx.Step(`^a source directory "([^"]*)" with no markdown files$`, dc.aSourceDirectoryWithNoMarkdownFiles)
	ctx.Step(`^a source directory "([^"]*)" with plain markdown files:$`, dc.aSourceDirectoryWithPlainMarkdownFiles)
	ctx.Step(`^I scan sources "([^"]*)"$`, dc.iScanSources)
	ctx.Step(`^the scan result should have (\d+) files$`, dc.theScanResultShouldHaveNFiles)
	ctx.Step(`^the scan frontmatter should show key "([^"]*)" appearing (\d+) times$`, dc.theScanFrontmatterShouldShowKeyAppearingNTimes)
	ctx.Step(`^the scan result should have (\d+) files without frontmatter$`, dc.theScanResultShouldHaveNFilesWithoutFrontmatter)
	ctx.Step(`^the scan result should include existing type "([^"]*)"$`, dc.theScanResultShouldIncludeExistingType)
}
