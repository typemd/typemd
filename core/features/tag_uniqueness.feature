Feature: Tag name uniqueness
  Tag names must be unique across the vault to enable unambiguous
  name-based resolution.

  Background:
    Given a vault is ready

  Scenario: Creating a tag with a unique name succeeds
    When I create a "tag" object named "go"
    Then no error should occur

  Scenario: Creating a tag with a duplicate name fails
    Given a "tag" object named "go" exists
    When I create a "tag" object named "go"
    Then an error should occur

  Scenario: Validation reports duplicate tag names
    Given a "tag" object named "go" exists
    And a raw duplicate tag named "go" exists
    When I validate tag name uniqueness
    Then there should be tag uniqueness errors
