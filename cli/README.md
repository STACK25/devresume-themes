# devresume-themes CLI

Local preview tool for developing devResume themes. Renders with the real bundled fonts, so what you see matches production pixel-for-pixel.

## Usage

From the repo root:

```bash
cd cli
go run .
```

Opens a live-reloading preview at http://localhost:7171. The theme and font come from the YAML's `theme:` and `font:` fields; change them and the browser reloads.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--yaml` | `../_examples/sample.yaml` | YAML resume file to render |
| `--theme` | *empty* | Override the YAML's `theme:` field (format `template-variant`, e.g. `modern-teal`). Leave empty to use the YAML. |
| `--port` | `7171` | HTTP port |
| `--themes-dir` | `..` | Directory containing template folders and `fonts/` |

## Live Reload

The CLI watches:
- The YAML file passed via `--yaml`
- Every template folder under `<themes-dir>/` (and its `themes/` subdirectory)
- The `<themes-dir>/fonts/` directory and each font subdirectory

Any write to these triggers a browser reload via Server-Sent Events. Font cache is rebuilt automatically on font edits. Because all template folders are watched, you can switch themes by editing the YAML's `theme:` field; the preview updates without a CLI restart.

Note: on macOS, `touch` without content change produces a Chmod-only event that is intentionally ignored. Editor saves (which write content) trigger reload as expected.

## What It Does Not Do

- No PDF export (production uses Gotenberg; not needed for theme development)

## Building a Binary

```bash
cd cli
go build -o devresume-themes .
./devresume-themes
```
