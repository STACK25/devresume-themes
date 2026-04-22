package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	var (
		yamlPath  = flag.String("yaml", "../_examples/sample.yaml", "Path to resume YAML file")
		themeName = flag.String("theme", "classic-navy", "Theme in format template-variant")
		port      = flag.Int("port", 7171, "HTTP port for preview server")
		themesDir = flag.String("themes-dir", "..", "Directory containing template folders and fonts/")
	)
	flag.Parse()

	fontCache, err := LoadFonts(FontsDir(*themesDir))
	if err != nil {
		log.Fatalf("load fonts: %v", err)
	}
	log.Printf("loaded %d font(s)", len(fontCache))

	srv := NewPreviewServer(*yamlPath, *themeName, *themesDir, fontCache)

	addr := fmt.Sprintf("localhost:%d", *port)
	log.Printf("devresume-themes preview: http://%s  (yaml=%s theme=%s)", addr, *yamlPath, *themeName)
	if err := http.ListenAndServe(addr, srv.Routes()); err != nil {
		log.Fatal(err)
	}
}
