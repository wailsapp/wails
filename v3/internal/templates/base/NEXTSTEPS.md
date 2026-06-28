# Wails Template — Next Steps

Your template skeleton has been generated. This document explains what was created,
how to customise it, and how to publish it so others can use it.

---

## What Was Generated

```text
<template-name>/
├── template.yaml          # Template metadata (name, author, wailsVersion, etc.)
├── NEXTSTEPS.md           # This file — delete it before publishing
├── README.md              # Shown to users after they create a project
├── main.go.tmpl           # Application entry point (template variables expanded at init time)
├── greetservice.go        # Example Go service bound to the frontend
├── go.mod.tmpl            # Go module file
├── go.sum.tmpl            # Go module checksums
├── gitignore.tmpl         # Becomes .gitignore in the generated project
├── Taskfile.tmpl.yml      # Build task definitions
└── frontend/              # Your frontend code goes here
```

> **Note:** Files ending in `.tmpl` are processed by the template engine when a user
> runs `wails init`. Template variables (e.g. `{{.ProjectName}}`) are replaced with
> the values the user supplies. Files without `.tmpl` are copied verbatim.

---

## Template Variables

The following variables are available inside `.tmpl` files:

| Variable            | Description                                      |
|---------------------|--------------------------------------------------|
| `{{.ProjectName}}`  | The project name supplied by the user (`-n`)     |
| `{{.ModulePath}}`   | The Go module path (`-mod` or derived from `-git`) |
| `{{.WailsVersion}}` | The Wails version used to initialise the project |
| `{{.ProductName}}`  | Product display name                             |
| `{{.ProductDescription}}` | Product description                        |
| `{{.ProductVersion}}` | Product version string                         |
| `{{.ProductCompany}}` | Company / author name                          |
| `{{.ProductIdentifier}}` | Reverse-DNS product identifier              |
| `{{.ProductCopyright}}` | Copyright string                             |
| `{{.Typescript}}`   | `true` for TypeScript templates — declared via `typescript: true` in `template.yaml` (built-in TS templates) or a `-ts` name suffix (community templates) |
| `{{.Opn}}`          | Literal `{{` — use inside templates to avoid parsing errors |
| `{{.Cls}}`          | Literal `}}` — use inside templates to avoid parsing errors |

---

## Customising Your Template

1. **Edit `template.yaml`** — update the `name`, `shortname`, `author`, `description`,
   and `helpurl` fields. The `wailsVersion` field must remain `3`.

2. **Replace `frontend/`** — drop in your framework of choice (Vite, React, Vue, Svelte,
   etc.). The frontend directory is copied verbatim; add `.tmpl` to any file you want
   to have template variables expanded.

3. **Modify `main.go.tmpl`** — adjust the application setup, window options, and
   bound services to match your template's needs.

4. **Update `README.md`** — this file is shown to users after they create a project
   from your template. Make it helpful.

5. **Delete `NEXTSTEPS.md`** — this file is only for you. Remove it before publishing
   so it does not appear in projects users create from your template.

---

## Build Assets

Build assets (platform-specific icons, manifests, signing configs) are **generated
automatically** by `wails init` after your template is applied. You do not need to
include a `build/` directory in your template.

If you want to ship custom build assets (e.g. a specific app icon), add a `build/`
directory to your template. Note that the auto-generated assets will be written on
top of whatever your template provides, so place only files that are not auto-generated.

---

## Publishing on GitHub

1. Commit your template directory as the root of a public GitHub repository.
   Your repository must contain `template.yaml` at the top level.

2. Tag a release:
   ```sh
   git tag v1.0.0
   git push origin v1.0.0
   ```

3. Users can now create projects from your template:
   ```sh
   # Latest commit on the default branch
   wails3 init -n myapp -t https://github.com/yourname/your-template

   # Pinned to a specific release tag
   wails3 init -n myapp -t https://github.com/yourname/your-template@v1.0.0
   ```

---

## Third-Party Template Disclaimer

When a user installs a remote template, Wails displays a warning explaining that the
template is third-party code and that the Wails project takes no responsibility for
its contents. Users must explicitly confirm before the project is created.

As a template author you are responsible for:
- The security and correctness of all code in your template
- Keeping your template up to date with new Wails versions
- Providing clear documentation (`README.md`, `helpurl` in `template.yaml`)

---

## Further Reading

- [Creating Custom Templates](https://v3.wails.io/guides/advanced/custom-templates)
- [Template Schema](https://v3.wails.io/reference/template-json)
- [Wails Init Reference](https://v3.wails.io/reference/cli#init)
