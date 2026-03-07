package core

import (
	"os"
	"strings"
	"testing"
)

func TestWriteFrontmatter(t *testing.T) {
	props := map[string]any{
		"title":  "Go in Action",
		"rating": 4.5,
	}

	data, err := writeFrontmatter(props, "", nil)
	if err != nil {
		t.Fatalf("writeFrontmatter() error = %v", err)
	}

	content := string(data)
	if !strings.HasPrefix(content, "---\n") {
		t.Error("expected content to start with ---")
	}
	if !strings.Contains(content, "title: Go in Action") {
		t.Error("expected content to contain title")
	}
}

func TestWriteFrontmatter_WithBody(t *testing.T) {
	props := map[string]any{"title": "Test"}

	data, err := writeFrontmatter(props, "Hello world", nil)
	if err != nil {
		t.Fatalf("writeFrontmatter() error = %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "---\n\nHello world") {
		t.Errorf("expected body after frontmatter, got:\n%s", content)
	}
}

func TestWriteFrontmatter_EmptyProps(t *testing.T) {
	data, err := writeFrontmatter(map[string]any{}, "", nil)
	if err != nil {
		t.Fatalf("writeFrontmatter() error = %v", err)
	}

	if string(data) != "---\n---\n" {
		t.Errorf("expected empty frontmatter, got:\n%s", string(data))
	}
}

func TestParseFrontmatter(t *testing.T) {
	content := "---\ntitle: Go in Action\nrating: 4.5\n---\n\nSome body here."

	props, body, err := parseFrontmatter([]byte(content))
	if err != nil {
		t.Fatalf("parseFrontmatter() error = %v", err)
	}

	if props["title"] != "Go in Action" {
		t.Errorf("title = %v, want %q", props["title"], "Go in Action")
	}
	if body != "Some body here." {
		t.Errorf("body = %q, want %q", body, "Some body here.")
	}
}

func TestParseFrontmatter_EmptyBody(t *testing.T) {
	content := "---\ntitle: Test\n---\n"

	props, body, err := parseFrontmatter([]byte(content))
	if err != nil {
		t.Fatalf("parseFrontmatter() error = %v", err)
	}

	if props["title"] != "Test" {
		t.Errorf("title = %v, want %q", props["title"], "Test")
	}
	if body != "" {
		t.Errorf("body = %q, want empty", body)
	}
}

func TestParseFrontmatter_NullValues(t *testing.T) {
	content := "---\ntitle: null\nrating: null\n---\n"

	props, _, err := parseFrontmatter([]byte(content))
	if err != nil {
		t.Fatalf("parseFrontmatter() error = %v", err)
	}

	if props["title"] != nil {
		t.Errorf("title = %v, want nil", props["title"])
	}
}

func TestFrontmatter_Roundtrip(t *testing.T) {
	original := map[string]any{
		"title":  "Go in Action",
		"status": "reading",
	}
	body := "Some notes here."

	data, err := writeFrontmatter(original, body, nil)
	if err != nil {
		t.Fatalf("writeFrontmatter() error = %v", err)
	}

	props, parsedBody, err := parseFrontmatter(data)
	if err != nil {
		t.Fatalf("parseFrontmatter() error = %v", err)
	}

	if props["title"] != "Go in Action" {
		t.Errorf("title = %v, want %q", props["title"], "Go in Action")
	}
	if props["status"] != "reading" {
		t.Errorf("status = %v, want %q", props["status"], "reading")
	}
	if parsedBody != body {
		t.Errorf("body = %q, want %q", parsedBody, body)
	}
}

func TestWriteFrontmatter_KeyOrder(t *testing.T) {
	props := map[string]any{
		"title":  "Go in Action",
		"status": "reading",
		"rating": 5,
	}

	data, err := writeFrontmatter(props, "", []string{"title", "status", "rating"})
	if err != nil {
		t.Fatalf("writeFrontmatter() error = %v", err)
	}

	content := string(data)
	titleIdx := strings.Index(content, "title:")
	statusIdx := strings.Index(content, "status:")
	ratingIdx := strings.Index(content, "rating:")

	if titleIdx == -1 || statusIdx == -1 || ratingIdx == -1 {
		t.Fatalf("missing keys in output:\n%s", content)
	}
	if titleIdx >= statusIdx || statusIdx >= ratingIdx {
		t.Errorf("keys not in schema order, got:\n%s", content)
	}
}

func TestVault_NewObject(t *testing.T) {
	v := setupTestVault(t)

	obj, err := v.NewObject("book", "golang-in-action")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	// Filename should have ULID suffix
	if !strings.HasPrefix(obj.Filename, "golang-in-action-") {
		t.Errorf("Filename should start with 'golang-in-action-', got %q", obj.Filename)
	}
	ulidPart := strings.TrimPrefix(obj.Filename, "golang-in-action-")
	if len(ulidPart) != 26 {
		t.Errorf("ULID part length = %d, want 26", len(ulidPart))
	}
	if obj.ID != "book/"+obj.Filename {
		t.Errorf("ID = %q, want %q", obj.ID, "book/"+obj.Filename)
	}
	if obj.Type != "book" {
		t.Errorf("Type = %q, want %q", obj.Type, "book")
	}

	// .md 檔存在
	if _, err := os.Stat(v.ObjectPath(obj.Type, obj.Filename)); os.IsNotExist(err) {
		t.Error("expected .md file to exist")
	}

	// Properties 從 schema 產生（book 有 title, status, rating）
	if len(obj.Properties) != 3 {
		t.Errorf("len(Properties) = %d, want 3", len(obj.Properties))
	}
}

func TestVault_NewObject_SameNameDifferentULID(t *testing.T) {
	v := setupTestVault(t)

	obj1, err := v.NewObject("book", "test")
	if err != nil {
		t.Fatalf("first NewObject() error = %v", err)
	}

	obj2, err := v.NewObject("book", "test")
	if err != nil {
		t.Fatalf("second NewObject() should succeed with ULID, error = %v", err)
	}

	if obj1.ID == obj2.ID {
		t.Errorf("two objects with same name should have different IDs, both got %q", obj1.ID)
	}
	if obj1.Filename == obj2.Filename {
		t.Errorf("two objects with same name should have different filenames, both got %q", obj1.Filename)
	}
}

func TestVault_NewObject_UnknownType(t *testing.T) {
	v := setupTestVault(t)

	_, err := v.NewObject("nonexistent", "test")
	if err == nil {
		t.Fatal("expected error for unknown type, got nil")
	}
}

func TestVault_NewObject_DBNotOpen(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	// 不呼叫 Open()

	_, err := v.NewObject("book", "test")
	if err == nil {
		t.Fatal("expected error when DB not opened, got nil")
	}
}

func TestVault_NewObject_DefaultValues(t *testing.T) {
	v := setupTestVault(t)

	// 寫入有 default 的自訂 schema
	yamlContent := []byte(`name: task
properties:
  - name: title
    type: string
  - name: status
    type: enum
    values: [todo, doing, done]
    default: todo
`)
	if err := os.WriteFile(v.TypesDir()+"/task.yaml", yamlContent, 0644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}

	obj, err := v.NewObject("task", "my-task")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	// Filename should have ULID suffix
	if !strings.HasPrefix(obj.Filename, "my-task-") {
		t.Errorf("Filename should start with 'my-task-', got %q", obj.Filename)
	}

	// title 沒有 default → nil
	if obj.Properties["title"] != nil {
		t.Errorf("title = %v, want nil", obj.Properties["title"])
	}

	// status 有 default → "todo"
	if obj.Properties["status"] != "todo" {
		t.Errorf("status = %v, want %q", obj.Properties["status"], "todo")
	}
}

func TestVault_GetObject(t *testing.T) {
	v := setupTestVault(t)

	created, err := v.NewObject("book", "golang-in-action")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	obj, err := v.GetObject(created.ID)
	if err != nil {
		t.Fatalf("GetObject() error = %v", err)
	}

	if obj.ID != created.ID {
		t.Errorf("ID = %q, want %q", obj.ID, created.ID)
	}
	if obj.Type != "book" {
		t.Errorf("Type = %q, want %q", obj.Type, "book")
	}
	if obj.Filename != created.Filename {
		t.Errorf("Filename = %q, want %q", obj.Filename, created.Filename)
	}
}

func TestVault_GetObject_NotFound(t *testing.T) {
	v := setupTestVault(t)

	_, err := v.GetObject("book/nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent object, got nil")
	}
}

