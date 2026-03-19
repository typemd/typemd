Feature: Vault statistics
  Users can view aggregate statistics about their vault, including per-type
  object counts and single-type property aggregations.

  Scenario: Vault-wide stats with multiple types
    Given a vault is ready
    And a "book" object named "book1" exists
    And a "book" object named "book2" exists
    And a "book" object named "book3" exists
    And a "person" object named "alice" exists
    And a "person" object named "bob" exists
    When I request vault stats
    Then the vault stats should show 2 types
    And the vault stats total should be 5
    And the vault stats for type "book" should show count 3
    And the vault stats for type "person" should show count 2

  Scenario: Vault-wide stats on empty vault
    Given a vault is ready
    When I request vault stats
    Then the vault stats should show 0 types
    And the vault stats total should be 0

  Scenario: Vault-wide stats includes built-in types
    Given a vault is ready
    And a "tag" object named "fiction" exists
    When I request vault stats
    Then the vault stats total should be 1
    And the vault stats for type "tag" should show count 1
    And the vault stats for type "tag" should show emoji "🏷️"
