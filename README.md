# devresume-themes

Themes, HTML templates, and bundled fonts for [devResume.io](https://devresume.io), a free, open-source "resume as code" tool for tech professionals.

This repo is consumed by the main [devresume](https://github.com/STACK25/devresume) backend as a git submodule. You can also use it standalone to develop new themes or tweak fonts.

## Repository Layout

```
classic/  compact/  modern/  technical/
  template.html          HTML scaffold (Go html/template)
  themes/*.css           color variants

fonts/                   bundled fonts (OFL 1.1; see fonts/LICENSES.md)
  inter/  jetbrains-mono/  merriweather/  ibm-plex-sans/
    meta.json, *.woff2

cli/                     Go preview CLI (see cli/README.md)
_examples/sample.yaml    example resume data for local preview
```

## Creating or Modifying a Theme

1. Pick the closest existing template (`classic`, `compact`, `modern`, `technical`).
2. Copy an existing variant CSS as a starting point:
   ```bash
   cp classic/themes/navy.css classic/themes/my-theme.css
   ```
3. Run the CLI against your variant (see `cli/README.md`).
4. Iterate on the CSS. The browser live-reloads on every save.
5. Open a PR.

## Template Data Contract

Templates are `html/template` files receiving a `ResumeData` struct. See [`cli/model.go`](cli/model.go) for the exact shape.

Two placeholders in every `template.html` get replaced before Go-template parsing:

- `{{FONT_CSS}}`: base64-inlined @font-face rules for the selected font, plus a `--font-resume` CSS variable.
- `{{THEME_CSS}}`: the selected `themes/<variant>.css` inlined verbatim.

The CLI and backend produce byte-identical output for these placeholders, so previews match production.

## Licensing

- **Code, templates, CSS, documentation:** MIT. See `LICENSE`.
- **Fonts in `fonts/`:** under their own SIL Open Font License 1.1 (with per-font license files inside each directory). See `fonts/LICENSES.md` for the summary.