func TestVault_GetObject_InvalidID(t *testing.T) {
	v := setupTestVault(t)

	_, err := v.GetObject("invalid-id-no-slash")
	if err == nil {
		t.Fatal("expected error for invalid ID, got nil")
	}
}

func TestVault_SetProperty(t *testing.T) {
	v := setupTestVault(t)

	created, err := v.NewObject("book", "test")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	if err := v.SetProperty(created.ID, "title", "Go in Action"); err != nil {
		t.Fatalf("SetProperty() error = %v", err)
	}

	// 重新讀取驗證
	obj, err := v.GetObject(created.ID)
	if err != nil {
		t.Fatalf("GetObject() error = %v", err)
	}
	if obj.Properties["title"] != "Go in Action" {
		t.Errorf("title = %v, want %q", obj.Properties["title"], "Go in Action")
	}
}

func TestVault_SetProperty_ValidationError(t *testing.T) {
	v := setupTestVault(t)

	created, err := v.NewObject("book", "test")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	// rating 應該是 number，給 string 應該報錯
	err = v.SetProperty(created.ID, "rating", "not-a-number")
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestVault_SetProperty_ExtraKey(t *testing.T) {
	v := setupTestVault(t)

	created, err := v.NewObject("book", "test")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	// custom_field 不在 schema 中，寬鬆模式允許寫入
	if err := v.SetProperty(created.ID, "custom_field", "anything"); err != nil {
		t.Fatalf("SetProperty() error = %v", err)
	}

	obj, err := v.GetObject(created.ID)
	if err != nil {
		t.Fatalf("GetObject() error = %v", err)
	}
	if obj.Properties["custom_field"] != "anything" {
		t.Errorf("custom_field = %v, want %q", obj.Properties["custom_field"], "anything")
	}
}

func TestVault_SetProperty_NotFound(t *testing.T) {
	v := setupTestVault(t)

	err := v.SetProperty("book/nonexistent", "title", "test")
	if err == nil {
		t.Fatal("expected error for nonexistent object, got nil")
	}
}

func TestVault_SetProperty_DBNotOpen(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	err := v.SetProperty("book/test", "title", "test")
	if err == nil {
		t.Fatal("expected error when DB not opened, got nil")
	}
}


func TestVault_SaveObject_WritesFileAndUpdatesDB(t *testing.T) {
	v := setupTestVault(t)
	obj, err := v.NewObject("book", "test-book")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	obj.Body = "New body content"
	obj.Properties["title"] = "Updated Title"

	if err := v.SaveObject(obj); err != nil {
		t.Fatalf("SaveObject() error = %v", err)
	}

	// File should contain updated content
	data, err := os.ReadFile(v.ObjectPath(obj.Type, obj.Filename))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "Updated Title") {
		t.Error("file should contain updated title")
	}
	if !strings.Contains(content, "New body content") {
		t.Error("file should contain updated body")
	}

	// DB should be updated
	reloaded, err := v.GetObject(obj.ID)
	if err != nil {
		t.Fatalf("GetObject() error = %v", err)
	}
	if reloaded.Body != "New body content" {
		t.Errorf("GetObject().Body = %q, want %q", reloaded.Body, "New body content")
	}
}

func TestVault_SaveObject_ErrorWhenNotOpened(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	obj := &Object{ID: "book/test", Type: "book", Filename: "test"}
	if err := v.SaveObject(obj); err == nil {
		t.Fatal("expected error when vault not opened, got nil")
	}
}
