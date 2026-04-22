package main

import (
	"strings"
	"testing"
)

const testFontsDir = "../fonts"

func TestLoadFonts_Inter(t *testing.T) {
	cache, err := LoadFonts(testFontsDir)
	if err != nil {
		t.Fatalf("LoadFonts: %v", err)
	}
	css, ok := cache["inter"]
	if !ok {
		t.Fatal("inter font missing from cache")
	}
	if !strings.Contains(css, "@font-face") {
		t.Error("expected @font-face rule in inter CSS")
	}
	if !strings.Contains(css, "--font-resume") {
		t.Error("expected --font-resume CSS variable")
	}
	if !strings.Contains(css, "base64,") {
		t.Error("expected base64-inlined font data")
	}
}

func TestGetFontCSS_FallbackToDefault(t *testing.T) {
	cache, err := LoadFonts(testFontsDir)
	if err != nil {
		t.Fatalf("LoadFonts: %v", err)
	}
	css, ok := GetFontCSS(cache, "nonexistent-font")
	if !ok {
		t.Fatal("expected fallback to default font")
	}
	if css == "" {
		t.Error("fallback CSS should not be empty")
	}
}

func TestGetFontCSS_EmptyIDUsesDefault(t *testing.T) {
	cache, err := LoadFonts(testFontsDir)
	if err != nil {
		t.Fatalf("LoadFonts: %v", err)
	}
	css1, _ := GetFontCSS(cache, "")
	css2, _ := GetFontCSS(cache, "inter")
	if css1 != css2 {
		t.Error("empty fontID should resolve to inter (default)")
	}
}
