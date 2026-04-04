Feature: Config key info
  As a TUI developer
  I want to access config key metadata
  So that I can display descriptions and defaults in the config settings page

  Scenario: All config keys have descriptions
    Given a vault is initialized
    When I open the vault
    Then every config key should have a description

  Scenario: Config key info includes current and default values
    Given a vault is initialized
    And a config file with content:
      """
      cli:
        default_type: idea
      """
    When I open the vault
    Then the config key info for "cli.default_type" should have value "idea"
    And the config key info for "tui.debounce_ms" should have a non-empty default
