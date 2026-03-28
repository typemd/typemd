package web

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/typemd/typemd/core"
)

// Server provides HTTP API endpoints backed by a core.Vault.
type Server struct {
	vault    *core.Vault
	mux      *http.ServeMux
	frontend fs.FS
}

// NewServer creates a new web server wrapping the given vault.
// If frontend is non-nil, it serves the SPA for non-API routes.
func NewServer(vault *core.Vault, frontend ...fs.FS) *Server {
	s := &Server{vault: vault, mux: http.NewServeMux()}
	if len(frontend) > 0 && frontend[0] != nil {
		s.frontend = frontend[0]
	}
	s.routes()
	return s
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /api/types", s.handleListTypes)
	s.mux.HandleFunc("GET /api/types/{name}", s.handleGetType)
	s.mux.HandleFunc("GET /api/objects", s.handleListObjects)
	s.mux.HandleFunc("POST /api/objects", s.handleCreateObject)
	s.mux.HandleFunc("GET /api/objects/{type}/{slug}", s.handleGetObject)
	s.mux.HandleFunc("PUT /api/objects/{type}/{slug}", s.handleUpdateObject)
	s.mux.HandleFunc("GET /api/properties/{type}/{slug}", s.handleGetDisplayProperties)
	s.mux.HandleFunc("PUT /api/properties/{type}/{slug}/{key}", s.handleUpdateProperty)
	s.mux.HandleFunc("GET /api/templates/{type}", s.handleListTemplates)
	s.mux.HandleFunc("GET /api/config", s.handleGetConfig)

	if s.frontend != nil {
		fileServer := http.FileServerFS(s.frontend)
		s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if path == "/" {
				path = "index.html"
			} else {
				path = strings.TrimPrefix(path, "/")
			}
			if _, err := fs.Stat(s.frontend, path); err == nil {
				fileServer.ServeHTTP(w, r)
				return
			}
			// SPA fallback
			r.URL.Path = "/"
			fileServer.ServeHTTP(w, r)
		})
	}
}

type typeItem struct {
	Name   string `json:"name"`
	Plural string `json:"plural"`
	Emoji  string `json:"emoji"`
	Color  string `json:"color"`
	Count  int    `json:"count"`
}

type propertyDef struct {
	Name    string       `json:"name"`
	Type    string       `json:"type"`
	Emoji   string       `json:"emoji,omitempty"`
	Pin     int          `json:"pin,omitempty"`
	Target  string       `json:"target,omitempty"`
	Options []optionItem `json:"options,omitempty"`
}

type optionItem struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type typeDetail struct {
	Name       string        `json:"name"`
	Plural     string        `json:"plural"`
	Emoji      string        `json:"emoji"`
	Color      string        `json:"color"`
	Properties []propertyDef `json:"properties"`
}

type objectItem struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Name   string `json:"name"`
	Locked bool   `json:"locked"`
}

type objectDetail struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Name       string         `json:"name"`
	Properties map[string]any `json:"properties"`
	Body       string         `json:"body"`
	Locked     bool           `json:"locked"`
}

type displayProp struct {
	Key        string `json:"key"`
	Value      any    `json:"value"`
	Display    string `json:"display"`
	Type       string `json:"type"`
	Emoji      string `json:"emoji,omitempty"`
	Pin        int    `json:"pin,omitempty"`
	IsRelation bool   `json:"isRelation,omitempty"`
	IsReverse  bool   `json:"isReverse,omitempty"`
	IsBacklink bool   `json:"isBacklink,omitempty"`
}

type createObjectRequest struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Template string `json:"template,omitempty"`
}

type updateObjectRequest struct {
	Body *string `json:"body,omitempty"`
}

type updatePropertyRequest struct {
	Value string `json:"value"`
}

