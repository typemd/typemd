Feature: AI availability detection
  As a user
  I want the vault to detect AI availability
  So that AI features are only shown when usable

  Scenario: AI enabled and claude binary available
    Given a vault is initialized
    And a config file with content:
      """
      ai:
        enabled: true
      """
    When I open the vault
    Then no error should occur
    And AI service should be available

  Scenario: AI not enabled
    Given a vault is initialized
    When I open the vault
    Then no error should occur
    And AI service should not be available

  Scenario: AI disabled explicitly
    Given a vault is initialized
    And a config file with content:
      """
      ai:
        enabled: false
      """
    When I open the vault
    Then no error should occur
    And AI service should not be available
