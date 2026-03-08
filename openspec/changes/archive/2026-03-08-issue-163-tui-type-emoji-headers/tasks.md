## 1. Add emoji to typeGroup struct

- [ ] 1.1 Add `Emoji` field to `typeGroup` struct in `tui/app.go`
- [ ] 1.2 Populate `Emoji` in `buildGroups()` by calling `vault.LoadType()` for each type

## 2. Update group header rendering

- [ ] 2.1 Modify `renderList()` in `tui/list.go` to include emoji prefix when present
- [ ] 2.2 Ensure no extra spacing when emoji is absent

## 3. Testing

- [ ] 3.1 Add unit test for `renderList` with emoji in group header
- [ ] 3.2 Add unit test for `renderList` without emoji (no visual change)
- [ ] 3.3 Verify `go test ./...` and `go build ./...` pass
