package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

//go:embed assets/reload.js
var reloadJS string

type PreviewServer struct {
	yamlPath     string
	themeOverride string // if non-empty, overrides the theme from YAML
	themesDir    string

	mu        sync.RWMutex
	fontCache map[string]string

	reload      <-chan struct{}
	subMu       sync.Mutex
	subscribers map[chan struct{}]struct{}
}

func NewPreviewServer(yamlPath, themeOverride, themesDir string, fontCache map[string]string, reload <-chan struct{}) *PreviewServer {
	s := &PreviewServer{
		yamlPath:      yamlPath,
		themeOverride: themeOverride,
		themesDir:     themesDir,
		fontCache:     fontCache,
		reload:        reload,
		subscribers:   make(map[chan struct{}]struct{}),
	}
	go s.fanout()
	return s
}

func (s *PreviewServer) fanout() {
	for range s.reload {
		s.refreshFontsIfNeeded()
		s.subMu.Lock()
		for ch := range s.subscribers {
			select {
			case ch <- struct{}{}:
			default:
			}
		}
		s.subMu.Unlock()
		log.Println("reload broadcast")
	}
}

// refreshFontsIfNeeded re-scans fonts/ on every reload tick. Cheap enough
// (base64-encoding four woff2 files is ~20ms total) and keeps preview in
// sync when contributors add or edit a font.
func (s *PreviewServer) refreshFontsIfNeeded() {
	cache, err := LoadFonts(FontsDir(s.themesDir))
	if err != nil {
		log.Printf("refresh fonts: %v", err)
		return
	}
	s.mu.Lock()
	s.fontCache = cache
	s.mu.Unlock()
}

func (s *PreviewServer) subscribe() chan struct{} {
	ch := make(chan struct{}, 1)
	s.subMu.Lock()
	s.subscribers[ch] = struct{}{}
	s.subMu.Unlock()
	return ch
}

func (s *PreviewServer) unsubscribe(ch chan struct{}) {
	s.subMu.Lock()
	delete(s.subscribers, ch)
	s.subMu.Unlock()
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

	// Theme source of truth is the YAML. A non-empty --theme flag overrides it
	// (useful for quickly previewing a theme without editing the yaml).
	theme := data.Theme
	if s.themeOverride != "" {
		theme = s.themeOverride
	}
	if theme == "" {
		return "", fmt.Errorf("no theme specified (set 'theme:' in yaml or pass --theme)")
	}

	s.mu.RLock()
	cache := s.fontCache
	s.mu.RUnlock()

	html, err := RenderTemplate(theme, s.themesDir, cache, data)
	if err != nil {
		return "", err
	}
	injected := "<script>\n" + reloadJS + "\n</script>"
	if idx := strings.LastIndex(html, "</body>"); idx >= 0 {
		html = html[:idx] + injected + html[idx:]
	} else {
		html += injected
	}
	return html, nil
}

func (s *PreviewServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	html, err := s.renderOnce()
	if err != nil {
		log.Printf("render error: %v", err)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `<!DOCTYPE html><html><head><title>devresume-themes: error</title>
<style>body{font:14px/1.5 ui-monospace,monospace;padding:2rem;color:#b91c1c}
pre{background:#fef2f2;padding:1rem;border-radius:8px;white-space:pre-wrap}</style>
</head><body><h1>Render error</h1><pre>%s</pre><script>%s</script></body></html>`,
			err.Error(), reloadJS)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func (s *PreviewServer) handleEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := s.subscribe()
	defer s.unsubscribe(ch)

	fmt.Fprintf(w, "data: connected\n\n")
	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ch:
			fmt.Fprintf(w, "data: reload\n\n")
			flusher.Flush()
		}
	}
}

func (s *PreviewServer) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/events", s.handleEvents)
	return mux
}

func FontsDir(themesDir string) string {
	return filepath.Join(themesDir, "fonts")
}
