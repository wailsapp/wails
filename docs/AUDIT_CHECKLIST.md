# Wails v3 Documentation Audit - Action Checklist

Generated: November 22, 2025

## Phase 1: Cleanup & Consolidation (Weeks 1-2)

### [ ] Task 1.1: Review & Delete "Learn" Folder Content
- [ ] Audit each of 24 files in `/learn/` folder
- [ ] Identify which are purely redundant
- [ ] Identify which have unique content to migrate
- [ ] Backup `/learn/` folder before deletion
- [ ] Delete redundant files (windows.mdx, application-menu.mdx, etc.)
- [ ] Migrate unique content (keybindings → /features/keyboard.mdx, etc.)
- [ ] Verify no broken links after deletion
- [ ] Delete entire `/learn/` folder

**Files to DELETE (redundant):**
- windows.mdx, application-menu.mdx, context-menu.mdx, systray.mdx
- dialogs.mdx, bindings.mdx, binding-system.mdx, binding-best-practices.mdx
- advanced-binding.mdx, services.mdx, clipboard.mdx, events.mdx
- screens.mdx, menu-reference.mdx, build.mdx
- (Total: 15-20 files)

**Files to MIGRATE (unique content):**
- keybindings.mdx → /features/keyboard.mdx
- notifications.mdx → /features/notifications.mdx
- environment.mdx → /features/environment.mdx
- dock.mdx → /features/platform/macos-dock.mdx
- browser.mdx → Decide: keep or migrate to /reference/

### [ ] Task 1.2: Remove Orphaned "Getting-Started" Folder
- [ ] Verify no links to `/getting-started/*` in sidebar
- [ ] Check for any incoming links from external documentation
- [ ] Delete `/getting-started/installation.mdx`
- [ ] Delete `/getting-started/your-first-app.mdx`
- [ ] Delete `/getting-started/` directory
- [ ] Verify `/quick-start/*` covers all onboarding scenarios

### [ ] Task 1.3: Archive Old Blog Content
- [ ] Create `/blog/archive/` directory
- [ ] Move v2 release notes to archive:
  - 2021-09-27-v2-beta1-release-notes.md
  - 2021-11-08-v2-beta2-release-notes.md
  - 2022-02-22-v2-beta3-release-notes.md
  - 2022-09-22-v2-release-notes.md
- [ ] Review/update or archive 2023-01-17-v3-roadmap.md
- [ ] Verify build system correctly ignores archived posts

---

## Phase 2: Structure Alignment (Weeks 3-4)

### [ ] Task 2.1: Create Guides Subdirectories & Reorganize
- [ ] Create `/guides/dev/` directory
- [ ] Create `/guides/build/` directory
- [ ] Create `/guides/distribution/` directory
- [ ] Create `/guides/patterns/` directory
- [ ] Create `/guides/advanced/` directory
- [ ] Reorganize existing guides into subdirectories:
  - building.mdx, cross-platform.mdx, signing.mdx → `/guides/build/`
  - auto-updates.mdx, single-instance.mdx, file-associations.mdx → `/guides/distribution/`
  - gin-routing.mdx, gin-services.mdx, menus.mdx → `/guides/patterns/`
  - custom-templates.mdx, panic-handling.mdx → `/guides/advanced/`
  - testing.mdx, e2e-testing.mdx → `/guides/dev/`
- [ ] Create missing stubs:
  - `/guides/dev/debugging.mdx` (stub, expand later)
  - `/guides/dev/project-structure.mdx` (stub, expand later)
  - `/guides/build/windows.mdx`, `/macos.mdx`, `/linux.mdx` (stubs)
  - `/guides/patterns/database.mdx`, `/rest-api.mdx` (stubs)
  - `/guides/advanced/security.mdx` (enhance from /guides/security.mdx)
  - `/guides/advanced/wml.mdx` (stub)
- [ ] Update all internal links
- [ ] Update astro.config.mjs sidebar references (should still work with autogenerate)
- [ ] Test build and verify sidebar renders correctly

