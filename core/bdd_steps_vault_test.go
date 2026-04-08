package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
)

// ── Vault steps ─────────────────────────────────────────────────────────────

func (dc *domainContext) setupVaultDir() {
	dc.rootDir = filepath.Join(os.TempDir(), "typemd-bdd-"+mustULID())
	os.MkdirAll(dc.rootDir, 0755)
	dc.vault = NewVault(dc.rootDir)
}

func (dc *domainContext) iInitializeANewVault() {
	dc.setupVaultDir()
	dc.lastErr = dc.vault.Init()
}

func (dc *domainContext) aVaultIsInitialized() {
	dc.setupVaultDir()
	if err := dc.vault.Init(); err != nil {
		panic(fmt.Sprintf("vault init failed: %v", err))
	}
}

func (dc *domainContext) theVaultDirectoryStructureShouldExist() error {
	for _, d := range []string{dc.vault.TypesDir(), dc.vault.ObjectsDir()} {
		if _, err := os.Stat(d); os.IsNotExist(err) {
			return fmt.Errorf("expected directory %s to exist", d)
		}
	}
	return nil
}

func (dc *domainContext) theSQLiteDatabaseShouldExist() error {
	if _, err := os.Stat(dc.vault.DBPath()); os.IsNotExist(err) {
		return fmt.Errorf("expected index.db to exist")
	}
	return nil
}

func (dc *domainContext) iInitializeTheVaultAgain() {
	dc.lastErr = dc.vault.Init()
}

func (dc *domainContext) iOpenTheVault() {
	dc.lastErr = dc.vault.Open()
}

func (dc *domainContext) iCloseTheVault() {
	if dc.lastErr == nil {
		dc.lastErr = dc.vault.Close()
	}
}

func (dc *domainContext) iOpenAnUninitializedVault() {
	dc.setupVaultDir()
	dc.lastErr = dc.vault.Open()
}

func (dc *domainContext) anObjectFileExistsOnDisk(relPath, title string) {
	fullPath := filepath.Join(dc.rootDir, "objects", relPath)
	os.MkdirAll(filepath.Dir(fullPath), 0755)
	content := fmt.Sprintf("---\ntitle: %s\n---\nHello world\n", title)
	os.WriteFile(fullPath, []byte(content), 0644)
}

func (dc *domainContext) theIndexShouldContainNObjects(expected int) error {
	var count int
	if err := dc.vault.db.QueryRow("SELECT COUNT(*) FROM objects").Scan(&count); err != nil {
		return fmt.Errorf("count query error: %v", err)
	}
	if count != expected {
		return fmt.Errorf("objects count = %d, want %d", count, expected)
	}
	return nil
}

func initVaultSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^I initialize a new vault$`, dc.iInitializeANewVault)
	ctx.Step(`^a vault is initialized$`, dc.aVaultIsInitialized)
	ctx.Step(`^the vault directory structure should exist$`, dc.theVaultDirectoryStructureShouldExist)
	ctx.Step(`^the SQLite database should exist$`, dc.theSQLiteDatabaseShouldExist)
	ctx.Step(`^I initialize the vault again$`, dc.iInitializeTheVaultAgain)
	ctx.Step(`^I open the vault$`, dc.iOpenTheVault)
	ctx.Step(`^I close the vault$`, dc.iCloseTheVault)
	ctx.Step(`^I open an uninitialized vault$`, dc.iOpenAnUninitializedVault)
	ctx.Step(`^an object file exists on disk at "([^"]*)" with title "([^"]*)"$`, dc.anObjectFileExistsOnDisk)
	ctx.Step(`^the index should contain (\d+) objects?$`, dc.theIndexShouldContainNObjects)
}
