Feature: Object locking
  Users can lock individual objects to prevent accidental editing.

  Scenario: IsLocked returns true for locked object
    Given a vault is ready
    When I create a "book" object named "locked-book"
    And I lock the object
    Then the object should be locked

  Scenario: IsLocked returns false for unlocked object
    Given a vault is ready
    When I create a "book" object named "unlocked-book"
    Then the object should not be locked

  Scenario: SaveObject on locked object returns error
    Given a vault is ready
    And a "book" object named "lock-save-test" exists
    And I lock the object
    When I save the object
    Then an error should occur

  Scenario: SetProperty on locked object returns error
    Given a vault is ready
    And a "book" object named "lock-setprop-test" exists
    And I lock the object
    When I set property "title" to "New Title" on the object
    Then an error should occur

  Scenario: SetPropertyMultiple on locked object returns error
    Given a vault is ready
    And a "book" object named "lock-setmultiple-test" exists
    And I lock the object
    When I set multiple properties on the object:
      | key   | value     |
      | title | New Title |
    Then an error should occur

  Scenario: Unlocked object can be modified normally
    Given a vault is ready
    And a "book" object named "normal-book" exists
    When I set property "title" to "New Title" on the object
    Then no error should occur

  Scenario: SetLocked locks an object
    Given a vault is ready
    And a "book" object named "setlocked-book" exists
    When I lock the object
    Then the object should be locked

  Scenario: SetLocked unlocks an object
    Given a vault is ready
    And a "book" object named "unlock-book" exists
    And I lock the object
    When I unlock the object
    Then the object should not be locked

  Scenario: LinkObjects on locked source returns error
    Given a vault is ready
    And a vault is ready with relation schemas
    And a "book" object named "link-lock-book" exists
    And a "person" object named "link-lock-person" exists
    And I lock the source object
    When I link "link-lock-book" to "link-lock-person" via "author"
    Then an error should occur

  Scenario: UnlinkObjects on locked source returns error
    Given a vault is ready
    And a vault is ready with relation schemas
    And a "book" object named "unlink-lock-book" exists
    And a "person" object named "unlink-lock-person" exists
    And I link "unlink-lock-book" to "unlink-lock-person" via "author"
    And I lock the source object
    When I unlink "unlink-lock-book" from "unlink-lock-person" via "author" without both flag
    Then an error should occur

