Feature: tmd fix wikilinks command
  The fix wikilinks command expands shorthand wiki-links to full IDs.

  Background:
    Given a vault is ready

  Scenario: Shorthand wiki-links already expanded on vault open
    Given a "note" object with a shorthand wiki-link exists
    When I run fix wikilinks
    Then the command should succeed
    And the output should contain "No changes needed"

  Scenario: No changes needed
    When I run fix wikilinks
    Then the command should succeed
    And the output should contain "No changes needed"

  Scenario: Ambiguous wiki-links are reported
    Given an object with an ambiguous shorthand wiki-link exists
    When I run fix wikilinks
    Then the command should succeed
    And the output should contain "ambiguous"
