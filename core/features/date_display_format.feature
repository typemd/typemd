Feature: Date display format
  As a user
  I want to configure how dates and datetimes are displayed
  So that I can see dates in my preferred format

  Background:
    Given a vault is initialized

  # ── Config keys ──────────────────────────────────────────────────────────

  Scenario: Default date format is empty
    When I open the vault
    Then the config value "date_format" should be ""

  Scenario: Default datetime format is empty
    When I open the vault
    Then the config value "datetime_format" should be ""

  Scenario: Set date_format via config
    When I open the vault
    And I set config "date_format" to "DD/MM/YYYY"
    Then no error should occur
    And the config value "date_format" should be "DD/MM/YYYY"

  Scenario: Set datetime_format via config
    When I open the vault
    And I set config "datetime_format" to "YYYY/MM/DD HH:mm"
    Then no error should occur
    And the config value "datetime_format" should be "YYYY/MM/DD HH:mm"

  Scenario: Config keys include date_format and datetime_format
    When I open the vault
    Then the known config keys should include "date_format"
    And the known config keys should include "datetime_format"

  Scenario: Date format loaded from config file
    And a config file with content:
      """
      date_format: "DD.MM.YYYY"
      datetime_format: "DD.MM.YYYY HH:mm:ss"
      """
    When I open the vault
    Then the config value "date_format" should be "DD.MM.YYYY"
    And the config value "datetime_format" should be "DD.MM.YYYY HH:mm:ss"

  # ── Display formatting ──────────────────────────────────────────────────

  Scenario: Date property uses configured format in display
    And a config file with content:
      """
      date_format: "MM/DD/YYYY"
      """
    And a type schema "event" with a date property
    When I open the vault
    And a "event" object named "birthday" exists with property "date" set to "2026-03-28"
    And I build display properties for "birthday"
    Then the display property "date" should have formatted value "03/28/2026"

  Scenario: Datetime property uses configured format in display
    And a config file with content:
      """
      datetime_format: "DD/MM/YYYY HH:mm:ss"
      """
    And a type schema "event" with a datetime property
    When I open the vault
    And a "event" object named "meeting" exists with property "due_at" set to "2026-03-28T14:30:00"
    And I build display properties for "meeting"
    Then the display property "due_at" should have formatted value "28/03/2026 14:30:00"

  Scenario: Default datetime format uses space separator
    And a type schema "event" with a datetime property
    When I open the vault
    And a "event" object named "meeting" exists with property "due_at" set to "2026-03-28T14:30:00"
    And I build display properties for "meeting"
    Then the display property "due_at" should have formatted value "2026-03-28 14:30:00"

  Scenario: Unrecognized tokens pass through as literals
    And a config file with content:
      """
      date_format: "YYYY年MM月DD日"
      """
    And a type schema "event" with a date property
    When I open the vault
    And a "event" object named "birthday" exists with property "date" set to "2026-03-28"
    And I build display properties for "birthday"
    Then the display property "date" should have formatted value "2026年03月28日"
