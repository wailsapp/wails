# AI Agent Instructions for Wails

## Repository intent and layout

This repository contains multiple generations of Wails and the supporting project sites:

- [`v2/`](v2/AGENTS.md) is the maintained Wails v2 codebase and release line.
- [`v3/`](v3/AGENTS.md) is the Wails v3 codebase, CLI, runtime, platform integrations, examples, and tests.
- [`docs/`](docs/AGENTS.md) is the current M-Press documentation source for Wails v3.
- [`website/`](website/AGENTS.md) is the Docusaurus website, including the v2 documentation and versioned historical documentation.
- [`scripts/`](scripts/AGENTS.md) and [`tools/`](tools/AGENTS.md) contain repository automation and supporting utilities.
- [`assets/`](assets/AGENTS.md) contains shared artwork and other repository-wide static assets.

Work in the smallest owning directory. Keep v2 and v3 code, documentation, and release assumptions separate unless a task explicitly spans both versions. The root Taskfile delegates to these areas; check the task file in the owning directory for its available commands.

The root is the source of repository-wide policy. A deeper `AGENTS.md` may add constraints for its subtree but does not override these instructions. Preserve existing project conventions and avoid unrelated changes.

## Wails release policy

For Wails release planning and issue/PR triage:

- v3.1.0 contains features only.
- Prioritise bugs by severity and release their fixes as nightly patch updates; keep them outside the v3.1.0 milestone/project.
- Triage issue and PR content rather than inferring priority or release placement from missing labels.
- GitHub Issues and pull requests are the authoritative trackers. Search for existing items before filing, link implementation work to its tracker, and do not create a parallel local issue list.

## Working boundaries

- Keep changes within the requested product/version scope. Do not backport or forward-port between v2 and v3 by assumption; identify compatibility and release intent first.
- Treat generated files according to the owning project's build workflow. When a workflow requires generated artifacts to be committed, update all coupled outputs together.
- Do not edit generated website output or dependency caches as source content. Edit the source and regenerate only when the task requires it.
- Store AI-generated planning documents in `history/`, not at the repository root.
- Do not push unless the user gives explicit manual confirmation.
- Use the `taliesin-ai` identity for commits. Before committing, run `coderabbit --plain` as required by this repository's contribution policy.

## Changelog ownership

- For ordinary v3 contributions, **do not edit `v3/UNRELEASED_CHANGELOG.md` or `docs/mpress/content/changelog.md`**. The former is the automation's staging file; the latter is the published release history. The root `CHANGELOG.md` is only a legacy link, not another v3 changelog.
- Describe the user-visible change accurately in the PR title and description. After merge, `.github/workflows/auto-changelog-v3.yml` generates an entry from the PR title and CodeRabbit walkthrough when available. Internal-only PRs may intentionally produce no entry.
- Release tooling consumes the staged entries, adds a versioned section to the published changelog, and resets the staging file. Do not delete, rename, reset or replace either file to simplify an ordinary PR.
- Changelog automation maintenance and explicitly requested release-note corrections are exceptions: change the owning generator/workflow and any coupled output together, preserve existing entries, and run the changelog publication and release tests.
- v2 retains its separate contributor workflow in `website/src/pages/changelog.mdx`; the v3 automation rule does not change it.

## Wails MCP tools for agents

Wails v3 provides two separate MCP servers. Prefer their connected tools for the tasks they support; if they are unavailable, use the documented CLI and Taskfile commands. Use a CLI built from the checkout you intend to test, rather than assuming an installed release includes local changes.

**Project lifecycle:** launch `wails3 mcp --stdio --root /absolute/path/to/app` through an MCP client. It can inspect projects, run diagnostics, list tasks, start builds/development jobs, and generate bindings. Relevant tools include `wails_project_inspect`, `wails_project_doctor`, `wails_project_task_list`, `wails_project_build`, `wails_project_dev_start`, `wails_project_generate_bindings` and `wails_project_task_run`. Builds return job IDs: use `wails_job_status` to check completion, output and exit status; a job starting is not a successful build. Use `wails_job_stop` for jobs this server started.

The project server is confined to its configured root. Mutating and process-control tools require the session token supplied in the MCP initialization instructions; HTTP transport also requires bearer authentication. A token does not grant permission for actions outside the user's request. See [CLI MCP documentation](docs/mpress/content/guides/cli.md#mcp) and `v3/internal/commands/mcp.go`.

