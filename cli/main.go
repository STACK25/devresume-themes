package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		yamlPath  = flag.String("yaml", "../_examples/sample.yaml", "Path to resume YAML file")
		themeName = flag.String("theme", "classic-navy", "Theme in format template-variant (e.g. classic-dark)")
		port      = flag.Int("port", 7171, "HTTP port for preview server")
		themesDir = flag.String("themes-dir", "..", "Directory containing template folders and fonts/")
	)
	flag.Parse()

	fmt.Printf("devresume-themes CLI\n  yaml: %s\n  theme: %s\n  themes-dir: %s\n  port: %d\n",
		*yamlPath, *themeName, *themesDir, *port)

	os.Exit(0)
}
