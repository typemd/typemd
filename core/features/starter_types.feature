Feature: Starter type templates
  Starter types are embedded YAML type schemas offered during vault initialization.
  Users can select which starter types to install; selected types are written as
  regular .typemd/types/*.yaml files.

  Scenario: List available starter types
    When I list available starter types
    Then I should get 3 starter types
    And the starter types should include "idea", "note", "book"
    And each starter type should have a name, emoji, and description

  Scenario: Starter type YAML is valid
    When I list available starter types
    Then each starter type YAML should parse as a valid TypeSchema

  Scenario: Write selected starter types to vault
    Given a vault is initialized
    When I write starter types "idea,book" to the vault
    Then the file ".typemd/types/idea.yaml" should exist
    And the file ".typemd/types/book.yaml" should exist
    And the file ".typemd/types/note.yaml" should not exist

  Scenario: Write all starter types to vault
    Given a vault is initialized
    When I write all starter types to the vault
    Then the file ".typemd/types/idea.yaml" should exist
    And the file ".typemd/types/note.yaml" should exist
    And the file ".typemd/types/book.yaml" should exist

  Scenario: Write no starter types
    Given a vault is initialized
    When I write starter types "" to the vault
    Then the file ".typemd/types/idea.yaml" should not exist
    And the file ".typemd/types/note.yaml" should not exist
    And the file ".typemd/types/book.yaml" should not exist

  Scenario: Written starter types are loadable
    Given a vault is initialized
    When I write all starter types to the vault
    And I open the vault
    Then I should be able to load type "book"
    And I should be able to load type "idea"
    And I should be able to load type "note"
