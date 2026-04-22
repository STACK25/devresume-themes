package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

// WatchPaths monitors the YAML file, every template folder (+ its themes/
// subfolder), and the fonts/ tree. All template folders are watched so the
// user can switch themes by editing the 'theme:' field in the YAML without
// restarting the CLI. Events are coalesced via a 120ms debounce.
func WatchPaths(yamlPath, themesDir string) (<-chan struct{}, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	targets := []string{yamlPath}

	// Every subdir of themesDir that contains a template.html is a template
	// folder; watch it plus its themes/ subfolder.
	if entries, err := os.ReadDir(themesDir); err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			dir := filepath.Join(themesDir, e.Name())
			if _, statErr := os.Stat(filepath.Join(dir, "template.html")); statErr != nil {
				continue
			}
			targets = append(targets, dir)
			themesSub := filepath.Join(dir, "themes")
			if _, err := os.Stat(themesSub); err == nil {
				targets = append(targets, themesSub)
			}
		}
	}

	fontsDir := FontsDir(themesDir)
	targets = append(targets, fontsDir)
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
