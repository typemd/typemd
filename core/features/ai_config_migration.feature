Feature: AI config backward-compatible migration
  As a user with existing flat AI config
  I want my config to work without changes
  So that upgrading typemd doesn't break my setup

  Scenario: Old flat config with ai.model migrates to providers
    Given a vault is initialized
    And a config file with content:
      """
      ai:
        enabled: true
        model: claude-haiku-4-5-20251001
      """
    When I open the vault
    Then no error should occur
    And the AI default provider should be "claude"
    And the AI provider "claude" should have type "cli"
    And the AI provider "claude" should have model "claude-haiku-4-5-20251001"

  Scenario: Old flat config without ai.model migrates to providers
    Given a vault is initialized
    And a config file with content:
      """
      ai:
        enabled: true
      """
    When I open the vault
    Then no error should occur
    And the AI default provider should be "claude"
    And the AI provider "claude" should have type "cli"

  Scenario: New config with providers takes precedence
    Given a vault is initialized
    And a config file with content:
      """
      ai:
        enabled: true
        model: old-model
        default: ollama
        providers:
          ollama:
            type: openai-compatible
            base_url: http://localhost:11434
            model: qwen3-coder:30b
      """
    When I open the vault
    Then no error should occur
    And the AI default provider should be "ollama"
    And the AI provider "ollama" should have model "qwen3-coder:30b"
