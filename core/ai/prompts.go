package ai

// Default system prompts for AI operations.
const (
	DefaultDescribePrompt = `You are a knowledge management assistant. Generate a concise, informative description for the given object based on its name, properties, body content, and the type schema context. The description should capture the essence of the object in 1-2 sentences. Do not repeat the object name verbatim. Focus on what makes this object noteworthy or useful.`

	DefaultTagPrompt = `You are a knowledge management assistant. Suggest relevant tags for the given object based on its name, properties, body content, and the type schema context. Prefer existing tags over creating new ones. Only suggest new tags when no existing tag adequately covers the concept. Each tag name should be lowercase, using hyphens for multi-word tags. Provide a brief reason for each suggestion.`

	DefaultExplorePrompt = `You are a knowledge management assistant analyzing a collection of objects to suggest schema improvements. Examine the provided objects and current type schema, then suggest property additions, modifications, or removals. Focus on patterns you observe: properties that many objects share but aren't in the schema, properties in the schema that no objects use, or properties whose type doesn't match the actual data. Be conservative — only suggest changes with clear evidence from the data.`
)
