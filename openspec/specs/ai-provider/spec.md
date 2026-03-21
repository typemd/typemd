### Requirement: Provider interface defines AI completion contract

The `core/ai` package SHALL define a `Provider` interface with a `Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)` method. `CompletionRequest` SHALL include `SystemPrompt` (string), `UserPrompt` (string), `JSONSchema` (optional json.RawMessage for structured output), and `Model` (optional string override). `CompletionResponse` SHALL include `Content` (string) and `JSONResult` (json.RawMessage, populated when JSONSchema was provided).

#### Scenario: Provider interface is implementable

- **WHEN** a struct implements the `Complete` method with the correct signature
- **THEN** it SHALL satisfy the `Provider` interface

#### Scenario: CompletionRequest with JSON schema

- **WHEN** a `CompletionRequest` has a non-nil `JSONSchema`
- **THEN** the provider implementation SHALL request structured output conforming to that schema
- **AND** the response `JSONResult` SHALL contain the parsed JSON

#### Scenario: CompletionRequest without JSON schema

- **WHEN** a `CompletionRequest` has a nil `JSONSchema`
- **THEN** the provider SHALL return free-text in `Content`
- **AND** `JSONResult` SHALL be nil

### Requirement: ClaudeCLI provider invokes claude binary

The `ClaudeCLI` struct SHALL implement `Provider` by invoking the `claude` CLI binary as a subprocess. It SHALL use `-p` (print mode), `--output-format json`, and `--tools ""` (disable all tools) flags. When `JSONSchema` is provided, it SHALL pass `--json-schema <schema>`. When `Model` is set, it SHALL pass `--model <model>`.

#### Scenario: Successful completion without JSON schema

- **WHEN** `Complete` is called with a prompt and no JSON schema
- **THEN** ClaudeCLI SHALL invoke `claude -p --output-format json --tools "" --system-prompt <system> <user_prompt>`
- **AND** SHALL parse the JSON output to extract the `result` field into `Content`

#### Scenario: Successful completion with JSON schema

- **WHEN** `Complete` is called with a prompt and a JSON schema
- **THEN** ClaudeCLI SHALL additionally pass `--json-schema <schema>`
- **AND** SHALL parse the response into `JSONResult`

#### Scenario: Model override

- **WHEN** `Complete` is called with `Model` set to `"claude-haiku-4-5-20251001"`
- **THEN** ClaudeCLI SHALL pass `--model claude-haiku-4-5-20251001` to the CLI

#### Scenario: Claude binary not found

- **WHEN** the `claude` binary is not in PATH
- **THEN** `Complete` SHALL return an error indicating the binary was not found

#### Scenario: Claude CLI returns non-zero exit code

- **WHEN** the `claude` subprocess exits with a non-zero code
- **THEN** `Complete` SHALL return an error containing the stderr output

#### Scenario: Context cancellation

- **WHEN** the context is cancelled while the subprocess is running
- **THEN** `Complete` SHALL terminate the subprocess and return a context cancellation error

### Requirement: AIService provides domain-specific AI operations

The `AIService` struct SHALL hold a `Provider` and provide domain-specific methods: `Describe(ctx, obj, schema) (string, error)`, `SuggestTags(ctx, obj, schema, existingTags) (*TagSuggestion, error)`, and `ExploreSchema(ctx, schema, objects, config) (*SchemaExploration, error)`. Each method SHALL assemble the appropriate prompt context and JSON schema, call the provider, and parse the structured response.

#### Scenario: AIService is initialized with a Provider

- **WHEN** `NewAIService(provider)` is called
- **THEN** it SHALL return an `AIService` that delegates to the given provider

#### Scenario: AIService methods assemble prompts from domain objects

- **WHEN** `Describe` is called with an Object and TypeSchema
- **THEN** it SHALL construct a prompt containing the object's name, properties, body, and the type schema's property descriptions
- **AND** SHALL call `Provider.Complete` with that prompt

### Requirement: AI availability detection on Vault open

The Vault SHALL detect AI availability during `Open()`. If `ai.enabled` is `true` in config, the Vault SHALL check for the `claude` binary via `exec.LookPath`. If found, the Vault SHALL initialize an `AIService` and store it. The Vault SHALL expose an `AIService()` method that returns the service or nil.

#### Scenario: AI enabled and claude binary found

- **WHEN** config has `ai.enabled: true` and `claude` is in PATH
- **THEN** `vault.AIService()` SHALL return a non-nil `*ai.AIService`

#### Scenario: AI enabled but claude binary not found

- **WHEN** config has `ai.enabled: true` but `claude` is not in PATH
- **THEN** `vault.AIService()` SHALL return nil

#### Scenario: AI not enabled

- **WHEN** config has `ai.enabled: false` or the key is absent
- **THEN** `vault.AIService()` SHALL return nil
- **AND** no binary lookup SHALL be performed
