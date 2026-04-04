Feature: Vault directory migration
  When opening a vault with the old directory layout (.typemd/types/ and
  .typemd/properties.yaml), typemd automatically migrates to the new
  root-level layout (types/ and properties/).

  Background:
    Given a vault is initialized

  Scenario: Migrate types from old to new location
    Given types exist at the old location ".typemd/types/"
    When I open the vault
    Then types should exist at "types/"
    And types should not exist at ".typemd/types/"

  Scenario: Migrate properties from old to new location
    Given a properties file exists at the old location ".typemd/properties.yaml"
    When I open the vault
    Then a properties file should exist at "properties/properties.yaml"
    And a properties file should not exist at ".typemd/properties.yaml"

  Scenario: No migration when already at new location
    When I open the vault
    Then no error should occur

  Scenario: Skip migration when old paths do not exist
    When I open the vault
    Then no error should occur

  Scenario: Conflict when both types directories exist
    Given types exist at the old location ".typemd/types/"
    And types exist at the new location "types/"
    When I open the vault
    Then an error should occur
    And the last error should mention "conflict"

  Scenario: Conflict when both properties paths exist
    Given a properties file exists at the old location ".typemd/properties.yaml"
    And a properties file exists at the new location "properties/properties.yaml"
    When I open the vault
    Then an error should occur
    And the last error should mention "conflict"

  Scenario: Conflict in one path prevents all migration
    Given types exist at the old location ".typemd/types/"
    And types exist at the new location "types/"
    And a properties file exists at the old location ".typemd/properties.yaml"
    When I open the vault
    Then an error should occur
    And a properties file should still exist at ".typemd/properties.yaml"
