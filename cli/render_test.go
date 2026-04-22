package main

import (
	"strings"
	"testing"
)

const testThemesDir = ".."

func TestRenderTemplate_ClassicNavy(t *testing.T) {
	fontCache, err := LoadFonts(testFontsDir)
	if err != nil {
		t.Fatalf("LoadFonts: %v", err)
	}
	data := &ResumeData{
		Name:  "Alice Example",
		Title: "Engineer",
		Font:  "inter",
		Contact: ContactInfo{
			Email:    "alice@example.com",
			Location: "Berlin",
		},
		Sections: []Section{
			{Type: "text", Title: "Summary", Content: "Hello world"},
		},
	}
	html, err := RenderTemplate("classic-navy", testThemesDir, fontCache, data)
	if err != nil {
		t.Fatalf("RenderTemplate: %v", err)
	}
	if !strings.Contains(html, "Alice Example") {
		t.Error("rendered HTML missing name")
	}
	if strings.Contains(html, "{{THEME_CSS}}") {
		t.Error("theme CSS placeholder not replaced")
	}
	if strings.Contains(html, "{{FONT_CSS}}") {
		t.Error("font CSS placeholder not replaced")
	}
	if !strings.Contains(html, "@font-face") {
		t.Error("expected real font CSS (base64 @font-face) in rendered HTML")
	}
}

func TestRenderTemplate_InvalidThemeFormat(t *testing.T) {
	_, err := RenderTemplate("classic", testThemesDir, map[string]string{}, &ResumeData{})
	if err == nil {
		t.Fatal("expected error for theme without variant, got nil")
	}
}

func TestRenderTemplate_UnknownTheme(t *testing.T) {
	_, err := RenderTemplate("classic-nope", testThemesDir, map[string]string{}, &ResumeData{})
	if err == nil {
		t.Fatal("expected error for unknown theme variant, got nil")
	}
}
