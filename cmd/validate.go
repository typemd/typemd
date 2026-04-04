package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/typemd/typemd/core"
)

var watchFlag bool

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate type schemas, objects, and relations",
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath)
		if err != nil {
			return err
		}
		defer vault.Close()

		if !watchFlag {
			return runValidation(vault)
		}

		ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		return runWatchValidation(ctx, vault)
	},
}

func init() {
	validateCmd.Flags().BoolVarP(&watchFlag, "watch", "w", false, "continuously watch for changes and re-validate")
	typeCmd.AddCommand(validateCmd)
}

func runValidation(vault *core.Vault) error {
	totalErrors := 0

	// Phase 1: Schema validation
	schemaErrs := core.ValidateAllSchemas(vault)
	if len(schemaErrs) > 0 {
		fmt.Println("Schema errors:")
		schemaNames := make([]string, 0, len(schemaErrs))
		for name := range schemaErrs {
			schemaNames = append(schemaNames, name)
		}
		sort.Strings(schemaNames)
		for _, name := range schemaNames {
			for _, e := range schemaErrs[name] {
				fmt.Printf("  %s: %s\n", name, e)
				totalErrors++
			}
		}
		fmt.Println()
	}

	// Phase 2: Object validation
	objectErrs := core.ValidateAllObjects(vault)
	if len(objectErrs) > 0 {
		fmt.Println("Object errors:")
		objectIDs := make([]string, 0, len(objectErrs))
		for id := range objectErrs {
			objectIDs = append(objectIDs, id)
		}
		sort.Strings(objectIDs)
		for _, id := range objectIDs {
			for _, e := range objectErrs[id] {
				fmt.Printf("  %s: %s\n", id, e)
				totalErrors++
			}
		}
		fmt.Println()
	}

	// Phase 3: Relation validation
	relationErrs := core.ValidateRelations(vault)
	if len(relationErrs) > 0 {
		fmt.Println("Relation errors:")
		for _, e := range relationErrs {
			fmt.Printf("  %s\n", e)
			totalErrors++
		}
		fmt.Println()
	}

	// Phase 3b: Unresolved relation references
	refErrs := core.ValidateRelationReferences(vault)
	if len(refErrs) > 0 {
		fmt.Println("Unresolved relation references:")
		for _, e := range refErrs {
			fmt.Printf("  %s\n", e)
			totalErrors++
		}
		fmt.Println()
	}

	// Phase 4: Wiki-link validation
	wikiLinkErrs := core.ValidateWikiLinks(vault)
	if len(wikiLinkErrs) > 0 {
		fmt.Println("Wiki-link errors:")
		for _, e := range wikiLinkErrs {
			fmt.Printf("  %s\n", e)
			totalErrors++
		}
		fmt.Println()
	}

	// Phase 5: Name uniqueness validation (all unique types)
	nameErrs := core.ValidateNameUniqueness(vault)
	if len(nameErrs) > 0 {
		fmt.Println("Name uniqueness errors:")
		for _, e := range nameErrs {
			fmt.Printf("  %s\n", e)
			totalErrors++
		}
		fmt.Println()
	}

	if totalErrors > 0 {
		return fmt.Errorf("found %d validation error(s)", totalErrors)
	}

	fmt.Println("Validation passed.")
	return nil
}

func runWatchValidation(ctx context.Context, vault *core.Vault) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("creating file watcher: %w", err)
	}
	defer watcher.Close()

	addWatchPaths(watcher, vault)

	clearTerminal()
	printWatchHeader()
	_ = runValidation(vault) // errors are printed; watch continues regardless

	const debounceMs = 200

	for {
		select {
		case <-ctx.Done():
			return nil
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) == 0 {
				continue
			}

			// Debounce: collect rapid changes into a single validation run
			select {
			case <-time.After(time.Duration(debounceMs) * time.Millisecond):
			case <-ctx.Done():
				return nil
			}

			// Drain queued events
			for len(watcher.Events) > 0 {
				<-watcher.Events
			}

			clearTerminal()
			printWatchHeader()
			events, _, syncErr := vault.Reconcile()
			if syncErr != nil {
				fmt.Printf("Sync error: %v\n\n", syncErr)
			} else {
				_ = vault.Project(events)
			}
			_ = runValidation(vault)

		case _, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
		}
	}
}

func addWatchPaths(watcher *fsnotify.Watcher, vault *core.Vault) {
	watchDirRecursive(watcher, vault.TypesDir())
	watchDirRecursive(watcher, vault.SharedPropertiesDir())
	watchDirRecursive(watcher, vault.ObjectsDir())
}

// watchDirRecursive adds a directory and all its subdirectories to the watcher.
// Missing directories are silently skipped.
func watchDirRecursive(watcher *fsnotify.Watcher, dir string) {
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return
	}
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			watcher.Add(path)
		}
		return nil
	})
}

func clearTerminal() {
	fmt.Print("\033[H\033[2J")
}

func printWatchHeader() {
	fmt.Printf("[%s] Validating...\n\n", time.Now().Format("15:04:05"))
}
