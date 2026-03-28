## ADDED Requirements

### Requirement: tmd serve starts HTTP server
The `tmd serve` command SHALL start a local HTTP server that serves both the REST API and the embedded frontend.

#### Scenario: Default port
- **WHEN** user runs `tmd serve` without flags
- **THEN** the server SHALL listen on port 3000

#### Scenario: Custom port
- **WHEN** user runs `tmd serve -p 8080`
- **THEN** the server SHALL listen on port 8080

#### Scenario: Invalid vault
- **WHEN** user runs `tmd serve` in a directory without `.typemd/`
- **THEN** the command SHALL return an error indicating the vault is not initialized

### Requirement: REST API lists types
The server SHALL expose `GET /api/types` returning all types with metadata.

#### Scenario: List types
- **WHEN** a GET request is made to `/api/types`
- **THEN** the response SHALL be a JSON array where each element contains `name`, `plural`, `emoji`, `color`, and `count` fields

#### Scenario: Type count reflects object count
- **WHEN** a vault has 4 objects of type "book"
- **THEN** the type entry for "book" SHALL have `count: 4`

### Requirement: REST API returns type schema
The server SHALL expose `GET /api/types/{name}` returning the full type schema with property definitions.

#### Scenario: Get type schema
- **WHEN** a GET request is made to `/api/types/book`
- **THEN** the response SHALL contain `name`, `plural`, `emoji`, `color`, and a `properties` array with each property's `name`, `type`, `emoji`, `pin`, `target`, and `options`

#### Scenario: Unknown type
- **WHEN** a GET request is made to `/api/types/nonexistent`
- **THEN** the response SHALL be HTTP 404 with an error message

### Requirement: REST API lists objects
The server SHALL expose `GET /api/objects` returning objects, optionally filtered by type.

#### Scenario: List all objects
- **WHEN** a GET request is made to `/api/objects`
- **THEN** the response SHALL be a JSON array of all objects with `id`, `type`, `name`, and `locked` fields, sorted by name ascending

#### Scenario: Filter by type
- **WHEN** a GET request is made to `/api/objects?type=book`
- **THEN** the response SHALL contain only objects of type "book"

### Requirement: REST API returns single object
The server SHALL expose `GET /api/objects/{type}/{slug}` returning the full object with properties and body.

#### Scenario: Get object
- **WHEN** a GET request is made to `/api/objects/book/clean-code-01xxx`
- **THEN** the response SHALL contain `id`, `type`, `name`, `properties` (serialized map), `body`, and `locked`

#### Scenario: Unknown object
- **WHEN** a GET request is made to `/api/objects/book/nonexistent`
- **THEN** the response SHALL be HTTP 404

### Requirement: REST API updates object body
The server SHALL expose `PUT /api/objects/{type}/{slug}` to update the object body.

#### Scenario: Update body
- **WHEN** a PUT request with `{"body": "new content"}` is made to `/api/objects/book/clean-code-01xxx`
- **THEN** the object's body SHALL be updated and the response SHALL contain the updated object

### Requirement: REST API creates objects
The server SHALL expose `POST /api/objects` to create a new object.

#### Scenario: Create object
- **WHEN** a POST request with `{"type": "book", "name": "New Book"}` is made to `/api/objects`
- **THEN** a new object SHALL be created and the response SHALL be HTTP 201 with the new object detail

#### Scenario: Create with template
- **WHEN** a POST request with `{"type": "book", "name": "New Book", "template": "review"}` is made to `/api/objects`
- **THEN** the object SHALL be created with the template applied

#### Scenario: Missing required fields
- **WHEN** a POST request with `{"type": ""}` is made to `/api/objects`
- **THEN** the response SHALL be HTTP 400 with an error message

### Requirement: REST API returns display properties
The server SHALL expose `GET /api/properties/{type}/{slug}` returning display-ready properties.

#### Scenario: Get display properties
- **WHEN** a GET request is made to `/api/properties/book/clean-code-01xxx`
- **THEN** the response SHALL be a JSON array where each element contains `key`, `value`, `display`, `type`, `emoji`, `pin`, `isRelation`, `isReverse`, and `isBacklink`

### Requirement: REST API updates single property
The server SHALL expose `PUT /api/properties/{type}/{slug}/{key}` to update a single property value.

#### Scenario: Update string property
- **WHEN** a PUT request with `{"value": "new title"}` is made to `/api/properties/book/clean-code-01xxx/title`
- **THEN** the property SHALL be updated and saved

#### Scenario: Invalid property value
- **WHEN** a PUT request with an invalid value is made (e.g., non-numeric string for a number property)
- **THEN** the response SHALL be HTTP 400 with a validation error

### Requirement: REST API lists templates
The server SHALL expose `GET /api/templates/{type}` returning available template names.

#### Scenario: List templates
- **WHEN** a GET request is made to `/api/templates/book`
- **THEN** the response SHALL be a JSON array of template name strings

### Requirement: Embedded frontend serving
The server SHALL serve the embedded React frontend for non-API routes with SPA fallback.

#### Scenario: Serve index
- **WHEN** a GET request is made to `/`
- **THEN** the server SHALL serve `index.html` from the embedded frontend

#### Scenario: SPA fallback
- **WHEN** a GET request is made to a path that does not match an API route or a static file
- **THEN** the server SHALL serve `index.html` for client-side routing

#### Scenario: Static assets
- **WHEN** a GET request is made to `/assets/index-xxx.js`
- **THEN** the server SHALL serve the exact file from the embedded frontend

### Requirement: Time values serialized as RFC3339
The server SHALL serialize `time.Time` values as RFC3339 strings in JSON responses.

#### Scenario: Date property
- **WHEN** an object has a `created_at` property with a time value
- **THEN** the JSON response SHALL contain the value as an RFC3339 string (e.g., `"2026-03-24T23:04:43+08:00"`)
