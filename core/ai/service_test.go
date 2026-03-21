package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// mockProvider is a test double for the Provider interface.
type mockProvider struct {
	completeFunc func(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)
	lastRequest  *CompletionRequest
}

func (m *mockProvider) Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	m.lastRequest = req
	if m.completeFunc != nil {
		return m.completeFunc(ctx, req)
	}
	return &CompletionResponse{}, nil
}

func newMockProvider(result string) *mockProvider {
	return &mockProvider{
		completeFunc: func(_ context.Context, _ *CompletionRequest) (*CompletionResponse, error) {
			return &CompletionResponse{
				Content:    result,
				JSONResult: json.RawMessage(result),
			}, nil
		},
	}
}

func newErrorProvider(err error) *mockProvider {
	return &mockProvider{
		completeFunc: func(_ context.Context, _ *CompletionRequest) (*CompletionResponse, error) {
			return nil, err
		},
	}
}

func TestNewAIService(t *testing.T) {
	p := &mockProvider{}
	svc := NewAIService(p, ServiceConfig{})
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
	if svc.provider != p {
		t.Error("provider not set correctly")
	}
}

func TestDescribe_Success(t *testing.T) {
	mock := newMockProvider(`{"description":"A comprehensive guide to Go concurrency patterns."}`)
	svc := NewAIService(mock, ServiceConfig{})

	obj := ObjectContext{
		Name:       "Go Concurrency",
		Properties: map[string]any{"author": "Rob Pike"},
		Body:       "This book covers goroutines, channels, and more.",
	}
	schema := SchemaContext{
		TypeName:        "book",
		TypeDescription: "Books in the library",
		PropertyDescriptions: map[string]string{
			"author": "The book's author",
		},
	}

	desc, err := svc.Describe(context.Background(), obj, schema)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if desc != "A comprehensive guide to Go concurrency patterns." {
		t.Errorf("unexpected description: %q", desc)
	}

	// Verify prompt contains expected context
	req := mock.lastRequest
	if !strings.Contains(req.UserPrompt, "Go Concurrency") {
		t.Error("prompt should contain object name")
	}
	if !strings.Contains(req.UserPrompt, "Rob Pike") {
		t.Error("prompt should contain property values")
	}
	if !strings.Contains(req.UserPrompt, "goroutines") {
		t.Error("prompt should contain body content")
	}
	if !strings.Contains(req.UserPrompt, "Books in the library") {
		t.Error("prompt should contain type description")
	}
	if !strings.Contains(req.UserPrompt, "The book's author") {
		t.Error("prompt should contain property descriptions")
	}
	if req.JSONSchema == nil {
		t.Error("expected JSON schema for structured output")
	}
}

