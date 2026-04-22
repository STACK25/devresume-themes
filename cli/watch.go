package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// WatchPaths monitors the YAML file, the active template folder (+ themes/),
// and the fonts/ tree. It coalesces bursts of filesystem events and sends a
// single tick on the returned channel per burst.
func WatchPaths(yamlPath, themesDir, themeName string) (<-chan struct{}, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	templateName := strings.SplitN(themeName, "-", 2)[0]
	templateDir := filepath.Join(themesDir, templateName)
	themesSubDir := filepath.Join(templateDir, "themes")
	fontsDir := FontsDir(themesDir)

	targets := []string{yamlPath, templateDir, themesSubDir, fontsDir}

	// Also watch each font subdir so we catch woff2/meta.json edits.
	if entries, err := os.ReadDir(fontsDir); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				targets = append(targets, filepath.Join(fontsDir, e.Name()))
			}
		}
	}

	for _, p := range targets {
		if err := w.Add(p); err != nil {
			log.Printf("watch %s: %v (skipping)", p, err)
		}
	}

	out := make(chan struct{}, 1)

	go func() {
		defer w.Close()
		var debounce *time.Timer
		for {
			select {
			case ev, ok := <-w.Events:
				if !ok {
					return
				}
				if ev.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) == 0 {
					continue
				}
				if debounce != nil {
					debounce.Stop()
				}
				debounce = time.AfterFunc(120*time.Millisecond, func() {
					select {
					case out <- struct{}{}:
					default:
					}
				})
			case err, ok := <-w.Errors:
				if !ok {
					return
				}
				log.Printf("watch error: %v", err)
			}
		}
	}()

	return out, nil
}
