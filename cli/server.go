package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type PreviewServer struct {
	yamlPath  string
	themeName string
	themesDir string
	fontCache map[string]string
}

func NewPreviewServer(yamlPath, themeName, themesDir string, fontCache map[string]string) *PreviewServer {
	return &PreviewServer{
		yamlPath:  yamlPath,
		themeName: themeName,
		themesDir: themesDir,
		fontCache: fontCache,
	}
}

func (s *PreviewServer) renderOnce() (string, error) {
	yamlBytes, err := os.ReadFile(s.yamlPath)
	if err != nil {
		return "", fmt.Errorf("read yaml: %w", err)
	}
	data, err := ParseYAML(string(yamlBytes))
	if err != nil {
		return "", err
	}
	data.HideWatermark = true
	return RenderTemplate(s.themeName, s.themesDir, s.fontCache, data)
}

func (s *PreviewServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	html, err := s.renderOnce()
	if err != nil {
		log.Printf("render error: %v", err)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "<h1>Render error</h1><pre>%s</pre>", err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func (s *PreviewServer) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleIndex)
	return mux
}

// FontsDir returns the fonts directory path inferred from themesDir.
func FontsDir(themesDir string) string {
	return filepath.Join(themesDir, "fonts")
}
