---
title: Creating Custom Templates
description: How to generate, customise, and host your own Wails v3 project templates
---

import { Tabs, TabItem, Steps } from '@astrojs/starlight/components';

Wails ships with a set of built-in templates, but you can create your own and share them with the community. A custom template is just a Git repository — once it is hosted publicly, anyone can scaffold a project from it with a single command.

## Generate a template skeleton

The `wails3 generate template` command produces a ready-to-customise template directory:

```bash
wails3 generate template -name MyTemplate
```

All flags:

| Flag           | Description                                        | Default           |
|----------------|----------------------------------------------------|-------------------|
| `-name`        | Template name (required)                           | —                 |
| `-author`      | Author name                                        | —                 |
| `-description` | Short description shown in the CLI                 | —                 |
| `-helpurl`     | URL to documentation for this template             | —                 |
| `-version`     | Initial version                                    | `v0.0.1`          |
| `-frontend`    | Copy an existing frontend directory into the template | —              |
| `-dir`         | Where to write the template directory              | Current directory |

Example with all flags:

```bash
wails3 generate template \
  -name "My Template" \
  -author "Your Name" \
  -description "React + custom setup" \
  -helpurl "https://github.com/yourname/my-template" \
  -version "v1.0.0" \
  -frontend ./my-existing-frontend
```

The generated directory looks like this:

```
MyTemplate/
├── template.yaml          # Template metadata — edit this
├── NEXTSTEPS.md           # Guidance for you as the template author — delete before publishing
├── README.md              # Shown to users after they create a project
├── main.go.tmpl           # Application entry point
├── greetservice.go        # Example Go service
├── go.mod.tmpl            # Go module file
├── go.sum.tmpl            # Go checksums
├── gitignore.tmpl         # Becomes .gitignore in generated projects
├── Taskfile.tmpl.yml      # Build task definitions
└── frontend/              # Your frontend code
```

:::tip[Read NEXTSTEPS.md]
The generated `NEXTSTEPS.md` contains detailed guidance on every part of the template. Read it before customising. Delete it before publishing — it must not appear in projects created from your template.
:::

## Configure template metadata

Open `template.yaml` to set your template's metadata:

```yaml
# yaml-language-server: $schema=https://v3.wails.io/schemas/template.v3.json
name: "My Template"
shortname: my-template
author: Your Name
description: A template with my preferred setup
helpurl: https://github.com/yourname/my-template
version: v1.0.0
wailsVersion: 3
```

