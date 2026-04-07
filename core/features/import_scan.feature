Feature: Import scan
  Scan source directories for markdown files, extracting frontmatter patterns
  and collecting file statistics to support onboarding workflows.

  # ── Basic scanning ─────────────────────────────────────────

  Scenario: Scan a directory with markdown files
    Given a vault is ready
    And a source directory "notes" with markdown files:
      | filename     | title       |
      | alpha.md     | Alpha Note  |
      | beta.md      | Beta Note   |
    When I scan sources "notes"
    Then no error should occur
    And the scan result should have 2 files

  Scenario: Scan multiple source directories
    Given a vault is ready
    And a source directory "notes" with markdown files:
      | filename     | title       |
      | note-one.md  | Note One    |
    And a source directory "docs" with markdown files:
      | filename     | title       |
      | doc-one.md   | Doc One     |
      | doc-two.md   | Doc Two     |
    When I scan sources "notes,docs"
    Then no error should occur
    And the scan result should have 3 files

  Scenario: Scan a non-existent path
    Given a vault is ready
    When I scan sources "nonexistent"
    Then an error should occur

  Scenario: Scan a directory with no markdown files
    Given a vault is ready
    And a source directory "images" with no markdown files
    When I scan sources "images"
    Then no error should occur
    And the scan result should have 0 files

  # ── Frontmatter extraction ────────────────────────────────

  Scenario: Extract frontmatter key patterns
    Given a vault is ready
    And a source directory "books" with markdown files:
      | filename      | title            | author         |
      | clean-code.md | Clean Code       | Robert Martin  |
      | go-book.md    | Go in Action     | William Kennedy|
      | rust-note.md  | Rust Notes       |                |
    When I scan sources "books"
    Then no error should occur
    And the scan frontmatter should show key "title" appearing 3 times
    And the scan frontmatter should show key "author" appearing 2 times

  Scenario: Scan files without frontmatter
    Given a vault is ready
    And a source directory "raw" with plain markdown files:
      | filename   | body            |
      | readme.md  | Hello world     |
      | notes.md   | Some notes here |
    When I scan sources "raw"
    Then no error should occur
    And the scan result should have 2 files
    And the scan result should have 2 files without frontmatter

  # ── Existing vault types ──────────────────────────────────

  Scenario: Scan includes existing vault types
    Given a vault is ready
    When I scan sources "objects"
    Then no error should occur
    And the scan result should include existing type "book"
    And the scan result should include existing type "tag"
