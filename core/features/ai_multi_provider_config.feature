Feature: AI multi-provider configuration
  As a user
  I want to configure multiple AI providers in my vault
  So that I can switch between local and cloud LLMs

  Scenario: AI config with providers and default
    Given a vault is initialized
    And a config file with content:
      """
      ai:
        enabled: true
        default: ollama
        providers:
          claude:
            type: cli
            model: claude-sonnet-4-20250514
          ollama:
            type: openai-compatible
            base_url: http://localhost:11434
            model: qwen3-coder:30b
      """
    When I open the vault
    Then no error should occur
    And the AI default provider should be "ollama"
    And the AI provider "claude" should have type "cli"
    And the AI provider "claude" should have model "claude-sonnet-4-20250514"
    And the AI provider "ollama" should have type "openai-compatible"
    And the AI provider "ollama" should have base_url "http://localhost:11434"
    And the AI provider "ollama" should have model "qwen3-coder:30b"

  Scenario: AI config with no providers section uses raw config
    Given a vault is initialized
    And a config file with content:
      """
      ai:
        enabled: false
      """
    When I open the vault
    Then no error should occur
    And the AI providers map should be empty

  Scenario: AI config with provider with api_key
    Given a vault is initialized
    And a config file with content:
      """
      ai:
        enabled: true
        default: my-server
        providers:
          my-server:
            type: openai-compatible
            base_url: http://192.168.1.100:8080
            model: llama3.2
            api_key: sk-test-key
      """
    When I open the vault
    Then no error should occur
    And the AI provider "my-server" should have api_key "sk-test-key"
