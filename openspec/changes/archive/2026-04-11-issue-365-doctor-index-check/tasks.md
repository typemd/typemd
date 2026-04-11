## 1. Doctor: Index category (BDD first)

- [x] 1.1 Update `core/features/doctor.feature` — change "Healthy vault has all categories passing" to assert **8 categories** and add an `"Index" category should pass` line
- [x] 1.2 Add a new BDD scenario to `core/features/doctor.feature`: "Out-of-band object file triggers index auto-fix" — uses an `out-of-band object "<path>" exists` step, asserts `Index` category has 0 issues and the doctor report has 1 auto-fixed
- [x] 1.3 Add the step definition `anOutOfBandObjectExists(relPath)` in `core/bdd_steps_doctor_test.go` — writes a minimal valid `.md` file (with `name:` frontmatter) directly under `vault.ObjectsDir()/<relPath>` without going through `ObjectService`
- [x] 1.4 Register the new step in `initDoctorSteps`
- [x] 1.5 Run `go test ./core/ -run TestFeatures` and confirm the new scenarios fail for the expected reason (7 categories, not 8 / no "Index" category)

## 2. Doctor: Implement `checkIndexSync`

- [x] 2.1 Add a `checkIndexSync(v *Vault) DoctorCategory` function in `core/doctor.go` that calls `v.Reconcile()`, applies events via `v.Project()`, and returns a `DoctorCategory{Name: "Index"}` with `AutoFixed` set to the number of events (or `Issues` containing a `SeverityError` if Reconcile/Project fails)
- [x] 2.2 Wire `checkIndexSync` into `RunDoctor` **between** `checkCorruptedFiles` (Files) and `checkOrphans` (Orphans) to match the category order documented in `cmd/doctor.go`
- [x] 2.3 Update the `RunDoctor` doc comment to say "8 categories" instead of "7 categories"
- [x] 2.4 Run `go test ./core/ -run TestFeatures` and confirm BDD scenarios now pass
- [x] 2.5 Run `go test ./core/...` and `go vet ./...` to ensure no regressions

## 3. Unit test coverage

- [x] 3.1 Add a unit test in `core/doctor_test.go` (create if missing) that constructs an in-memory-vault, verifies `RunDoctor` returns a category with `Name == "Index"`, and that the category is in the expected position (index 6, between "Files" at 5 and "Orphans" at 7)
- [x] 3.2 Add a unit test asserting that when an out-of-band file is written and `RunDoctor` runs, the Index category's `AutoFixed` count is ≥ 1 and `HasUnresolvedIssues()` returns false

## 4. Verification

- [x] 4.1 Run `make test` (builds frontend + `go build` + `go test ./...` + `go vet ./...`) and confirm green
- [x] 4.2 Manual smoke: `go run ./cmd/tmd doctor -C <tmpvault>` on a clean vault — verify output lists 8 category lines including "✓ Index" and the summary reports "No issues found"
