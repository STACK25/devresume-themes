package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	var (
		yamlPath      = flag.String("yaml", "../_examples/sample.yaml", "Path to resume YAML file")
		themeOverride = flag.String("theme", "", "Override theme from YAML (format: template-variant, e.g. classic-navy). If empty, use YAML's theme field.")
		port          = flag.Int("port", 7171, "HTTP port for preview server")
		themesDir     = flag.String("themes-dir", "..", "Directory containing template folders and fonts/")
	)
	flag.Parse()

	fontCache, err := LoadFonts(FontsDir(*themesDir))
	if err != nil {
		log.Fatalf("load fonts: %v", err)
	}
	log.Printf("loaded %d font(s)", len(fontCache))

	reload, err := WatchPaths(*yamlPath, *themesDir)
	if err != nil {
		log.Fatalf("start watcher: %v", err)
	}

	srv := NewPreviewServer(*yamlPath, *themeOverride, *themesDir, fontCache, reload)

	addr := fmt.Sprintf("localhost:%d", *port)
	if *themeOverride != "" {
		log.Printf("devresume-themes preview: http://%s  (yaml=%s, theme override=%s)", addr, *yamlPath, *themeOverride)
	} else {
		log.Printf("devresume-themes preview: http://%s  (yaml=%s, theme from yaml)", addr, *yamlPath)
	}
	log.Printf("watching yaml, all template folders, and fonts/; edit any file to trigger reload")
	if err := http.ListenAndServe(addr, srv.Routes()); err != nil {
		log.Fatal(err)
	}
}
