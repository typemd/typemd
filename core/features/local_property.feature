Feature: Local property identification
  Properties in object frontmatter that are not defined in the type schema
  should be identified as "local" in display properties.

  Scenario: Schema-defined property is not local
    Given a vault is ready
    And a type schema "item" with a "status" string property
    And a "item" object named "test-item" exists with property "status" set to "active"
    When I build display properties for "test-item"
    Then the display property "status" should not be local

  Scenario: System property is not local
    Given a vault is ready
    And a type schema "item" with a "status" string property
    And a "item" object named "test-item" exists
    When I build display properties for "test-item"
    Then the display property "created_at" should not be local

  Scenario: Extra property is local
    Given a vault is ready
    And a type schema "item" with a "status" string property
    And a "item" object named "test-item" exists with extra property "custom_field" set to "hello"
    When I build display properties for "test-item"
    Then the display property "custom_field" should be local

  Scenario: Object with no local properties
    Given a vault is ready
    And a type schema "item" with a "status" string property
    And a "item" object named "test-item" exists with property "status" set to "done"
    When I build display properties for "test-item"
    Then no display property should be local

  Scenario: Object with no schema has all non-system properties as local
    Given a vault is ready
    And a type schema "item" with a "status" string property
    And a "item" object named "no-schema-item" exists with extra property "foo" set to "bar"
    And the type schema "item" is removed
    When I build display properties for "no-schema-item"
    Then the display property "foo" should be local
