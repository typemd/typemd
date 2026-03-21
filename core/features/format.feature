Feature: Object and schema formatting
  The format command normalizes frontmatter property ordering and YAML style
  without modifying body content or updated_at timestamps.

  Background:
    Given a vault is initialized

  Scenario: Format object with out-of-order properties
    Given a type "book" with properties:
      | name   | type   |
      | author | string |
      | rating | number |
    And an object "book/test" with raw frontmatter:
      """
      ---
      author: Alice
      name: Test Book
      rating: 5
      created_at: "2025-01-01T00:00:00Z"
      updated_at: "2025-01-01T00:00:00Z"
      ---
      """
    When I format all objects
    Then 1 object should be formatted
    And the object "book/test" frontmatter should have "name" before "author"

  Scenario: Already-formatted objects are not rewritten
    Given a type "book" with properties:
      | name   | type   |
      | author | string |
    And an object "book/test" created through the vault with name "Test Book"
    When I format all objects
    Then 0 objects should be formatted

  Scenario: Body content is preserved after formatting
    Given a type "book" with properties:
      | name   | type   |
      | author | string |
    And an object "book/test" with raw frontmatter:
      """
      ---
      author: Alice
      name: Test Book
      created_at: "2025-01-01T00:00:00Z"
      updated_at: "2025-01-01T00:00:00Z"
      ---

      This is the body content.
      """
    When I format all objects
    Then the object "book/test" body should be "This is the body content."

  Scenario: updated_at is not modified during formatting
    Given a type "book" with properties:
      | name   | type   |
      | author | string |
    And an object "book/test" with raw frontmatter:
      """
      ---
      author: Alice
      name: Test Book
      created_at: "2025-01-01T00:00:00Z"
      updated_at: "2025-06-15T12:00:00Z"
      ---
      """
    When I format all objects
    Then the object "book/test" property "updated_at" should be "2025-06-15T12:00:00Z"

  Scenario: Format only objects of a specific type
    Given a type "book" with properties:
      | name   | type   |
      | author | string |
    And a type "note" with properties:
      | name   | type   |
      | topic  | string |
    And an object "book/test" with raw frontmatter:
      """
      ---
      author: Alice
      name: Test Book
      created_at: "2025-01-01T00:00:00Z"
      updated_at: "2025-01-01T00:00:00Z"
      ---
      """
    And an object "note/test" with raw frontmatter:
      """
      ---
      topic: Go
      name: Test Note
      created_at: "2025-01-01T00:00:00Z"
      updated_at: "2025-01-01T00:00:00Z"
      ---
      """
    When I format objects of type "book"
    Then 1 object should be formatted

  Scenario: Dry-run mode lists files without writing
    Given a type "book" with properties:
      | name   | type   |
      | author | string |
    And an object "book/test" with raw frontmatter:
      """
      ---
      author: Alice
      name: Test Book
      created_at: "2025-01-01T00:00:00Z"
      updated_at: "2025-01-01T00:00:00Z"
      ---
      """
    When I format all objects in dry-run mode
    Then 1 object should be formatted
    And the object "book/test" file should still have "author" before "name"

  Scenario: Format schema with non-canonical YAML
    Given a type "book" with raw schema:
      """
      name: book
      emoji: "\U0001F4DA"
      plural: books
      properties:
        - name: author
          type: string
        - name: rating
          type: number
      """
    When I format all schemas
    Then 1 schema should be formatted

  Scenario: Already-formatted schema is not rewritten
    Given a type "book" with properties:
      | name   | type   |
      | author | string |
    When I format all schemas
    Then 0 schemas should be formatted

  Scenario: Format invalid type name returns error
    When I format objects of type "nonexistent"
    Then an error should occur
