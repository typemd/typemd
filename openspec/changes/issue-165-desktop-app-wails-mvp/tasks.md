## 1. Scaffold Wails project

- [ ] 1.1 Install Wails v3 CLI if not present (`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`)
- [ ] 1.2 Scaffold Wails v3 project under `app/` with vanilla TS frontend template
- [ ] 1.3 Update `go.mod` to include Wails v3 dependency
- [ ] 1.4 Verify `app/` directory structure: `main.go`, `frontend/`, Wails config

## 2. Create AppService backend

- [ ] 2.1 Create `app/service.go` with `AppService` struct wrapping `core.Vault`
- [ ] 2.2 Implement `ListObjects()` method returning objects with type, ID, and display name
- [ ] 2.3 Implement `GetObject(id string)` method returning full object detail
- [ ] 2.4 Wire `AppService` into Wails app binding in `main.go`

## 3. Build frontend object list

- [ ] 3.1 Generate Wails bindings for `AppService` methods
- [ ] 3.2 Create object list page that calls `ListObjects` binding on load
- [ ] 3.3 Display objects grouped by type with count headers
- [ ] 3.4 Handle empty vault state with "no objects" message

## 4. Build and verify

- [ ] 4.1 Verify `go build ./app/` compiles without errors
- [ ] 4.2 Run the binary and confirm the window opens with object list
- [ ] 4.3 Test with a vault containing objects and verify data displays correctly
