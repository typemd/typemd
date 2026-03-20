Feature: Incremental sync
  The projector can sync specific files incrementally instead of walking all objects.

  Scenario: SyncFiles upserts a new object
    Given a vault is ready
    And a "book" type schema exists
    And I create an object "book" named "Go Programming"
    When I sync files for the created object
    Then the object should be searchable by "Go Programming"

  Scenario: SyncFiles removes a deleted object
    Given a vault is ready
    And a "book" type schema exists
    And I create an object "book" named "Temp Book"
    And a full sync is performed
    And the object file is deleted from disk
    When I sync files for the deleted object path
    Then the object should not be in the index

  Scenario: SyncFiles with empty paths falls back to full sync
    Given a vault is ready
    And a "book" type schema exists
    And I create an object "book" named "Full Sync Book"
    When I sync with empty paths
    Then the object should be searchable by "Full Sync Book"
