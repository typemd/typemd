Feature: AI provider resolution
  As a user
  I want the vault to resolve my configured AI provider
  So that the correct AI backend is used

  Scenario: OpenAI-compatible provider configured
    Given a vault is initialized
    And a config file with content:
      """
      ai:
        enabled: true
        default: ollama
        providers:
          ollama:
            type: openai-compatible
            base_url: http://localhost:11434
            model: qwen3-coder:30b
      """
    When I open the vault
    Then no error should occur
    And AI service should be available

  Scenario: Enabled without providers yields no AI service
    Given a vault is initialized
    And a config file with content:
      """
      ai:
        enabled: true
      """
    When I open the vault
    Then no error should occur
    And AI service should not be available

  Scenario: Default points to undefined provider
    Given a vault is initialized
    And a config file with content:
      """
      ai:
        enabled: true
        default: missing
        providers:
          claude:
            type: cli
      """
    When I open the vault
    Then no error should occur
    And AI service should not be available

  Scenario: Unknown provider type
    Given a vault is initialized
    And a config file with content:
      """
      ai:
        enabled: true
        default: custom
        providers:
          custom:
            type: unknown-type
            base_url: http://example.com
      """
    When I open the vault
    Then no error should occur
    And AI service should not be available

  Scenario: OpenAI-compatible provider missing base_url
    Given a vault is initialized
    And a config file with content:
      """
      ai:
        enabled: true
        default: broken
        providers:
          broken:
            type: openai-compatible
            model: llama3.2
      """
    When I open the vault
    Then no error should occur
    And AI service should not be available
