Feature: Schema Cache
  The vault caches type schemas in memory to avoid repeated disk reads.

  Scenario: LoadType returns consistent results
    Given a vault is ready
    And a type schema file "project" with emoji "📁"
    When I load type "project"
    Then the loaded schema should have emoji "📁"

  Scenario: SaveType makes updated schema available via LoadType
    Given a vault is ready
    And a type schema file "project" with emoji "📁"
    When I load type "project"
    And I save type "project" with emoji "🗂️"
    And I load type "project"
    Then the loaded schema should have emoji "🗂️"

  Scenario: DeleteType makes type unavailable via LoadType
    Given a vault is ready
    And a type schema file "project" with emoji "📁"
    When I load type "project"
    And I delete type "project"
    And I load type "project"
    Then an error should occur

  Scenario: InvalidateSchemaCache forces reload from disk
    Given a vault is ready
    And a type schema file "project" with emoji "📁"
    When I load type "project"
    And the type schema file "project" is changed to emoji "🔄" on disk
    And I invalidate the schema cache
    And I load type "project"
    Then the loaded schema should have emoji "🔄"