### [ ] Task 2.2: Create Reference Subdirectories
- [ ] Create `/reference/application/` directory
- [ ] Create `/reference/window/` directory
- [ ] Create `/reference/menu/` directory
- [ ] Create `/reference/events/` directory
- [ ] Create `/reference/dialogs/` directory
- [ ] Create `/reference/runtime/` directory
- [ ] Create `/reference/cli/` directory
- [ ] Move existing reference files:
  - `/reference/application-api.mdx` → `/reference/application/overview.mdx`
  - `/reference/window-api.mdx` → `/reference/window/overview.mdx`
  - `/reference/menu-api.mdx` → `/reference/menu/overview.mdx`
- [ ] Create autogenerate-ready stub files in each directory
- [ ] Test sidebar autogeneration
- [ ] Verify no broken links

### [ ] Task 2.3: Reorganize Contributing Documentation
- [ ] Create `/contributing/architecture/` directory
- [ ] Create `/contributing/codebase/` directory
- [ ] Create `/contributing/workflows/` directory
- [ ] Move/organize existing contributing files:
  - `/contributing/architecture.mdx` → `/contributing/architecture/overview.mdx`
  - `/contributing/codebase-layout.mdx` → `/contributing/codebase/structure.mdx`
- [ ] Create missing stubs:
  - `/contributing/getting-started.mdx`
  - `/contributing/setup.mdx`
  - `/contributing/standards.mdx`
  - `/contributing/pull-requests.mdx`
  - `/contributing/documentation.mdx`
  - `/contributing/architecture/cli.mdx`, `/runtime.mdx`, `/platform.mdx`, `/build.mdx`
  - `/contributing/codebase/application.mdx`, `/internal.mdx`, `/platform.mdx`, `/testing.mdx`
  - `/contributing/workflows/building.mdx`, `/testing.mdx`, `/debugging.mdx`, `/features.mdx`, `/bugs.mdx`
- [ ] Update sidebar configuration with new nested structure
- [ ] Test build and sidebar rendering

---

## Phase 3: Documentation Completion (Weeks 5-8)

### [ ] Task 3.1: Create Missing Feature Documentation (Priority 1 - Core)

**Drag & Drop:**
- [ ] Create `/features/drag-drop.mdx`
- [ ] Include: basic setup, file drop handling, custom drop zones
- [ ] Add code examples (Go + JavaScript)
- [ ] Reference implementation details

**Keyboard Shortcuts:**
- [ ] Create `/features/keyboard.mdx`
- [ ] Migrate content from `/learn/keybindings.mdx` (if not already deleted)
- [ ] Include: accelerators, key bindings, platform-specific shortcuts
- [ ] Add examples for all platforms

**Environment Variables:**
- [ ] Create `/features/environment.mdx`
- [ ] Include: accessing environment, app config, system info
- [ ] Add cross-platform examples

**Notifications:**
- [ ] Create `/features/notifications.mdx`
- [ ] Include: notification display, callbacks, styling
- [ ] Add examples for all platforms

### [ ] Task 3.2: Create Event System Documentation (Priority 2)
- [ ] Create `/features/events/application.mdx`
  - App lifecycle events (startup, shutdown, etc.)
  
- [ ] Create `/features/events/window.mdx`
  - Window events (close, resize, focus, etc.)
  
- [ ] Create `/features/events/custom.mdx`
  - Creating custom events
  - Emitting events
  - Listening to custom events

### [ ] Task 3.3: Create Platform-Specific Features (Priority 3)
- [ ] Create `/features/platform/macos-dock.mdx`
  - Dock icon, badge, menu
  - Launch on startup
  - Keyboard shortcuts
  
- [ ] Create `/features/platform/windows-uac.mdx`
  - Move `/guides/windows-uac.mdx` → `/features/platform/windows-uac.mdx`
  - Update links
  
- [ ] Create `/features/platform/linux.mdx`
  - Linux-specific features
  - Desktop integration
  - System tray
  
- [ ] Create `/features/platform/macos-toolbar.mdx` (if implemented)
  - Or mark as "Coming Soon"

### [ ] Task 3.4: Expand Reference/API Documentation
- [ ] Create `/reference/events/overview.mdx` with:
  - Event types
  - Event signatures
  - Examples for each event type
  