func TestDescribe_CustomPrompt(t *testing.T) {
	mock := newMockProvider(`{"description":"test"}`)
	svc := NewAIService(mock, ServiceConfig{
		DescribePrompt: "Custom system prompt",
	})

	_, err := svc.Describe(context.Background(), ObjectContext{Name: "test"}, SchemaContext{TypeName: "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if mock.lastRequest.SystemPrompt != "Custom system prompt" {
		t.Errorf("expected custom system prompt, got %q", mock.lastRequest.SystemPrompt)
	}
}

func TestDescribe_DefaultPrompt(t *testing.T) {
	mock := newMockProvider(`{"description":"test"}`)
	svc := NewAIService(mock, ServiceConfig{})

	_, err := svc.Describe(context.Background(), ObjectContext{Name: "test"}, SchemaContext{TypeName: "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if mock.lastRequest.SystemPrompt != DefaultDescribePrompt {
		t.Error("expected default describe prompt")
	}
}

func TestDescribe_ProviderError(t *testing.T) {
	mock := newErrorProvider(fmt.Errorf("connection refused"))
	svc := NewAIService(mock, ServiceConfig{})

	_, err := svc.Describe(context.Background(), ObjectContext{Name: "test"}, SchemaContext{TypeName: "test"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "connection refused") {
		t.Errorf("expected provider error in message, got: %v", err)
	}
}

func TestDescribe_InvalidJSON(t *testing.T) {
	mock := newMockProvider(`not json`)
	svc := NewAIService(mock, ServiceConfig{})

	_, err := svc.Describe(context.Background(), ObjectContext{Name: "test"}, SchemaContext{TypeName: "test"})
	if err == nil {
		t.Fatal("expected error for invalid JSON response")
	}
	if !strings.Contains(err.Error(), "parse response") {
		t.Errorf("expected parse error, got: %v", err)
	}
}

func TestSuggestTags_Success(t *testing.T) {
	mock := newMockProvider(`{"tags":[{"name":"go","is_new":false,"reason":"About Go programming"},{"name":"concurrency","is_new":true,"reason":"Covers concurrent patterns"}]}`)
	svc := NewAIService(mock, ServiceConfig{})

	obj := ObjectContext{
		Name: "Go Concurrency Patterns",
		Body: "A deep dive into goroutines and channels.",
	}
	schema := SchemaContext{TypeName: "book"}
	tags := []TagInfo{
		{Name: "go", Description: "Go programming language"},
		{Name: "rust", Description: "Rust programming language"},
	}

	result, err := svc.SuggestTags(context.Background(), obj, schema, tags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(result.Tags))
	}
	if result.Tags[0].Name != "go" || result.Tags[0].IsNew {
		t.Errorf("first tag should be 'go' (existing), got: %+v", result.Tags[0])
	}
	if result.Tags[1].Name != "concurrency" || !result.Tags[1].IsNew {
		t.Errorf("second tag should be 'concurrency' (new), got: %+v", result.Tags[1])
	}

	// Verify prompt includes existing tags
	req := mock.lastRequest
	if !strings.Contains(req.UserPrompt, "go") {
		t.Error("prompt should contain existing tag names")
	}
	if !strings.Contains(req.UserPrompt, "Go programming language") {
		t.Error("prompt should contain tag descriptions")
	}
}

func TestSuggestTags_ProviderError(t *testing.T) {
	mock := newErrorProvider(fmt.Errorf("rate limited"))
	svc := NewAIService(mock, ServiceConfig{})

	_, err := svc.SuggestTags(context.Background(), ObjectContext{Name: "test"}, SchemaContext{TypeName: "test"}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "rate limited") {
		t.Errorf("expected provider error, got: %v", err)
	}
}

func TestExploreSchema_Success(t *testing.T) {
	mock := newMockProvider(`{"suggestions":[{"type":"add","property_name":"publisher","property_type":"string","reason":"3 of 5 objects mention a publisher","description":"The book publisher"},{"type":"remove","property_name":"isbn","reason":"No objects use this property"}]}`)
	svc := NewAIService(mock, ServiceConfig{})

	schema := SchemaContext{
		TypeName:      "book",
		PropertyNames: []string{"author", "isbn"},
		PropertyTypes: map[string]string{"author": "string", "isbn": "string"},
	}
	objects := []ObjectContext{
		{Name: "Book A", Properties: map[string]any{"author": "Author A"}, Body: "Published by Acme..."},
		{Name: "Book B", Properties: map[string]any{"author": "Author B"}, Body: "A great read."},
	}

	result, err := svc.ExploreSchema(context.Background(), schema, objects)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Suggestions) != 2 {
		t.Fatalf("expected 2 suggestions, got %d", len(result.Suggestions))
	}
	if result.Suggestions[0].Type != "add" || result.Suggestions[0].PropertyName != "publisher" {
		t.Errorf("first suggestion should be add publisher, got: %+v", result.Suggestions[0])
	}
	if result.Suggestions[1].Type != "remove" || result.Suggestions[1].PropertyName != "isbn" {
		t.Errorf("second suggestion should be remove isbn, got: %+v", result.Suggestions[1])
	}

	// Verify prompt includes schema and objects
	req := mock.lastRequest
	if !strings.Contains(req.UserPrompt, "author") {
		t.Error("prompt should contain property names")
	}
	if !strings.Contains(req.UserPrompt, "Book A") {
		t.Error("prompt should contain object names")
	}
}

func TestExploreSchema_ProviderError(t *testing.T) {
	mock := newErrorProvider(fmt.Errorf("timeout"))
	svc := NewAIService(mock, ServiceConfig{})

	_, err := svc.ExploreSchema(context.Background(), SchemaContext{TypeName: "test"}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDescribe_ModelPassedThrough(t *testing.T) {
	mock := newMockProvider(`{"description":"test"}`)
	svc := NewAIService(mock, ServiceConfig{Model: "custom-model"})

	_, _ = svc.Describe(context.Background(), ObjectContext{Name: "test"}, SchemaContext{TypeName: "test"})

	if mock.lastRequest.Model != "custom-model" {
		t.Errorf("expected model 'custom-model', got %q", mock.lastRequest.Model)
	}
}

func TestBuildDescribePrompt_MinimalObject(t *testing.T) {
	svc := NewAIService(&mockProvider{}, ServiceConfig{})
	prompt := svc.buildDescribePrompt(
		ObjectContext{Name: "Test"},
		SchemaContext{TypeName: "note"},
	)
	if !strings.Contains(prompt, "Test") {
		t.Error("prompt should contain object name")
	}
	if !strings.Contains(prompt, "note") {
		t.Error("prompt should contain type name")
	}
}

func TestBuildTagPrompt_IncludesExistingTags(t *testing.T) {
	svc := NewAIService(&mockProvider{}, ServiceConfig{})
	prompt := svc.buildTagPrompt(
		ObjectContext{Name: "Test"},
		SchemaContext{TypeName: "note"},
		[]TagInfo{
			{Name: "important", Description: "High priority items"},
			{Name: "draft"},
		},
	)
	if !strings.Contains(prompt, "important") {
		t.Error("prompt should contain tag name")
	}
	if !strings.Contains(prompt, "High priority items") {
		t.Error("prompt should contain tag description")
	}
	if !strings.Contains(prompt, "draft") {
		t.Error("prompt should contain tag without description")
	}
}

func TestBuildExplorePrompt_IncludesObjectsAndSchema(t *testing.T) {
	svc := NewAIService(&mockProvider{}, ServiceConfig{})
	prompt := svc.buildExplorePrompt(
		SchemaContext{
			TypeName:      "book",
			PropertyNames: []string{"author"},
			PropertyTypes: map[string]string{"author": "string"},
			PropertyDescriptions: map[string]string{"author": "The author"},
		},
		[]ObjectContext{
			{Name: "Book 1", Body: "Content here"},
		},
	)
	if !strings.Contains(prompt, "book") {
		t.Error("prompt should contain type name")
	}
	if !strings.Contains(prompt, "author") {
		t.Error("prompt should contain property name")
	}
	if !strings.Contains(prompt, "The author") {
		t.Error("prompt should contain property description")
	}
	if !strings.Contains(prompt, "Book 1") {
		t.Error("prompt should contain object name")
	}
	if !strings.Contains(prompt, "Content here") {
		t.Error("prompt should contain object body")
	}
}
