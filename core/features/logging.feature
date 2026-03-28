Feature: Structured logging initialization
  typemd provides centralized logging via Go's log/slog package.

  Scenario: Initialize logging with debug level
    When I initialize logging at debug level
    And I log a debug message "sync started"
    Then the log output should contain "sync started"
    And the log output should contain "DEBUG"

  Scenario: Initialize logging with warn level
    When I initialize logging at warn level
    And I log a debug message "sync started"
    Then the log output should be empty

  Scenario: Warn messages pass warn level filter
    When I initialize logging at warn level
    And I log a warn message "object not found"
    Then the log output should contain "object not found"
    And the log output should contain "WARN"

  Scenario: Log output is valid JSON
    When I initialize logging at debug level
    And I log a debug message "test" with attribute "component" value "projector"
    Then the log output should be valid JSON
    And the log JSON should have field "msg" with value "test"
    And the log JSON should have field "component" with value "projector"
    And the log JSON should have field "level"
    And the log JSON should have field "time"
