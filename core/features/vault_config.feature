Feature: Vault configuration
  As a user
  I want to configure vault-level settings
  So that CLI commands can use sensible defaults

  Scenario: Config file with valid content
    Given a vault is initialized
    And a config file with content:
      """
      cli:
        default_type: idea
      """
    When I open the vault
    Then no error should occur
    And the default type should be "idea"

  Scenario: Config file does not exist
    Given a vault is initialized
    When I open the vault
    Then no error should occur
    And the default type should be ""

  Scenario: Config file is empty
    Given a vault is initialized
    And a config file with content:
      """
      """
    When I open the vault
    Then no error should occur
    And the default type should be ""

  Scenario: Config file has invalid YAML
    Given a vault is initialized
    And a config file with content:
      """
      [invalid: yaml: content
      """
    When I open the vault
    Then an error should occur

  Scenario: Unknown keys are ignored
    Given a vault is initialized
    And a config file with content:
      """
      unknown_key: some_value
      cli:
        default_type: note
      """
    When I open the vault
    Then no error should occur
    And the default type should be "note"
