Feature: Tag reference resolution
  SyncIndex resolves tag references in object frontmatter by ID or name,
  auto-creates missing tags, and writes tag relations.

  Background:
    Given a vault is ready

  Scenario: Tag reference resolved by full ID
    Given a "tag" object named "go" exists
    And a "book" object named "golang-book" exists with tag reference by ID
    When I sync the index
    Then the book should have a tag relation to the tag

  Scenario: Tag reference resolved by name
    Given a "tag" object named "go" exists
    And a "book" object named "golang-book" exists with tag reference by name "go"
    When I sync the index
    Then the book should have a tag relation to the tag

  Scenario: Missing tag is auto-created
    Given a "book" object named "golang-book" exists with tag reference by name "auto-tag"
    When I sync the index
    Then a tag object named "auto-tag" should exist on disk
