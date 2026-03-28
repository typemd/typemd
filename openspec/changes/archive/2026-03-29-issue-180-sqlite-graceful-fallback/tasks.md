## 1. Core: In-memory filter matching

- [x] 1.1 Write BDD scenarios for in-memory filter matching (all operators: is, is_not, contains, does_not_contain, starts_with, ends_with, eq, neq, gt, gte, lt, lte, before, after, on_or_before, on_or_after, is_empty, is_not_empty, type filter)
- [x] 1.2 Implement step definitions for filter matching scenarios
- [x] 1.3 Create `core/filter_match.go` with `MatchFilter(obj *Object, rule FilterRule) bool` and `MatchFilters(obj *Object, rules []FilterRule) bool`
- [x] 1.4 Add unit tests for filter edge cases (nil properties, numeric string comparison, empty values)

## 2. Core: In-memory sort

- [x] 2.1 Write BDD scenarios for in-memory sorting (single sort, multi sort, asc/desc, missing values)
- [x] 2.2 Implement step definitions for sort scenarios
- [x] 2.3 Create `core/sort_objects.go` with `SortObjects(objects []*Object, rules []SortRule)` using sort.SliceStable
- [x] 2.4 Add unit tests for sort edge cases (nil property values, mixed types, stable ordering)

## 3. Core: QueryService fallback for Query()

- [x] 3.1 Write BDD scenarios for query fallback (index unavailable → filesystem scan with filter+sort, warning logged)
- [x] 3.2 Implement step definitions for query fallback scenarios
- [x] 3.3 Add fallback logic in `QueryService.Query()`: catch index error → repo.Walk() + MatchFilters + SortObjects
- [x] 3.4 Add unit tests for query fallback edge cases (empty vault, walk error)

## 4. Core: QueryService fallback for Search()

- [x] 4.1 Write BDD scenarios for search fallback (index unavailable → substring match on name/description/body)
- [x] 4.2 Implement step definitions for search fallback scenarios
- [x] 4.3 Add fallback logic in `QueryService.Search()`: catch index error → repo.Walk() + case-insensitive substring match
- [x] 4.4 Add unit tests for search fallback edge cases (empty keyword, no matches, special characters)

## 5. Core: VaultStats and TypeStats fallback

- [x] 5.1 Write BDD scenarios for VaultStats and TypeStats in fallback mode
- [x] 5.2 Implement step definitions for stats fallback scenarios
- [x] 5.3 Verify VaultStats and TypeStats work via existing Query() fallback (changed VaultStats/TypeStats to use s.Query() instead of s.index.Query())
