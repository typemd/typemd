Feature: Doctor
  The doctor command performs a comprehensive vault health check, reporting issues
  across multiple categories and auto-fixing what it can.

  Scenario: Healthy vault has all categories passing
    Given a vault is ready
    When I run doctor
    Then the doctor report should have 8 categories
    And the "Schemas" category should pass
    And the "Objects" category should pass
    And the "Relations" category should pass
    And the "Wiki-links" category should pass
    And the "Uniqueness" category should pass
    And the "Files" category should pass
    And the "Index" category should pass
    And the "Orphans" category should pass
    And the doctor report should have 0 total issues

  Scenario: Corrupted object file is detected
    Given a vault is ready
    And a corrupted object file exists at "book/bad-file.md"
    When I run doctor
    Then the "Files" category should have 1 issue
    And the doctor report should have 1 total issues

  Scenario: Orphan type directory is detected
    Given a vault is ready
    And an orphan object directory "ghost" exists
    When I run doctor
    Then the "Orphans" category should have 1 issue

  Scenario: Index out of sync is auto-fixed
    Given a vault is ready
    And the index is out of sync
    When I run doctor
    Then the doctor report should have 1 auto-fixed
    And the "Index" category should pass

  Scenario: Mixed issues across categories
    Given a vault is ready
    And a corrupted object file exists at "book/bad-file.md"
    And an orphan object directory "ghost" exists
    When I run doctor
    Then the "Files" category should have 1 issue
    And the "Orphans" category should have 1 issue
    And the doctor report should have 2 total issues
