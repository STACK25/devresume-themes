package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const DefaultFontID = "inter"

type fontWeight struct {
	File   string `json:"file"`
	Weight int    `json:"weight"`
	Style  string `json:"style"`
}

type fontMeta struct {
	ID       string       `json:"id"`
	Label    string       `json:"label"`
	Category string       `json:"category"`
	Family   string       `json:"family"`
	Fallback string       `json:"fallback"`
	Weights  []fontWeight `json:"weights"`
}

// LoadFonts scans fontsDir and returns a map of fontID to @font-face CSS
// (base64-inlined, identical output to the backend fonts service).
func LoadFonts(fontsDir string) (map[string]string, error) {
	entries, err := os.ReadDir(fontsDir)
	if err != nil {
		return nil, fmt.Errorf("read fonts dir %s: %w", fontsDir, err)
	}

	cache := make(map[string]string)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		fontID := entry.Name()
		metaPath := filepath.Join(fontsDir, fontID, "meta.json")
		metaBytes, err := os.ReadFile(metaPath)
		if err != nil {
			continue // dir without meta.json is ignored
		}
		var meta fontMeta
		if err := json.Unmarshal(metaBytes, &meta); err != nil {
			return nil, fmt.Errorf("parse %s: %w", metaPath, err)
		}
		css, err := buildFontCSS(fontsDir, fontID, &meta)
		if err != nil {
			return nil, fmt.Errorf("build CSS for %s: %w", fontID, err)
		}
		cache[meta.ID] = css
	}
	if len(cache) == 0 {
		return nil, fmt.Errorf("no fonts found in %s", fontsDir)
	}
	return cache, nil
}

// GetFontCSS returns the CSS for fontID, falling back to DefaultFontID if
// fontID is empty or unknown. Second return is false only when the default
// is also missing (effectively: no fonts loaded at all).
func GetFontCSS(cache map[string]string, fontID string) (string, bool) {
	if fontID == "" {
		fontID = DefaultFontID
	}
	if css, ok := cache[fontID]; ok {
		return css, true
	}
	css, ok := cache[DefaultFontID]
	return css, ok
}

func buildFontCSS(fontsDir, fontID string, meta *fontMeta) (string, error) {
	var b strings.Builder
	for _, w := range meta.Weights {
		path := filepath.Join(fontsDir, fontID, w.File)
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("read %s: %w", path, err)
		}
		encoded := base64.StdEncoding.EncodeToString(data)
		fmt.Fprintf(&b, "@font-face{font-family:'%s';src:url(data:font/woff2;base64,%s) format('woff2');font-weight:%d;font-style:%s;font-display:swap;}\n",
			meta.Family, encoded, w.Weight, w.Style)
	}
	fallback := meta.Fallback
	if fallback == "" {
		fallback = "sans-serif"
	}
	fmt.Fprintf(&b, ":host,:root{--font-resume:'%s',%s;}\n", meta.Family, fallback)
	return b.String(), nil
}