**Running application:** `WAILS_MCP=1 wails3 dev` enables the experimental in-app server. Connect using the endpoint printed by the application (normally `http://127.0.0.1:9099/mcp`). Its tools inspect windows/DOM, interact with the UI, call bound methods and observe events; they do not replace the project build tools. `screenshot_dom` is a structural DOM snapshot, not a pixel screenshot. Keep the `mcp` build tag out of shipped builds. See [application MCP documentation](docs/mpress/content/guides/mcp-service.md).

## Issue Tracking with GitHub

**IMPORTANT**: GitHub Issues and pull requests are the authoritative trackers for this project. Do not create a parallel local issue database or markdown task list.

### Workflow for AI Agents

1. Search existing GitHub issues and pull requests before creating new work.
2. Link implementation work to the relevant issue or pull request.
3. File newly discovered actionable work as a GitHub issue with reproduction details, scope, and acceptance criteria.
4. Use labels, milestones, and cross-references to express priority, release scope, and dependencies.
5. Close issues only after the work is complete and verified; otherwise document why the item is deferred or no longer applicable.

### Important Rules

- Do not duplicate existing GitHub issues.
- Keep issue descriptions current when scope or reproduction steps change.
- Store AI-generated planning documents in `history/`, not the repository root.
- **ALWAYS run `coderabbit --plain` before committing** to catch issues early.
- All commits must use the `taliesin-ai` identity.
- Never push until the user gives explicit manual confirmation.

### Managing AI-Generated Planning Documents

AI assistants often create planning and design documents during development:
- PLAN.md, ARCHITECTURE.md
- DESIGN.md, CODEBASE_SUMMARY.md, INTEGRATION_PLAN.md
- TESTING_GUIDE.md, TECHNICAL_DESIGN.md, and similar files

**Best Practice: Use a dedicated directory for these ephemeral files**

**Recommended approach:**
- Create a `history/` directory in the project root
- Store all ephemeral AI-generated planning/design docs in `history/`
- Keep the repository root clean and focused on permanent project files
- Only access `history/` when explicitly asked to review past planning

**Example .gitignore entry (optional):**
```
# AI planning documents (ephemeral)
history/
```

**Benefits:**
- Clean repository root
- Clear separation between ephemeral and permanent documentation
- Easy to exclude from version control if desired
- Preserves planning history for archeological research
- Reduces noise when browsing the project

## Frontend Runtime: Two Build Outputs

The TypeScript runtime in `v3/internal/runtime/desktop/@wailsio/runtime` produces **two**
independent artifacts, and rebuilding one but not the other is a common and confusing
mistake:

| task | output | consumed by |
|---|---|---|
| `task v3:runtime:build:assets` | `v3/internal/assetserver/bundledassets/runtime.js` (+ `.debug.js`) | the webview, served at `/wails/runtime.js` |
| `task v3:runtime:build:package` | `dist/` in the package directory | an app's frontend, via `node_modules` |

After changing anything under `src/`, rebuild **both**. CI verifies the committed bundles
match `build:assets` output exactly, so the bundles must be committed with the change.

An application imports `@wailsio/runtime` from npm, so it will not see runtime changes made
in this checkout. To test an app against the working tree:

```bash
task v3:install-runtime -- ./path/to/your-app/frontend
```

Undo with `npm install @wailsio/runtime@latest` in the same directory.

## Subsystem References

Some subsystems have a dedicated internals page written for agents. Read the relevant one
before changing that code — several of its decisions look arbitrary until you know which
measured bug they prevent.

- **Streams** (`v3/pkg/application/stream*.go`, `v3/internal/runtime/.../stream.ts`):
  `docs/mpress/content/guides/advanced/streams-internals.md`. Covers the held-poll
  design, the buffer constants and how to pick them, session and connection lifecycle,
  transport selection, and what is unfinished. To convert an existing WebSocket
  implementation, follow `docs/mpress/content/guides/streams-from-websockets.md` —
  a mechanical checklist, including the differences that break silently.

## Landing the Plane (Session Completion)

**When ending a work session**, complete the applicable steps below. Never push without explicit manual confirmation from the user.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **Prepare for remote sync**:
   ```bash
   git status
   ```
   When the worktree is clean and synchronization is intended, run
   `git pull --rebase`.
   Run `git push` only after the user explicitly confirms it.
5. **Clean up** - Review stashes and remove only obsolete ones; prune remote branches
6. **Verify** - All intended changes are present and committed when requested
7. **Hand off** - Provide context for next session

**CRITICAL RULES:**
- Use the `taliesin-ai` identity for every commit.
- Do not push without explicit manual confirmation.
- Report clearly whether changes are uncommitted, committed locally, or pushed.
