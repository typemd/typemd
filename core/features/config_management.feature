Feature: Config management
  As a user
  I want to get and set vault config values via dot-notation keys
  So that I can manage vault settings without editing files directly

  Scenario: Set a known config key
    Given a vault is initialized
    When I open the vault
    And I set config "cli.default_type" to "idea"
    Then no error should occur
    And the config value "cli.default_type" should be "idea"

  Scenario: Set updates existing config value
    Given a vault is initialized
    And a config file with content:
      """
      cli:
        default_type: idea
      """
    When I open the vault
    And I set config "cli.default_type" to "note"
    Then no error should occur
    And the config value "cli.default_type" should be "note"

  Scenario: Set creates config file if missing
    Given a vault is initialized
    When I open the vault
    And I set config "cli.default_type" to "idea"
    Then no error should occur
    And the config file should exist

  Scenario: Set with unknown key returns error
    Given a vault is initialized
    When I open the vault
    And I set config "unknown.key" to "value"
    Then an error should occur

  Scenario: Get a set config value
    Given a vault is initialized
    And a config file with content:
      """
      cli:
        default_type: idea
      """
    When I open the vault
    Then the config value "cli.default_type" should be "idea"

  Scenario: Get an unset known key returns empty
    Given a vault is initialized
    When I open the vault
    Then the config value "cli.default_type" should be ""

  Scenario: Get an unknown key is not known
    Given a vault is initialized
    When I open the vault
    Then the config key "nonexistent" should not be known

  Scenario: List known config keys
    Given a vault is initialized
    When I open the vault
    Then the known config keys should include "cli.default_type"
