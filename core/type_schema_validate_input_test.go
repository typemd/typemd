package core

import (
	"testing"
)

func TestValidatePropertyValue_String(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"hello", false},
		{"", false},
		{"with spaces and 123", false},
	}
	for _, tt := range tests {
		err := ValidatePropertyValue("string", nil, tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidatePropertyValue(string, %q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
	}
}

func TestValidatePropertyValue_Number(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"42", false},
		{"-7", false},
		{"3.14", false},
		{"-0.5", false},
		{"0", false},
		{"abc", true},
		{"12abc", true},
		{"", true},
		{".", true},
		{"-", true},
		{"-.", true},
		{"123.", true},
		{"12.34.56", true},
	}
	for _, tt := range tests {
		err := ValidatePropertyValue("number", nil, tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidatePropertyValue(number, %q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
	}
}

func TestValidatePropertyValue_Date(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"2024-01-15", false},
		{"2024-12-31", false},
		{"not-a-date", true},
		{"2024/01/15", true},
		{"2024-13-01", true},
		{"2024-01-32", true},
		{"", true},
	}
	for _, tt := range tests {
		err := ValidatePropertyValue("date", nil, tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidatePropertyValue(date, %q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
	}
}

func TestValidatePropertyValue_Datetime(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"2024-01-15T10:30:00", false},
		{"2024-01-15T10:30:00Z", false},
		{"2024-01-15T10:30:00+08:00", false},
		{"2024-01-15", true},
		{"not-a-datetime", true},
		{"", true},
	}
	for _, tt := range tests {
		err := ValidatePropertyValue("datetime", nil, tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidatePropertyValue(datetime, %q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
	}
}

func TestValidatePropertyValue_URL(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"https://example.com", false},
		{"http://localhost:8080", false},
		{"https://example.com/path?q=1", false},
		{"ftp://example.com", true},
		{"not-a-url", true},
		{"", true},
	}
	for _, tt := range tests {
		err := ValidatePropertyValue("url", nil, tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidatePropertyValue(url, %q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
	}
}

func TestValidatePropertyValue_Select(t *testing.T) {
	options := []Option{
		{Value: "draft"},
		{Value: "published"},
		{Value: "archived"},
	}
	tests := []struct {
		input   string
		wantErr bool
	}{
		{"draft", false},
		{"published", false},
		{"archived", false},
		{"unknown", true},
		{"", true},
	}
	for _, tt := range tests {
		err := ValidatePropertyValue("select", options, tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidatePropertyValue(select, %q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
		}
	}
}

func TestValidatePropertyValue_UnknownType(t *testing.T) {
	err := ValidatePropertyValue("unknown_type", nil, "anything")
	if err != nil {
		t.Errorf("ValidatePropertyValue(unknown_type) should accept any input, got error: %v", err)
	}
}
