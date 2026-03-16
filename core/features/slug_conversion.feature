Feature: Slug conversion
  As a user
  I want to enter natural-language names
  So that they are automatically converted to valid slugs

  Scenario: Simple name with spaces
    When I slugify "Some Thought"
    Then the slug should be "some-thought"

  Scenario: Name with mixed case
    When I slugify "Clean Code"
    Then the slug should be "clean-code"

  Scenario: Name with underscores
    When I slugify "my_great_idea"
    Then the slug should be "my-great-idea"

  Scenario: Name with special characters
    When I slugify "What's the plan?"
    Then the slug should be "whats-the-plan"

  Scenario: Name with consecutive spaces
    When I slugify "too   many   spaces"
    Then the slug should be "too-many-spaces"

  Scenario: Name with leading and trailing whitespace
    When I slugify "  padded name  "
    Then the slug should be "padded-name"

  Scenario: Already slugified input is idempotent
    When I slugify "clean-code"
    Then the slug should be "clean-code"

  Scenario: Name with numbers
    When I slugify "Chapter 3 Notes"
    Then the slug should be "chapter-3-notes"

  Scenario: Non-ASCII letters are preserved
    When I slugify "café latte"
    Then the slug should be "café-latte"

  Scenario: CJK characters are preserved
    When I slugify "我的日記"
    Then the slug should be "我的日記"
