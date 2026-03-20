Feature: Vault initialization
  As a user
  I want to initialize a new vault via CLI
  So that I can start managing objects

  Scenario: Initialize a new vault
    Given an empty directory
    When I run init with no-starters
    Then the command should succeed
    And the output should contain "Initialized vault at"
    And the .typemd directory should exist
    And the config should have default_type "page"

  Scenario: Initialize an already-initialized vault
    Given a vault is already initialized
    When I run init
    Then the command should fail
