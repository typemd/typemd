Feature: Skill instructions
  Skills are embedded SKILL.md files with YAML frontmatter (name, description)
  and markdown body. The instructions system provides a registry of available
  skills and can output them enriched with vault context.

  Scenario: List available skills
    When I list available skills
    Then I should get 3 skills
    And the skills should include "explore", "importer", "guide"
    And each skill should have a name and description

  Scenario: Get skill by name
    When I get the skill "explore"
    Then the skill name should be "explore"
    And the skill description should not be empty
    And the skill instructions should not be empty
    And the skill instructions should not contain "---"

  Scenario: Get raw skill content
    When I get the raw content of skill "explore"
    Then the raw content should start with "---"
    And the raw content should contain "name: explore"

  Scenario: Unknown skill returns error
    When I get the skill "nonexistent"
    Then an error should occur

  Scenario: Build skill context with types
    Given a vault is initialized
    And a type "book" exists with emoji "📚" and description "Track your reading"
    And the type "book" has a property "author" of type "string"
    And I open the vault
    When I build skill context from the vault
    Then the context should have types
    And the context types should include "book"
    And the type "book" in context should have property "author"

  Scenario: Build skill context with empty vault
    Given a vault is initialized
    And I open the vault
    When I build skill context from the vault
    Then the context should have types
    And the context types should include "tag"
    And the context types should include "page"

  Scenario: Vault override replaces embedded skill
    Given a vault is initialized
    And a vault override exists for skill "explore" with body "Custom instructions"
    And I open the vault
    When I get the skill "explore" with vault override
    Then the skill instructions should contain "Custom instructions"

  Scenario: Vault override without frontmatter
    Given a vault is initialized
    And a vault override exists for skill "explore" without frontmatter
    And I open the vault
    When I get the skill "explore" with vault override
    Then the skill name should be "explore"
    And the skill description should not be empty
