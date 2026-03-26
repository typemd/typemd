Feature: Object lock and unlock commands
  Users can lock and unlock objects via CLI.

  Scenario: Lock an unlocked object
    Given a vault with a "book" object named "lock-test"
    When I run tmd object lock on the object
    Then the output should contain "Locked"
    And the object should be locked

  Scenario: Lock an already locked object
    Given a vault with a "book" object named "already-locked"
    And the object is locked
    When I run tmd object lock on the object
    Then the output should contain "already locked"

  Scenario: Unlock a locked object
    Given a vault with a "book" object named "unlock-test"
    And the object is locked
    When I run tmd object unlock on the object
    Then the output should contain "Unlocked"
    And the object should not be locked

  Scenario: Unlock a non-locked object
    Given a vault with a "book" object named "not-locked"
    When I run tmd object unlock on the object
    Then the output should contain "not locked"
