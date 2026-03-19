Feature: View configuration
  ViewConfig defines how objects of a type are filtered, sorted, grouped, and displayed.

  # ── ViewConfig creation ─────────────────────────────────────

  Scenario: Create a ViewConfig with all fields
    Given a view config "by-rating" with layout "list"
    And the view has filter property "status" operator "is" value "reading"
    And the view has sort property "rating" direction "desc"
    And the view has group_by property "genre"
    Then the view name should be "by-rating"
    And the view layout should be "list"
    And the view should have 1 filter rule
    And the view should have 1 sort rule
    And the view should have 1 group rule
    And the view group_by property 1 should be "genre"

  Scenario: Create a ViewConfig with multiple group rules
    Given a view config "multi-group" with layout "list"
    And the view has group_by property "genre"
    And the view has group_by property "status"
    Then the view should have 2 group rules
    And the view group_by property 1 should be "genre"
    And the view group_by property 2 should be "status"

  Scenario: Create a minimal ViewConfig
    Given a view config "default" with layout "list"
    Then the view name should be "default"
    And the view should have 0 filter rules
    And the view should have 0 sort rules
    And the view should have 0 group rules

  # ── YAML serialization ──────────────────────────────────────

  Scenario: Serialize ViewConfig with group rules to YAML
    Given a view config "by-rating" with layout "list"
    And the view has sort property "rating" direction "desc"
    And the view has group_by property "genre"
    When I serialize the view config to YAML
    Then the view YAML should contain "name: by-rating"
    And the view YAML should contain "layout: list"
    And the view YAML should contain "property: rating"
    And the view YAML should contain "group_by:"
    And the view YAML should contain "property: genre"

  Scenario: Serialize ViewConfig without group rules omits group_by
    Given a view config "simple" with layout "list"
    And the view has sort property "name" direction "asc"
    When I serialize the view config to YAML
    Then the view YAML should not contain "group_by:"

  Scenario: Deserialize ViewConfig with new array group_by format
    Given view YAML content:
      """
      name: reading-now
      layout: list
      filter:
        - property: status
          operator: is
          value: reading
      sort:
        - property: name
          direction: asc
      group_by:
        - property: genre
      """
    When I deserialize the view YAML
    Then the deserialized view name should be "reading-now"
    And the deserialized view should have 1 filter rule
    And the deserialized view should have 1 sort rule
    And the deserialized view should have 1 group rule
    And the deserialized view group_by property 1 should be "genre"

  Scenario: Deserialize ViewConfig with legacy string group_by
    Given view YAML content:
      """
      name: legacy-view
      layout: list
      group_by: genre
      """
    When I deserialize the view YAML
    Then the deserialized view name should be "legacy-view"
    And the deserialized view should have 1 group rule
    And the deserialized view group_by property 1 should be "genre"

  Scenario: Deserialize ViewConfig with multiple group rules
    Given view YAML content:
      """
      name: multi-group
      layout: list
      group_by:
        - property: genre
        - property: status
      """
    When I deserialize the view YAML
    Then the deserialized view should have 2 group rules
    And the deserialized view group_by property 1 should be "genre"
    And the deserialized view group_by property 2 should be "status"

  Scenario: Deserialize ViewConfig without group_by
    Given view YAML content:
      """
      name: no-group
      layout: list
      """
    When I deserialize the view YAML
    Then the deserialized view should have 0 group rules
