package tui

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/fsnotify/fsnotify"
)

const defaultDebounceMs = 200

// fileChangedMsg is sent when file changes are detected in the objects directory.
// Paths contains the deduplicated list of changed file paths.
type fileChangedMsg struct {
	Paths []string
}

// schemaChangedMsg is sent when file changes are detected in the types directory.
type schemaChangedMsg struct{}

// deduplicatePaths returns a deduplicated copy of the input paths.
func deduplicatePaths(paths []string) []string {
	seen := make(map[string]bool, len(paths))
	result := make([]string, 0, len(paths))
	for _, p := range paths {
		if !seen[p] {
			seen[p] = true
			result = append(result, p)
		}
	}
	return result
}

// watchDir creates a watcher on a directory and calls buildMsg when changes are detected.
// It debounces rapid changes and collects changed file paths during the window.
func watchDir(dir string, debounceMs int, buildMsg func(paths []string) tea.Msg) tea.Cmd {
	if debounceMs <= 0 {
		debounceMs = defaultDebounceMs
	}
	return func() tea.Msg {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return nil
		}
		defer watcher.Close()

		// Watch dir and all subdirectories
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				watcher.Add(path)
			}
			return nil
		})

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return nil
				}
				if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
					// Collect this event's path
					paths := []string{event.Name}

					// Debounce: wait for more changes
					time.Sleep(time.Duration(debounceMs) * time.Millisecond)

					// Drain queued events, collecting their paths
					for len(watcher.Events) > 0 {
						e := <-watcher.Events
						if e.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
							paths = append(paths, e.Name)
						}
					}

					return buildMsg(deduplicatePaths(paths))
				}
			case _, ok := <-watcher.Errors:
				if !ok {
					return nil
				}
			}
		}
	}
}

// watchObjects starts watching the objects directory for changes.
// Returns a tea.Cmd that sends fileChangedMsg with changed .md file paths.
// Non-.md paths are filtered out; if no .md paths remain, Paths will be nil
// to signal that a full sync should be used.
func watchObjects(objectsDir string, debounceMs int) tea.Cmd {
	return watchDir(objectsDir, debounceMs, func(paths []string) tea.Msg {
		var mdPaths []string
		for _, p := range paths {
			if strings.HasSuffix(p, ".md") {
				mdPaths = append(mdPaths, p)
			}
		}
		return fileChangedMsg{Paths: mdPaths}
	})
}

// watchTypes starts watching the types directory for schema changes.
// Returns a tea.Cmd that sends schemaChangedMsg when changes are detected.
func watchTypes(typesDir string, debounceMs int) tea.Cmd {
	return watchDir(typesDir, debounceMs, func(_ []string) tea.Msg {
		return schemaChangedMsg{}
	})
}
