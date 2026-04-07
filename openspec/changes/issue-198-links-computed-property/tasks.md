## 1. Core: Add IsLink to DisplayProperty and FormatValue

- [x] 1.1 Write BDD scenarios for links display property (links appear, format, ordering)
- [x] 1.2 Implement BDD step definitions for links display scenarios
- [x] 1.3 Add `IsLink` field to `DisplayProperty` struct and update `FormatValue()` to handle it
- [x] 1.4 Add unit tests for FormatValue with IsLink (arrow prefix, display ID stripping)

## 2. Core: Integrate links into BuildDisplayProperties

- [x] 2.1 Add links query to `QueryService.BuildDisplayProperties` (between reverse relations and backlinks)
- [x] 2.2 Make BDD scenarios pass (links appear in display properties with correct fields and ordering)

## 3. TUI: Ensure links are read-only in property editor

- [x] 3.1 Add unit test verifying links display properties are not editable in prop editor
- [x] 3.2 Update prop editor read-only check to include `IsLink` (if not already covered by existing logic)
