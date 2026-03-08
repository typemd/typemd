Feature: Resolve object by prefix
  Users can type a shortened object ID (e.g. "book/clean-code") and have it
  resolve to the full ULID-suffixed ID.

  Background:
    Given a vault is ready

  Scenario: Resolve object by exact full ID
    Given a "book" object named "clean-code" exists
    When I resolve the object by its full ID
    Then the resolved ID should match the original

  Scenario: Resolve object by name prefix without ULID
    Given a "book" object named "clean-code" exists
    When I resolve the object by prefix "book/clean-code"
    Then the resolved object should match the created one

  Scenario: Resolve object by partial ULID prefix
    Given a "book" object named "clean-code" exists
    When I resolve the object by a partial ULID prefix
    Then the resolved object should match the created one

  Scenario: Ambiguous prefix matches multiple objects
    Given a "book" object named "clean-code" exists
    And a "book" object named "clean-code-second-edition" exists
    When I resolve the object by prefix "book/clean-code"
    Then an ambiguous match error should occur with 2 candidates

  Scenario: No object matches prefix
    When I resolve the object by prefix "book/nonexistent"
    Then an error should occur
