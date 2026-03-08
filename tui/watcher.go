package tui

import (
	"os"
	"path/filepath"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/fsnotify/fsnotify"
)

// fileChangedMsg is sent when a file change is detected in the objects directory.
type fileChangedMsg struct{}

// watchObjects starts watching the objects directory for changes.
// Returns a tea.Cmd that sends fileChangedMsg when changes are detected.
// Uses debouncing to batch rapid changes (e.g. editor save).
func watchObjects(objectsDir string) tea.Cmd {
	return func() tea.Msg {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return nil
		}

		// Watch objects dir and all subdirectories
		filepath.Walk(objectsDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				watcher.Add(path)
			}
			return nil
		})

		// Wait for first change, then debounce
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return nil
				}
				if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
					// Debounce: wait a bit for more changes
					time.Sleep(200 * time.Millisecond)
					// Drain any queued events
					for len(watcher.Events) > 0 {
						<-watcher.Events
					}
					// Re-add new subdirectories (for Create events)
					filepath.Walk(objectsDir, func(path string, info os.FileInfo, err error) error {
						if err == nil && info.IsDir() {
							watcher.Add(path)
						}
						return nil
					})
					return fileChangedMsg{}
				}
			case _, ok := <-watcher.Errors:
				if !ok {
					return nil
				}
			}
		}
	}
}
