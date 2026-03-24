## ADDED Requirements

### Requirement: OpenAICompatible provider implements Provider interface via HTTP

The `OpenAICompatible` struct in `core/ai` SHALL implement the `Provider` interface by sending HTTP POST requests to `{BaseURL}/v1/chat/completions`. It SHALL accept `BaseURL` (string, required), `Model` (string, required), `APIKey` (string, optional), and `HTTPClient` (optional `*http.Client` for testing). If no `HTTPClient` is provided, a default client with 60-second timeout SHALL be used.

#### Scenario: Successful completion with structured JSON output

- **WHEN** `Complete` is called with a `CompletionRequest` that has a non-nil `JSONSchema`
- **THEN** the provider SHALL send a POST to `{BaseURL}/v1/chat/completions` with `response_format: { type: "json_schema", json_schema: { name: "response", strict: true, schema: <JSONSchema> } }`
- **AND** the request body SHALL include `model`, `messages` (system + user), and `response_format`
- **AND** the response `JSONResult` SHALL contain the parsed JSON from `choices[0].message.content`

#### Scenario: Successful completion without JSON schema

- **WHEN** `Complete` is called with a `CompletionRequest` that has a nil `JSONSchema`
- **THEN** the provider SHALL send a POST without `response_format`
- **AND** the response `Content` SHALL contain `choices[0].message.content`
- **AND** `JSONResult` SHALL be nil

#### Scenario: Model override in request

- **WHEN** `Complete` is called with `CompletionRequest.Model` set to `"llama3.2"`
- **THEN** the request body `model` field SHALL be `"llama3.2"` (overriding the struct's default `Model`)

#### Scenario: Model from struct default

- **WHEN** `Complete` is called with `CompletionRequest.Model` empty
- **THEN** the request body `model` field SHALL use the struct's `Model` value

#### Scenario: API key authentication

- **WHEN** `APIKey` is set on the `OpenAICompatible` struct
- **THEN** all HTTP requests SHALL include an `Authorization: Bearer <APIKey>` header

#### Scenario: No API key

- **WHEN** `APIKey` is empty
- **THEN** HTTP requests SHALL NOT include an `Authorization` header

### Requirement: OpenAICompatible handles connection errors with helpful messages

The `OpenAICompatible` provider SHALL wrap common HTTP errors with user-friendly messages that help diagnose the problem.

#### Scenario: Server not running (connection refused)

- **WHEN** `Complete` is called and the HTTP request fails with a connection refused error
- **THEN** the error message SHALL contain the base URL
- **AND** the error message SHALL contain "connection refused"

#### Scenario: HTTP non-200 response

- **WHEN** the server returns a non-200 HTTP status code
- **THEN** the error message SHALL contain the HTTP status code
- **AND** the error message SHALL contain the response body (truncated to 500 bytes if longer)

#### Scenario: Invalid JSON response

- **WHEN** the server returns HTTP 200 but the response body is not valid JSON
- **THEN** the error message SHALL indicate a JSON parse failure

#### Scenario: Context cancellation

- **WHEN** the context is cancelled during the HTTP request
- **THEN** `Complete` SHALL return a context cancellation error
