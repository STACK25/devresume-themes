# devresume-themes CLI

Local preview tool for developing devResume themes. Renders with the real bundled fonts, so what you see matches production pixel-for-pixel.

## Usage

From the repo root:

```bash
cd cli
go run . --theme classic-navy
```

Opens a live-reloading preview at http://localhost:7171.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--yaml` | `../_examples/sample.yaml` | YAML resume file to render |
| `--theme` | `classic-navy` | `template-variant` (e.g. `modern-teal`, `technical-dark`) |
| `--port` | `7171` | HTTP port |
| `--themes-dir` | `..` | Directory containing template folders and `fonts/` |

## Live Reload

The CLI watches:
- The YAML file passed via `--yaml`
- The selected template's folder (`<themes-dir>/<template>/`)
- Its `themes/` subdirectory
- The `<themes-dir>/fonts/` directory and each font subdirectory

Any write to these triggers a browser reload via Server-Sent Events. Font cache is rebuilt automatically on font edits.

Note: on macOS, `touch` without content change produces a Chmod-only event that is intentionally ignored. Editor saves (which write content) trigger reload as expected.

## What It Does Not Do

- No PDF export (production uses Gotenberg; not needed for theme development)

## Building a Binary

```bash
cd cli
go build -o devresume-themes .
./devresume-themes --theme modern-teal
```