func newObjectDetail(obj *core.Object) objectDetail {
	return objectDetail{
		ID:         obj.ID,
		Type:       obj.Type,
		Name:       obj.GetName(),
		Properties: serializeProperties(obj.Properties),
		Body:       obj.Body,
		Locked:     obj.IsLocked(),
	}
}

func (s *Server) handleListTypes(w http.ResponseWriter, r *http.Request) {
	names := s.vault.ListTypes()
	items := make([]typeItem, 0, len(names))
	for _, name := range names {
		item := typeItem{Name: name}
		if ts, err := s.vault.LoadType(name); err == nil {
			item.Plural = ts.Plural
			item.Emoji = ts.Emoji
			item.Color = ts.Color
		}
		count, _ := s.vault.CountObjectsByType(name)
		item.Count = count
		items = append(items, item)
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) handleGetType(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if !safePath(name) {
		writeError(w, http.StatusBadRequest, "invalid type name")
		return
	}
	ts, err := s.vault.LoadType(name)
	if err != nil {
		writeError(w, http.StatusNotFound, "type not found")
		return
	}

	props := make([]propertyDef, len(ts.Properties))
	for i, p := range ts.Properties {
		opts := make([]optionItem, len(p.Options))
		for j, o := range p.Options {
			label := o.Label
			if label == "" {
				label = o.Value
			}
			opts[j] = optionItem{Value: o.Value, Label: label}
		}
		props[i] = propertyDef{
			Name:    p.Name,
			Type:    p.Type,
			Emoji:   p.Emoji,
			Pin:     p.Pin,
			Target:  p.Target,
			Options: opts,
		}
	}

	writeJSON(w, http.StatusOK, typeDetail{
		Name:       ts.Name,
		Plural:     ts.Plural,
		Emoji:      ts.Emoji,
		Color:      ts.Color,
		Properties: props,
	})
}

func (s *Server) handleListObjects(w http.ResponseWriter, r *http.Request) {
	typeName := r.URL.Query().Get("type")

	var filter []core.FilterRule
	if typeName != "" {
		filter = core.TypeFilter(typeName)
	}
	sort := core.SortRule{Property: core.NameProperty, Direction: "asc"}

	objects, err := s.vault.QueryObjects(filter, sort)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "query failed")
		return
	}

	items := make([]objectItem, 0, len(objects))
	for _, obj := range objects {
		items = append(items, objectItem{
			ID:     obj.ID,
			Type:   obj.Type,
			Name:   obj.GetName(),
			Locked: obj.IsLocked(),
		})
	}
	writeJSON(w, http.StatusOK, items)
}

// safePath validates a URL path parameter contains no traversal sequences.
func safePath(segment string) bool {
	return !strings.Contains(segment, "..") && !strings.Contains(segment, "/") && !strings.Contains(segment, "\\")
}

func objectID(r *http.Request) (string, bool) {
	t := r.PathValue("type")
	s := r.PathValue("slug")
	if !safePath(t) || !safePath(s) {
		return "", false
	}
	return t + "/" + s, true
}

func (s *Server) handleGetObject(w http.ResponseWriter, r *http.Request) {
	id, ok := objectID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid object ID")
		return
	}
	obj, err := s.vault.GetObject(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "object not found")
		return
	}
	writeJSON(w, http.StatusOK, newObjectDetail(obj))
}

func (s *Server) handleUpdateObject(w http.ResponseWriter, r *http.Request) {
	id, ok := objectID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid object ID")
		return
	}

	var req updateObjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	obj, err := s.vault.GetObject(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "object not found")
		return
	}

	if req.Body != nil {
		obj.Body = *req.Body
	}

	if err := s.vault.SaveObject(obj); err != nil {
		writeError(w, http.StatusInternalServerError, "save failed")
		return
	}
	writeJSON(w, http.StatusOK, newObjectDetail(obj))
}

