Feature: Vault initialization and lifecycle
  A vault stores objects as Markdown files and manages a SQLite index.

  Scenario: Initialize a new vault
    When I initialize a new vault
    Then the vault directory structure should exist
    And the SQLite database should exist

  Scenario: Initialization creates .gitignore
    When I initialize a new vault
    Then the .gitignore should contain "index.db"

  Scenario: Double initialization fails
    Given a vault is initialized
    When I initialize the vault again
    Then an error should occur

  Scenario: Open and close vault
    Given a vault is initialized
    When I open the vault
    And I close the vault
    Then no error should occur

  Scenario: Opening an uninitialized vault fails
    When I open an uninitialized vault
    Then an error should occur

  Scenario: Auto-sync on open when index is empty
    Given a vault is initialized
    And an object file exists on disk at "book/test-book.md" with title "Test Book"
    When I open the vault
    Then the index should contain 1 object
