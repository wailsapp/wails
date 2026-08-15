# AI Agent Instructions for Wails v3

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

- **Streams** (`pkg/application/stream*.go`, `runtime/.../stream.ts`):
  `docs/src/content/docs/guides/advanced/streams-internals.mdx`. Covers the held-poll
  design, the buffer constants and how to pick them, session and connection lifecycle,
  transport selection, and what is unfinished. To convert an existing WebSocket
  implementation, follow `docs/src/content/docs/guides/streams-from-websockets.mdx` —
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

## Agent skills

### Issue tracker

Planning tickets use local Markdown under `history/wayfinder/`. See
`docs/agents/issue-tracker.md` before creating, claiming, blocking, or resolving
a ticket.

### Triage labels

Local tickets use the canonical triage vocabulary. See
`docs/agents/triage-labels.md`.

### Domain docs

Domain documentation uses a single root context with system ADRs under