func (s *Server) handleCreateObject(w http.ResponseWriter, r *http.Request) {
	var req createObjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if req.Type == "" || req.Name == "" {
		writeError(w, http.StatusBadRequest, "type and name are required")
		return
	}

	obj, err := s.vault.Objects.Create(req.Type, req.Name, req.Template)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "create failed")
		return
	}

	writeJSON(w, http.StatusCreated, newObjectDetail(obj))
}

func (s *Server) handleGetDisplayProperties(w http.ResponseWriter, r *http.Request) {
	id, ok := objectID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid object ID")
		return
	}
	obj, err := s.vault.GetObject(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "object not found")
		return
	}

	dps, err := s.vault.BuildDisplayProperties(obj)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "build properties failed")
		return
	}

	result := make([]displayProp, len(dps))
	for i, dp := range dps {
		result[i] = displayProp{
			Key:        dp.Key,
			Value:      serializeValue(dp.Value),
			Display:    dp.FormatValue(),
			Type:       dp.Type,
			Emoji:      dp.Emoji,
			Pin:        dp.Pin,
			IsRelation: dp.IsRelation,
			IsReverse:  dp.IsReverse,
			IsBacklink: dp.IsBacklink,
		}
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleUpdateProperty(w http.ResponseWriter, r *http.Request) {
	id, ok := objectID(r)
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid object ID")
		return
	}
	key := r.PathValue("key")
	if !safePath(key) {
		writeError(w, http.StatusBadRequest, "invalid property key")
		return
	}

	var req updatePropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	obj, err := s.vault.GetObject(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "object not found")
		return
	}

	ts, err := s.vault.LoadType(obj.Type)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "load type failed")
		return
	}

	var prop *core.Property
	for i := range ts.Properties {
		if ts.Properties[i].Name == key {
			prop = &ts.Properties[i]
			break
		}
	}

	if prop != nil {
		if err := core.ValidatePropertyValue(prop.Type, prop.Options, req.Value); err != nil {
			writeError(w, http.StatusBadRequest, "validation: "+err.Error())
			return
		}
	}

	parsed := parsePropertyValue(prop, req.Value)
	if _, err := obj.SetProperty(key, parsed, ts); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := s.vault.SaveObject(obj); err != nil {
		writeError(w, http.StatusInternalServerError, "save failed")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleListTemplates(w http.ResponseWriter, r *http.Request) {
	typeName := r.PathValue("type")
	if !safePath(typeName) {
		writeError(w, http.StatusBadRequest, "invalid type name")
		return
	}
	templates, err := s.vault.ListTemplates(typeName)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "list templates failed")
		return
	}
	writeJSON(w, http.StatusOK, templates)
}

type webConfigResponse struct {
	Theme string `json:"theme"`
}

func (s *Server) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	cfg := s.vault.Config()
	theme := cfg.Web.Theme
	if theme == "" {
		theme = "warm"
	}
	writeJSON(w, http.StatusOK, webConfigResponse{Theme: theme})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func serializeProperties(props map[string]any) map[string]any {
	result := make(map[string]any, len(props))
	for k, v := range props {
		result[k] = serializeValue(v)
	}
	return result
}

func serializeValue(v any) any {
	switch val := v.(type) {
	case time.Time:
		return val.Format(time.RFC3339)
	case []any:
		items := make([]any, len(val))
		for i, item := range val {
			items[i] = serializeValue(item)
		}
		return items
	default:
		return v
	}
}

func parsePropertyValue(prop *core.Property, value string) any {
	if prop == nil {
		return value
	}
	switch prop.Type {
	case "checkbox":
		return value == "true"
	case "number":
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f
		}
		return value
	case "integer":
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
		return value
	case "multi_select":
		value = strings.TrimSpace(value)
		if value == "" {
			return []any{}
		}
		parts := strings.Split(value, ",")
		result := make([]any, len(parts))
		for i, p := range parts {
			result[i] = strings.TrimSpace(p)
		}
		return result
	default:
		return value
	}
}
