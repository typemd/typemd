Feature: In-memory filter matching
  MatchFilter evaluates filter rules against objects in memory,
  as a fallback when SQLite is unavailable.

  # ── String operators ───────────────────────────────────────────

  Scenario: "is" operator matches exact value
    Given an object with type "book" and property "status" set to "reading"
    When I match filter property "status" operator "is" value "reading"
    Then the filter should match

  Scenario: "is" operator rejects different value
    Given an object with type "book" and property "status" set to "reading"
    When I match filter property "status" operator "is" value "finished"
    Then the filter should not match

  Scenario: "is_not" operator rejects matching value
    Given an object with type "book" and property "status" set to "reading"
    When I match filter property "status" operator "is_not" value "reading"
    Then the filter should not match

  Scenario: "is_not" operator accepts different value
    Given an object with type "book" and property "status" set to "reading"
    When I match filter property "status" operator "is_not" value "finished"
    Then the filter should match

  Scenario: "contains" operator matches substring
    Given an object with type "book" and property "author" set to "J.R.R. Tolkien"
    When I match filter property "author" operator "contains" value "tolkien"
    Then the filter should match

  Scenario: "contains" operator rejects missing substring
    Given an object with type "book" and property "author" set to "J.R.R. Tolkien"
    When I match filter property "author" operator "contains" value "rowling"
    Then the filter should not match

  Scenario: "does_not_contain" operator accepts missing substring
    Given an object with type "book" and property "author" set to "J.R.R. Tolkien"
    When I match filter property "author" operator "does_not_contain" value "rowling"
    Then the filter should match

  Scenario: "does_not_contain" operator rejects present substring
    Given an object with type "book" and property "author" set to "J.R.R. Tolkien"
    When I match filter property "author" operator "does_not_contain" value "tolkien"
    Then the filter should not match

  Scenario: "starts_with" operator matches prefix
    Given an object with type "book" and property "title" set to "The Lord of the Rings"
    When I match filter property "title" operator "starts_with" value "the lord"
    Then the filter should match

  Scenario: "starts_with" operator rejects non-prefix
    Given an object with type "book" and property "title" set to "The Lord of the Rings"
    When I match filter property "title" operator "starts_with" value "lord"
    Then the filter should not match

  Scenario: "ends_with" operator matches suffix
    Given an object with type "book" and property "title" set to "The Lord of the Rings"
    When I match filter property "title" operator "ends_with" value "rings"
    Then the filter should match

  Scenario: "ends_with" operator rejects non-suffix
    Given an object with type "book" and property "title" set to "The Lord of the Rings"
    When I match filter property "title" operator "ends_with" value "lord"
    Then the filter should not match

  # ── Numeric operators ──────────────────────────────────────────

  Scenario: "eq" operator matches equal numeric value
    Given an object with type "book" and property "rating" set to "4.5"
    When I match filter property "rating" operator "eq" value "4.5"
    Then the filter should match

  Scenario: "neq" operator rejects equal numeric value
    Given an object with type "book" and property "rating" set to "4.5"
    When I match filter property "rating" operator "neq" value "4.5"
    Then the filter should not match

  Scenario: "neq" operator matches different numeric value
    Given an object with type "book" and property "rating" set to "4.5"
    When I match filter property "rating" operator "neq" value "3.0"
    Then the filter should match

  Scenario: "gt" operator matches greater value
    Given an object with type "book" and property "rating" set to "4.5"
    When I match filter property "rating" operator "gt" value "4"
    Then the filter should match

  Scenario: "gt" operator rejects equal value
    Given an object with type "book" and property "rating" set to "4"
    When I match filter property "rating" operator "gt" value "4"
    Then the filter should not match

  Scenario: "gte" operator matches equal value
    Given an object with type "book" and property "rating" set to "4"
    When I match filter property "rating" operator "gte" value "4"
    Then the filter should match

  Scenario: "lt" operator matches lesser value
    Given an object with type "book" and property "rating" set to "3"
    When I match filter property "rating" operator "lt" value "4"
    Then the filter should match

  Scenario: "lte" operator matches equal value
    Given an object with type "book" and property "rating" set to "4"
    When I match filter property "rating" operator "lte" value "4"
    Then the filter should match

  # ── Date operators ─────────────────────────────────────────────

  Scenario: "before" operator matches earlier date
    Given an object with type "book" and property "published" set to "2024-01-15"
    When I match filter property "published" operator "before" value "2025-01-01"
    Then the filter should match

  Scenario: "before" operator rejects later date
    Given an object with type "book" and property "published" set to "2025-06-01"
    When I match filter property "published" operator "before" value "2025-01-01"
    Then the filter should not match

  Scenario: "after" operator matches later date
    Given an object with type "book" and property "published" set to "2025-06-01"
    When I match filter property "published" operator "after" value "2025-01-01"
    Then the filter should match

  Scenario: "on_or_before" operator matches equal date
    Given an object with type "book" and property "published" set to "2025-01-01"
    When I match filter property "published" operator "on_or_before" value "2025-01-01"
    Then the filter should match

  Scenario: "on_or_after" operator matches equal date
    Given an object with type "book" and property "published" set to "2025-01-01"
    When I match filter property "published" operator "on_or_after" value "2025-01-01"
    Then the filter should match

  # ── Empty checks ───────────────────────────────────────────────

  Scenario: "is_empty" matches nil property
    Given an object with type "book" and no property "author"
    When I match filter property "author" operator "is_empty" value ""
    Then the filter should match

  Scenario: "is_empty" matches empty string property
    Given an object with type "book" and property "author" set to ""
    When I match filter property "author" operator "is_empty" value ""
    Then the filter should match

  Scenario: "is_empty" matches "null" string property
    Given an object with type "book" and property "author" set to "null"
    When I match filter property "author" operator "is_empty" value ""
    Then the filter should match

  Scenario: "is_not_empty" matches filled property
    Given an object with type "book" and property "author" set to "Tolkien"
    When I match filter property "author" operator "is_not_empty" value ""
    Then the filter should match

  Scenario: "is_not_empty" rejects nil property
    Given an object with type "book" and no property "author"
    When I match filter property "author" operator "is_not_empty" value ""
    Then the filter should not match

  # ── Type filter ────────────────────────────────────────────────

  Scenario: Type filter matches object type
    Given an object with type "book" and property "status" set to "reading"
    When I match filter property "type" operator "is" value "book"
    Then the filter should match

  Scenario: Type filter rejects different type
    Given an object with type "book" and property "status" set to "reading"
    When I match filter property "type" operator "is" value "person"
    Then the filter should not match

  # ── Multiple filters (AND logic) ──────────────────────────────

  Scenario: Multiple filters all match
    Given an object with type "book" and property "status" set to "reading"
    And the object also has property "rating" set to "5"
    When I match filters "type=book status=reading"
    Then the filter should match

  Scenario: Multiple filters with one mismatch
    Given an object with type "book" and property "status" set to "reading"
    And the object also has property "rating" set to "5"
    When I match filters "type=person status=reading"
    Then the filter should not match
