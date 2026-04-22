package main

import (
	"strings"
	"testing"
)

func TestParseYAML_Minimal(t *testing.T) {
	input := `
name: Alex
title: Engineer
theme: classic-navy
contact:
  email: alex@example.com
sections:
  - type: text
    title: Summary
    content: Hello
`
	data, err := ParseYAML(input)
	if err != nil {
		t.Fatalf("ParseYAML: %v", err)
	}
	if data.Name != "Alex" {
		t.Errorf("Name = %q, want %q", data.Name, "Alex")
	}
	if len(data.Sections) != 1 || data.Sections[0].Type != "text" {
		t.Errorf("Sections = %+v, want one text section", data.Sections)
	}
}

func TestParseYAML_InvalidReturnsError(t *testing.T) {
	_, err := ParseYAML("!!invalid: [unclosed")
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
	if !strings.Contains(err.Error(), "yaml") {
		t.Errorf("error should mention yaml, got: %v", err)
	}
}
