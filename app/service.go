package main

import (
	"fmt"

	"github.com/typemd/typemd/core"
)

// ObjectItem represents an object for the frontend.
type ObjectItem struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	DisplayName string `json:"displayName"`
	Emoji       string `json:"emoji"`
}

// ObjectDetail represents a full object for the frontend.
type ObjectDetail struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"`
	DisplayName string            `json:"displayName"`
	Properties  map[string]any    `json:"properties"`
	Body        string            `json:"body"`
}

// AppService wraps core.Vault and exposes methods for the Wails frontend.
type AppService struct {
	vaultPath string
	vault     *core.Vault
}

// NewAppService creates a new AppService with the given vault path.
func NewAppService(vaultPath string) *AppService {
	return &AppService{vaultPath: vaultPath}
}

// Init initializes the vault. Called when the app starts.
func (s *AppService) Init() error {
	v := core.NewVault(s.vaultPath)
	if err := v.Open(); err != nil {
		return fmt.Errorf("open vault: %w", err)
	}
	s.vault = v
	return nil
}

// ListObjects returns all objects grouped for the frontend.
func (s *AppService) ListObjects() ([]ObjectItem, error) {
	if s.vault == nil {
		if err := s.Init(); err != nil {
			return nil, err
		}
	}

	objects, err := s.vault.QueryObjects("")
	if err != nil {
		return nil, fmt.Errorf("query objects: %w", err)
	}

	items := make([]ObjectItem, 0, len(objects))
	for _, obj := range objects {
		var emoji string
		if ts, err := s.vault.LoadType(obj.Type); err == nil {
			emoji = ts.Emoji
		}
		items = append(items, ObjectItem{
			ID:          obj.ID,
			Type:        obj.Type,
			DisplayName: obj.DisplayName(),
			Emoji:       emoji,
		})
	}
	return items, nil
}

// GetObject returns the full detail of an object by ID.
func (s *AppService) GetObject(id string) (*ObjectDetail, error) {
	if s.vault == nil {
		if err := s.Init(); err != nil {
			return nil, err
		}
	}

	obj, err := s.vault.GetObject(id)
	if err != nil {
		return nil, fmt.Errorf("get object: %w", err)
	}

	return &ObjectDetail{
		ID:          obj.ID,
		Type:        obj.Type,
		DisplayName: obj.DisplayName(),
		Properties:  obj.Properties,
		Body:        obj.Body,
	}, nil
}