- [ ] Create `/reference/dialogs/overview.mdx` with:
  - Dialog options
  - Return types
  - Examples
  
- [ ] Create `/reference/runtime/overview.mdx` with:
  - JavaScript runtime API
  - Available functions
  - Type definitions
  
- [ ] Create `/reference/cli/overview.mdx` with:
  - Move content from `/guides/cli.mdx`
  - Or maintain in guides and link from reference

---

## Phase 4: Quality Improvement (Weeks 9-10)

### [ ] Task 4.1: Create Missing Tutorials
- [ ] Create `/tutorials/react-todo.mdx` - React example
- [ ] Create `/tutorials/vue-todo.mdx` - Vue example
- [ ] Create `/tutorials/svelte-todo.mdx` - Svelte example
- [ ] Create `/tutorials/multi-window-app.mdx` - Complex example

### [ ] Task 4.2: Expand Troubleshooting
- [ ] Create `/troubleshooting/build-errors.mdx` - Common build issues
- [ ] Create `/troubleshooting/windows.mdx` - Windows-specific issues
- [ ] Create `/troubleshooting/linux.mdx` - Linux-specific issues
- [ ] Create `/troubleshooting/macos.mdx` - macOS-specific issues
- [ ] Create `/troubleshooting/debugging.mdx` - Debugging guide
- [ ] Create `/troubleshooting/performance.mdx` - Performance issues

### [ ] Task 4.3: Create Migration Guides
- [ ] Create `/migration/from-electron.mdx` with:
  - Architecture comparison
  - API mapping guide
  - Migration patterns
  - Common gotchas
  - Example migration

### [ ] Task 4.4: Consolidate Platform-Specific Documentation
- [ ] Audit scattered platform docs
- [ ] Consolidate into `/features/platform/` structure
- [ ] Move `/guides/windows-uac.mdx` (done in Phase 3)
- [ ] Ensure consistent naming and organization

---

## Final Verification & Testing

### [ ] Verification Checklist
- [ ] Build runs without errors: `npm run build`
- [ ] No broken links in sidebar
- [ ] No 404s in generated HTML
- [ ] Search index includes all new files
- [ ] Mobile navigation works correctly
- [ ] All code examples are syntactically correct
- [ ] No console errors in dev tools
- [ ] Sidebar structure matches expected organization

### [ ] Quality Assurance
- [ ] Spellcheck all new files
- [ ] Grammar review of new content
- [ ] Verify code examples work/compile
- [ ] Check consistency in terminology
- [ ] Verify all links point to valid files
- [ ] Review sidebar order and grouping

### [ ] Documentation
- [ ] Update this checklist with completion dates
- [ ] Document any deviations from plan
- [ ] Note lessons learned
- [ ] Estimate timeline accuracy

---

## Metrics Tracking

### Before Audit
- Total files: 142
- Redundant files: 24
- Missing files: 56
- Lines of duplicate content: 7,125
- Broken sidebar links: 56+

### Target After Audit
- Total files: 160+
- Redundant files: 0
- Missing files: 0
- Lines of duplicate content: 0
- Broken sidebar links: 0

### Progress Tracking
- [ ] Phase 1 Complete (Cleanup)
- [ ] Phase 2 Complete (Reorganization)
- [ ] Phase 3 Complete (New Docs)
- [ ] Phase 4 Complete (Enhancement)
- [ ] Verification & Testing Complete

---

## Notes & Blockers

### Current Blockers
- (None identified - ready to proceed)

### Questions to Clarify
- Should `/guides/menus.mdx` be kept or consolidated into `/features/menus/`?
- Should `/learn/browser.mdx` be migrated or is it outdated for v3?
- Should DEVELOPER_GUIDE.md be migrated to contributing section?
- What framework examples should tutorials prioritize?

### Reference Documents
- Full Audit Report: `DOCUMENTATION_AUDIT_REPORT.md`
- Sidebar Config: `astro.config.mjs`
- Current Docs Structure: `/src/content/docs/`

---

**Status:** Ready to begin Phase 1  
**Last Updated:** November 22, 2025  
**Next Review:** After Phase 1 completion
