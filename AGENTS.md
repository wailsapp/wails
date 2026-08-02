# AI Agent Instructions for Wails v3

## Managing AI-Generated Planning Documents

AI assistants often create planning and design documents during development:
- PLAN.md, IMPLEMENTATION.md, ARCHITECTURE.md
- DESIGN.md, CODEBASE_SUMMARY.md, INTEGRATION_PLAN.md
- TESTING_GUIDE.md, TECHNICAL_DESIGN.md, and similar files

**Best Practice: Use a dedicated directory for these ephemeral files**

**Recommended approach:**
- Create a `history/` directory in the project root
- Store ALL AI-generated planning/design docs in `history/`
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

### Important Rules

- Store AI planning docs in `history/` directory
- **ALWAYS run `coderabbit --plain` before committing** to get code analysis and catch issues early
- Do NOT clutter repo root with planning documents

## Implementation Tracking (IMPLEMENTATION.md)

**IMPORTANT**: The `IMPLEMENTATION.md` file at the repository root is a **persistent tracking document** for the GTK4 / WebKitGTK 6.0 / GTK3-legacy implementation work. It is NOT an ephemeral planning document.

As of 2026-05-16 (issue #5459), GTK4 + WebKitGTK 6.0 is the **default** Linux stack; GTK3 + WebKit2GTK 4.1 is a legacy opt-in (`-tags gtk3`) for one v3 cycle and is scheduled for removal in v3.1. The default-flip rationale is recorded in `IMPLEMENTATION.md` Decision 1.1.

### Requirements

1. **Update with EVERY commit** that touches GTK4/WebKitGTK 6.0 or legacy GTK3 code
2. **Track all architectural decisions** with context, decision, and rationale
3. **Maintain progress status** for each implementation phase
4. **Document API differences** between the GTK4 default and GTK3 legacy paths
5. **Keep file references** accurate and up-to-date

### What to Update

- Phase completion status (✅ COMPLETE, 🔄 IN PROGRESS, 📋 PENDING)
- New decisions made during implementation
- Files created or modified
- Changelog entries with dates
- TODO items discovered during work

### Commit Message Pattern

When updating IMPLEMENTATION.md:
```
docs: update implementation tracker for [phase/feature]
```

## Landing the Plane (Session Completion)

**When ending a work session**, you MUST complete ALL steps below. Work is NOT complete until `git push` succeeds.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **PUSH TO REMOTE** - This is MANDATORY:
   ```bash
   git pull --rebase
   git push
   git status  # MUST show "up to date with origin"
   ```
5. **Clean up** - Clear stashes, prune remote branches
6. **Verify** - All changes committed AND pushed
7. **Hand off** - Provide context for next session

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing - that leaves work stranded locally
- NEVER say "ready to push when you are" - YOU must push
- If push fails, resolve and retry until it succeeds
