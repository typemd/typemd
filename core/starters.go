package core

import (
	"embed"
	"fmt"
	"sync"
)

//go:embed starters/*.yaml
var starterFS embed.FS

// StarterType describes an available starter type template.
type StarterType struct {
	Name        string
	Emoji       string
	Description string
	YAML        []byte
}

// starterMeta holds display metadata for each starter type.
// The order here determines the display order in the picker.
var starterMeta = []struct {
	name        string
	filename    string
	emoji       string
	description string
}{
	{"idea", "idea.yaml", "💡", "Capture and develop ideas"},
	{"note", "note.yaml", "📝", "Quick notes and thoughts"},
	{"book", "book.yaml", "📚", "Track your reading"},
}

var (
	starterOnce  sync.Once
	starterCache []StarterType
)

// StarterTypes returns all available starter type templates.
// Results are cached after the first call since the data is embedded and immutable.
func StarterTypes() []StarterType {
	starterOnce.Do(func() {
		starterCache = make([]StarterType, 0, len(starterMeta))
		for _, m := range starterMeta {
			data, err := starterFS.ReadFile("starters/" + m.filename)
			if err != nil {
				panic(fmt.Sprintf("embedded starter %s unreadable: %v", m.filename, err))
			}
			starterCache = append(starterCache, StarterType{
				Name:        m.name,
				Emoji:       m.emoji,
				Description: m.description,
				YAML:        data,
			})
		}
	})
	return starterCache
}
