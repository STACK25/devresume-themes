package main

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// RenderTemplate loads template.html + the named theme CSS from themesDir,
// picks the requested font CSS from fontCache (falling back to default),
// injects both, and executes the html/template with data.
//
// On-disk layout expected (identical to backend):
//
//	{themesDir}/{templateName}/template.html
//	{themesDir}/{templateName}/themes/{variant}.css
func RenderTemplate(themeName, themesDir string, fontCache map[string]string, data *ResumeData) (string, error) {
	parts := strings.SplitN(themeName, "-", 2)
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid theme format %q (want template-variant, e.g. classic-navy)", themeName)
	}
	templateName, variant := parts[0], parts[1]

	tmplPath := filepath.Join(themesDir, templateName, "template.html")
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		return "", fmt.Errorf("read template %s: %w", tmplPath, err)
	}

	themePath := filepath.Join(themesDir, templateName, "themes", variant+".css")
	themeBytes, err := os.ReadFile(themePath)
	if err != nil {
		return "", fmt.Errorf("read theme %s: %w", themePath, err)
	}

	fontCSS, ok := GetFontCSS(fontCache, data.Font)
	if !ok {
		fontCSS = "" // No fonts loaded; render anyway, browser will fall back
	}

	html := string(tmplBytes)
	html = strings.Replace(html, "{{FONT_CSS}}", fontCSS, 1)
	html = strings.Replace(html, "{{THEME_CSS}}", string(themeBytes), 1)

	tmpl, err := template.New("resume").Parse(html)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	return buf.String(), nil
}
