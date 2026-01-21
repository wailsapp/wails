# AI Agent Instructions for Wails v3

## Issue Tracking with bd (beads)

**IMPORTANT**: This project uses **bd (beads)** for ALL issue tracking. Do NOT use markdown TODOs, task lists, or other tracking methods.

### Why bd?

- Dependency-aware: Track blockers and relationships between issues
- Git-friendly: Auto-syncs to JSONL for version control
- Agent-optimized: JSON output, ready work detection, discovered-from links
- Prevents duplicate tracking systems and confusion

### Quick Start

**Check for ready work:**
```bash
bd ready --json
```

**Create new issues:**
```bash
bd create "Issue title" -t bug|feature|task -p 0-4 --json
bd create "Issue title" -p 1 --deps discovered-from:bd-123 --json
bd create "Subtask" --parent <epic-id> --json  # Hierarchical subtask (gets ID like epic-id.1)
```

**Claim and update:**
```bash
bd update bd-42 --status in_progress --json
bd update bd-42 --priority 1 --json
```

**Complete work:**
```bash
bd close bd-42 --reason "Completed" --json
```

### Issue Types

- `bug` - Something broken
- `feature` - New functionality
- `task` - Work item (tests, docs, refactoring)
- `epic` - Large feature with subtasks
- `chore` - Maintenance (dependencies, tooling)

### Priorities

- `0` - Critical (security, data loss, broken builds)
- `1` - High (major features, important bugs)
- `2` - Medium (default, nice-to-have)
- `3` - Low (polish, optimization)
- `4` - Backlog (future ideas)

### Workflow for AI Agents

1. **Check ready work**: `bd ready` shows unblocked issues
2. **Claim your task**: `bd update <id> --status in_progress`
3. **Work on it**: Implement, test, document
4. **Discover new work?** Create linked issue:
   - `bd create "Found bug" -p 1 --deps discovered-from:<parent-id>`
5. **Complete**: `bd close <id> --reason "Done"`
6. **Commit together**: Always commit the `.beads/issues.jsonl` file together with the code changes so issue state stays in sync with code state

### Auto-Sync

bd automatically syncs with git:
- Exports to `.beads/issues.jsonl` after changes (5s debounce)
- Imports from JSONL when newer (e.g., after `git pull`)
- No manual export/import needed!

### GitHub Copilot Integration

If using GitHub Copilot, also create `.github/copilot-instructions.md` for automatic instruction loading.
Run `bd onboard` to get the content, or see step 2 of the onboard instructions.

### MCP Server (Recommended)

If using Claude or MCP-compatible clients, install the beads MCP server:

```bash
pip install beads-mcp
```

Add to MCP config (e.g., `~/.config/claude/config.json`):
```json
{
  "beads": {
    "command": "beads-mcp",
    "args": []
  }
}
```

Then use `mcp__beads__*` functions instead of CLI commands.

### Managing AI-Generated Planning Documents

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

### CLI Help

Run `bd <command> --help` to see all available flags for any command.
For example: `bd create --help` shows `--parent`, `--deps`, `--assignee`, etc.

### Important Rules

- Use bd for ALL task tracking
- Always use `--json` flag for programmatic use
- Link discovered work with `discovered-from` dependencies
- Check `bd ready` before asking "what should I work on?"
- Store AI planning docs in `history/` directory
- Run `bd <cmd> --help` to discover available flags
- **ALWAYS run `coderabbit --plain` before committing** to get code analysis and catch issues early
- Do NOT create markdown TODO lists
- Do NOT use external issue trackers
- Do NOT duplicate tracking systems
- Do NOT clutter repo root with planning documents

For more details, see README.md and QUICKSTART.md.

## Event System

### Adding New Events

**IMPORTANT**: Events are code-generated. Do NOT manually edit the generated event files.

#### Event Files (Generated - DO NOT EDIT)

These files are automatically generated from `v3/pkg/events/events.txt`:
- `v3/pkg/events/events.go` - Go event types and constants
- `v3/pkg/events/events_darwin.h` - macOS C header defines
- `v3/pkg/events/events_linux.h` - Linux C header defines
- `v3/pkg/events/events_ios.h` - iOS C header defines
- `v3/internal/runtime/desktop/@wailsio/runtime/src/event_types.ts` - TypeScript event types

#### How to Add New Events

1. **Edit the source file**: `v3/pkg/events/events.txt`
   - Format: `platform:EventName`
   - Platforms: `common`, `mac`, `windows`, `linux`, `ios`
   - Add `!` suffix to prevent auto-generating delegate methods (e.g., `mac:WindowDidZoom!`)

2. **Run the generator**:
   ```bash
   cd v3/tasks/events
   go run generate.go
   go fmt ../../pkg/events/events.go
   ```

   Or use the task runner:
   ```bash
   cd v3/internal/runtime
   task generate:events
   ```

3. **Rebuild the JavaScript runtime** (if needed):
   ```bash
   cd v3/internal/runtime
   task build
   ```

#### Event Naming Conventions

- `common:` - Cross-platform events (emitted on all platforms)
- `mac:` - macOS-specific events
- `windows:` - Windows-specific events
- `linux:` - Linux-specific events
- `ios:` - iOS-specific events

#### Example: Adding a New Common Event

```
# In v3/pkg/events/events.txt, add:
common:WindowMouseEnter
common:WindowMouseLeave
```

Then run the generator to create:
- Go: `events.Common.WindowMouseEnter`, `events.Common.WindowMouseLeave`
- TypeScript: `"common:WindowMouseEnter"`, `"common:WindowMouseLeave"`

#### Generated File Headers

All generated files include the standard Wails header with both English and Welsh comments:
```
/*
 _       __      _ __
| |     / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The electron alternative for Go
(c) Lea Anthony 2019-present
*/

// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
```