The `wailsVersion` field is **required** and must be `3`. The `# yaml-language-server` comment at the top enables autocomplete and inline validation in VS Code (with the [YAML extension](https://marketplace.visualstudio.com/items?itemName=redhat.vscode-yaml)) and JetBrains IDEs — you can leave it or remove it; it has no effect at runtime.

## Customise the template

### Frontend

The `frontend/` directory is copied verbatim into every project created from your template. Replace the placeholder content with your actual frontend:

<Tabs>
  <TabItem label="Start from scratch">
    ```bash
    cd MyTemplate/frontend
    npm create vite@latest .
    ```
    Follow the prompts, then install dependencies:
    ```bash
    npm install
    ```
  </TabItem>
  <TabItem label="Use an existing project">
    Pass `-frontend` when generating the template to copy an existing frontend in one step:
    ```bash
    wails3 generate template -name MyTemplate -frontend ./my-app/frontend
    ```
    Or copy it manually into the `frontend/` directory afterwards.
  </TabItem>
</Tabs>

### Build tasks

`Taskfile.tmpl.yml` defines the build workflow. Update the `install:frontend:deps` and `build:frontend` tasks to match your frontend toolchain:

```yaml
tasks:
  install:frontend:deps:
    dir: frontend
    cmds:
      - npm install       # replace with pnpm install, yarn, etc.

  build:frontend:
    dir: frontend
    deps: [install:frontend:deps, generate:bindings]
    cmds:
      - npm run build     # replace with your build command
```

### Go application

The `main.go.tmpl` file is the application entry point. It is processed by Wails' template engine when a project is created — template variables like `{{.ProductName}}` are replaced with the values the user supplies.

To edit it as a real Go file (with IDE support), temporarily rename it to `main.go`, make your changes, then rename it back to `main.go.tmpl` before committing.

#### Template variables

These variables are available in any `.tmpl` file:

| Variable               | Description                              | Example                         |
|------------------------|------------------------------------------|---------------------------------|
| `{{.ProjectName}}`     | Project name supplied by the user        | `"MyApp"`                       |
| `{{.BinaryName}}`      | Binary filename                          | `"myapp"`                       |
| `{{.ProductName}}`     | Product display name                     | `"My Application"`              |
| `{{.ProductDescription}}` | Product description                   | `"An awesome application"`      |
| `{{.ProductVersion}}`  | Product version                          | `"1.0.0"`                       |
| `{{.ProductCompany}}`  | Company / author name                    | `"My Company Ltd"`              |
| `{{.ProductCopyright}}`| Copyright string                         | `"Copyright 2024 My Company Ltd"` |
| `{{.ProductComments}}` | Additional product comments              | `"Built with Wails"`            |
| `{{.ProductIdentifier}}`| Reverse-DNS product identifier          | `"com.mycompany.myapp"`         |
| `{{.ModulePath}}`      | Go module path                           | `"github.com/you/myapp"`        |
| `{{.WailsVersion}}`    | Wails version used to create the project | `"3.0.0"`                       |
| `{{.Typescript}}`      | `true` if the template name ends in `-ts` | `true`                         |
| `{{.Opn}}`             | Literal `{{` — escape inside templates   | `{{`                            |
| `{{.Cls}}`             | Literal `}}` — escape inside templates   | `}}`                            |

:::tip
Any file in your template can be a `.tmpl` file — including HTML, JSON, and YAML files. Files without the `.tmpl` suffix are copied verbatim.
:::

## Test your template locally

Before publishing, test the template by creating a project from a local path:

```bash
wails3 init -n testproject -t /path/to/MyTemplate
```

Then verify the project works:

```bash
cd testproject
wails3 dev    # development mode with hot reload
wails3 build  # production binary
```

Check that:
- Frontend hot reload works
- Go code changes rebuild and relaunch the app
- The production binary in `bin/` runs correctly

## Publish on GitHub

<Steps>

1. **Create a public GitHub repository** for your template. The repository root must contain `template.yaml`.

2. **Delete `NEXTSTEPS.md`** — this file is guidance for template authors and must not appear in projects users create from your template.

3. **Commit and push** the template directory contents as the repository root:

   ```bash
   git init
   git add .
   git commit -m "Initial template"
   git remote add origin https://github.com/yourname/my-template.git
   git push -u origin main
   ```

4. **Tag a release** using semantic versioning:

   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

</Steps>

Users can now create projects from your template:

```bash
# Latest commit on the default branch
wails3 init -n myapp -t https://github.com/yourname/my-template

# Pinned to a specific release tag
wails3 init -n myapp -t https://github.com/yourname/my-template@v1.0.0
```

:::caution[Third-party template warning]
When a user installs a remote template, Wails displays a warning explaining that the template is third-party code and that the Wails project takes no responsibility for its contents. Users must explicitly confirm before the project is created.

As a template author you are responsible for the security and correctness of all code in your template.
:::

## Best practices

- **Write a clear `README.md`** — this is shown to users after they create a project. Explain how to run, build, and customise the project.
- **Fill in `helpurl`** — link to your repository or dedicated documentation. Users see it in the Wails CLI template list.
- **Pin frontend dependency versions** in `package.json` to avoid broken installs from upstream updates.
- **Test before tagging** — create a fresh project from the tagged release before announcing it to the community.
- **Keep `wailsVersion: 3`** — this field tells Wails which major version the template targets. Do not change it.
- **Update regularly** — keep dependencies current and test against new Wails releases.
