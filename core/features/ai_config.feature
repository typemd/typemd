Feature: AI configuration
  As a user
  I want to configure AI features in my vault
  So that I can enable AI-powered assistance

  Scenario: AI config with enabled flag
    Given a vault is initialized
    And a config file with content:
      """
      ai:
        enabled: true
      """
    When I open the vault
    Then no error should occur
    And AI should be enabled

  Scenario: AI config defaults when absent
    Given a vault is initialized
    When I open the vault
    Then no error should occur
    And AI should not be enabled

  Scenario: AI config with custom prompts
    Given a vault is initialized
    And a config file with content:
      """
      ai:
        prompts:
          describe: "Custom describe prompt"
          tag: "Custom tag prompt"
          explore: "Custom explore prompt"
      """
    When I open the vault
    Then no error should occur
    And the AI describe prompt should be "Custom describe prompt"
    And the AI tag prompt should be "Custom tag prompt"
    And the AI explore prompt should be "Custom explore prompt"

  Scenario: AI config with explore settings
    Given a vault is initialized
    And a config file with content:
      """
      ai:
        explore:
          sample_count: 20
          body_truncate: 1000
      """
    When I open the vault
    Then no error should occur
    And the AI explore sample count should be 20
    And the AI explore body truncate should be 1000

  Scenario: AI config with default explore settings
    Given a vault is initialized
    And a config file with content:
      """
      ai:
        enabled: true
      """
    When I open the vault
    Then no error should occur
    And the AI explore sample count should be 0
    And the AI explore body truncate should be 0
