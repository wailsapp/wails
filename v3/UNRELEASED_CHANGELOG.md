# Unreleased Changes

<!--
AUTO-GENERATED FILE: do not edit in ordinary contribution PRs.

The Auto Changelog (v3) workflow generates entries after PRs merge.
Describe the user-visible change in the PR title and description instead.
Release tooling moves staged entries into docs/mpress/content/changelog.md
and resets this file. Do not delete or reset it by hand.

For changes to this automation or explicitly requested release-note corrections,
follow the changelog ownership policy in AGENTS.md and preserve existing entries.
-->

## Added
<!-- New features, capabilities, or enhancements -->

## Changed
<!-- Changes in existing functionality -->
- Improve changelog generation speed by reusing compiled regex patterns in [PR](https://github.com/wailsapp/wails/pull/6490) by @taliesin-ai
- Reduce allocations when serving localized fallback pages in [PR](https://github.com/wailsapp/wails/pull/6485) by @taliesin-ai
- Speed up AppImage runtime discovery by avoiding unnecessary metadata lookups in [PR](https://github.com/wailsapp/wails/pull/6482) by @taliesin-ai
- Reuse compiled command routing patterns to improve wake performance in [PR](https://github.com/wailsapp/wails/pull/6475) by @taliesin-ai
- Improve performance by avoiding rescanning of unfinished output lines in [PR](https://github.com/wailsapp/wails/pull/6473) by @taliesin-ai
- Stream executable copies during MSIX package staging to reduce memory use and speed up builds in [PR](https://github.com/wailsapp/wails/pull/6472) by @taliesin-ai
- Avoid copying unfinished command-output lines to improve performance in [PR](https://github.com/wailsapp/wails/pull/6460) by @taliesin-ai
- Reduce allocations when generating name-based UUIDs in [PR](https://github.com/wailsapp/wails/pull/6455) by @taliesin-ai
- Reduce allocations when formatting SQLite statement handles in [PR](https://github.com/wailsapp/wails/pull/6454) by @taliesin-ai
- Reduce UUID string conversion to a single allocation in [PR](https://github.com/wailsapp/wails/pull/6453) by @taliesin-ai
- Reduce allocations when formatting changelog bullet entries in [PR](https://github.com/wailsapp/wails/pull/6449) by @taliesin-ai
- Reduce allocations when formatting terminal hyperlinks in [PR](https://github.com/wailsapp/wails/pull/6443) by @taliesin-ai
- Reduce allocations when serving assets by sniffing MIME buffers lazily in [PR](https://github.com/wailsapp/wails/pull/6428) by @taliesin-ai
- Skip debug log preparation in the asset server when debug logging is disabled in [PR](https://github.com/wailsapp/wails/pull/6420) by @taliesin-ai
- Reduce allocations in webview request handling by avoiding discarded request contexts in [PR](https://github.com/wailsapp/wails/pull/6419) by @taliesin-ai
- Speed up asset serving by skipping MIME sniffing for known file extensions in [PR](https://github.com/wailsapp/wails/pull/6414) by @taliesin-ai
- Reduce allocations when serving assets by avoiding a reader for content sniffing in [PR](https://github.com/wailsapp/wails/pull/6410) by @taliesin-ai
- Reduce allocations when serving non-seekable asset files in [PR](https://github.com/wailsapp/wails/pull/6408) by @taliesin-ai
- Speed up asset server index root selection by avoiding candidate collection in [PR](https://github.com/wailsapp/wails/pull/6404) by @taliesin-ai
- Improve asset server performance by formatting content lengths without fmt in [PR](https://github.com/wailsapp/wails/pull/6393) by @taliesin-ai
- Reduce memory copying when reading desktop file explorer metadata in [PR](https://github.com/wailsapp/wails/pull/6391) by @taliesin-ai

## Fixed
<!-- Bug fixes -->
- Keep shell command chains out of the implicit Go build cache in [PR](https://github.com/wailsapp/wails/pull/6493) by @taliesin-ai
- Fix notification category maps failing to load when saved data is null in [PR](https://github.com/wailsapp/wails/pull/6491) by @taliesin-ai
- Show diagnostics from failed variable shell commands in [PR](https://github.com/wailsapp/wails/pull/6492) by @taliesin-ai
- Fix wake template matching when multiple equality candidates exist in [PR](https://github.com/wailsapp/wails/pull/6489) by @taliesin-ai
- Serialize concurrent cache-file saves to prevent corrupted cache writes in [PR](https://github.com/wailsapp/wails/pull/6488) by @taliesin-ai
- Prevent races when replacing the active reporter in [PR](https://github.com/wailsapp/wails/pull/6487) by @taliesin-ai
- Report resource-generation failures accurately in the CLI in [PR](https://github.com/wailsapp/wails/pull/6486) by @taliesin-ai
- Recognize stable versions that include build metadata in the CLI in [PR](https://github.com/wailsapp/wails/pull/6484) by @taliesin-ai
- Serve assets whose request paths contain repeated leading slashes in [PR](https://github.com/wailsapp/wails/pull/6483) by @taliesin-ai
- Fix shell fallback execution for wake on each platform in [PR](https://github.com/wailsapp/wails/pull/6481) by @taliesin-ai
- Preserve previous Linux packages when a build fails in [PR](https://github.com/wailsapp/wails/pull/6480) by @taliesin-ai
- Route command chains and redirections through the shell in [PR](https://github.com/wailsapp/wails/pull/6479) by @taliesin-ai
- Preserve quoted regex and path backslashes in wake command parsing in [PR](https://github.com/wailsapp/wails/pull/6477) by @taliesin-ai
- Preserve existing DMG files when packaging builds fail in [PR](https://github.com/wailsapp/wails/pull/6476) by @taliesin-ai
- Replace symlinks with generated regular files in bindings output in [PR](https://github.com/wailsapp/wails/pull/6474) by @taliesin-ai
- Preserve literal percent signs in Linux explorer commands in [PR](https://github.com/wailsapp/wails/pull/6470) by @taliesin-ai
- Align precondition and task command directories so they resolve consistently in [PR](https://github.com/wailsapp/wails/pull/6471) by @taliesin-ai
- Run default tasks in the executor base directory in [PR](https://github.com/wailsapp/wails/pull/6469) by @taliesin-ai
- Prevent legacy checks from hiding cache misses in [PR](https://github.com/wailsapp/wails/pull/6468) by @taliesin-ai
- Fix CLI prerelease numbering to increment without integer overflow limits in [PR](https://github.com/wailsapp/wails/pull/6467) by @taliesin-ai
- Avoid a panic when cache entries decode as null in [PR](https://github.com/wailsapp/wails/pull/6466) by @taliesin-ai
- Probe ports once per check in the CLI in [PR](https://github.com/wailsapp/wails/pull/6465) by @taliesin-ai
- Include the task directory in cache identity to avoid stale wake results in [PR](https://github.com/wailsapp/wails/pull/6464) by @taliesin-ai
- Prevent data loss when packaging copies a file onto itself in [PR](https://github.com/wailsapp/wails/pull/6463) by @taliesin-ai
- Make status checks respect task directories in [PR](https://github.com/wailsapp/wails/pull/6462) by @taliesin-ai
- Detect TODO, FIXME and TBD placeholders in changelog entries in [PR](https://github.com/wailsapp/wails/pull/6459) by @taliesin-ai
- Generate runtime assets from embedded bundles in the CLI in [PR](https://github.com/wailsapp/wails/pull/6458) by @taliesin-ai
- Separate stderr from parsed HEAD output in git commands in [PR](https://github.com/wailsapp/wails/pull/6457) by @taliesin-ai
- Preserve the default in-memory SQLite database under concurrent access in [PR](https://github.com/wailsapp/wails/pull/6456) by @taliesin-ai
- Return correct responses to HEAD requests for the fallback index page in [PR](https://github.com/wailsapp/wails/pull/6451) by @taliesin-ai
- Handle SQLite statement preparation racing with database close or reopen in [PR](https://github.com/wailsapp/wails/pull/6450) by @taliesin-ai
- Release prepared-statement handles when SQLite database is closed or reopened in [PR](https://github.com/wailsapp/wails/pull/6448) by @taliesin-ai
- Keep extraction destination paths literal when using gosod in [PR](https://github.com/wailsapp/wails/pull/6447) by @taliesin-ai
- Propagate errors from SQLite queries after partial row iteration in [PR](https://github.com/wailsapp/wails/pull/6446) by @taliesin-ai
- Return snapshots for whole-store reads in kvstore in [PR](https://github.com/wailsapp/wails/pull/6444) by @taliesin-ai
- Prevent write panics in kvstore after loading a JSON null value in [PR](https://github.com/wailsapp/wails/pull/6442) by @taliesin-ai
- Print plain terminal messages without interpreting format verbs in [PR](https://github.com/wailsapp/wails/pull/6441) by @taliesin-ai
- Surface errors when looking up the latest stable release in the updater in [PR](https://github.com/wailsapp/wails/pull/6440) by @taliesin-ai
- Ignore hidden entries and headers in changelog comment blocks in [PR](https://github.com/wailsapp/wails/pull/6439) by @taliesin-ai
- Install wails3 once and surface installer errors in the updater in [PR](https://github.com/wailsapp/wails/pull/6438) by @taliesin-ai
- Preserve literal bracketed template paths in gosod in [PR](https://github.com/wailsapp/wails/pull/6437) by @taliesin-ai
- Reject failed GitHub tag responses before decoding them in the updater in [PR](https://github.com/wailsapp/wails/pull/6436) by @taliesin-ai
- Preserve cancellation errors during parallel library searches in [PR](https://github.com/wailsapp/wails/pull/6435) by @taliesin-ai
- Close HTTP response bodies on every return path to prevent resource leaks in [PR](https://github.com/wailsapp/wails/pull/6434) by @taliesin-ai
- Preserve request bodies when the development proxy retries requests in [PR](https://github.com/wailsapp/wails/pull/6432) by @taliesin-ai
- Honor the launch deadline during desktop file explorer discovery in [PR](https://github.com/wailsapp/wails/pull/6430) by @taliesin-ai
- Reject HTTP errors before saving AppImage tooling downloads in [PR](https://github.com/wailsapp/wails/pull/6429) by @taliesin-ai
- Order oversized numeric prerelease identifiers correctly in semver comparisons in [PR](https://github.com/wailsapp/wails/pull/6427) by @taliesin-ai
- Preserve escaped paths when routing services in the asset server in [PR](https://github.com/wailsapp/wails/pull/6426) by @taliesin-ai
- Handle HEAD requests for the bundled runtime in the asset server in [PR](https://github.com/wailsapp/wails/pull/6425) by @taliesin-ai
- Return 404 for missing bundled assets instead of an incorrect response in [PR](https://github.com/wailsapp/wails/pull/6424) by @taliesin-ai
- Stop dev proxy retries when a request is cancelled in [PR](https://github.com/wailsapp/wails/pull/6423) by @taliesin-ai
- Avoid shared start URL mutation and stale ports in the asset server in [PR](https://github.com/wailsapp/wails/pull/6422) by @taliesin-ai
- Honor the optional logger configuration in the asset server in [PR](https://github.com/wailsapp/wails/pull/6421) by @taliesin-ai
- Remove data race in library path cache invalidation in [PR](https://github.com/wailsapp/wails/pull/6418) by @taliesin-ai
- Avoid false library matches when searching ldconfig output in [PR](https://github.com/wailsapp/wails/pull/6416) by @taliesin-ai
- Prevent asset server queue growth and shrink oscillation under load in [PR](https://github.com/wailsapp/wails/pull/6415) by @taliesin-ai
- Release callback captures after the debounce timer fires to avoid retaining stale references in [PR](https://github.com/wailsapp/wails/pull/6413) by @taliesin-ai
- Preserve file explorer targets found inside executable paths in [PR](https://github.com/wailsapp/wails/pull/6412) by @taliesin-ai
- Skip directories that shadow desktop entries in file explorer in [PR](https://github.com/wailsapp/wails/pull/6411) by @taliesin-ai
- Fix asset server headers not being sent when a response body is empty in [PR](https://github.com/wailsapp/wails/pull/6409) by @taliesin-ai
- Respect XDG desktop-file search directories when locating application icons in [PR](https://github.com/wailsapp/wails/pull/6407) by @taliesin-ai
- Preserve file paths during field expansion in the file explorer in [PR](https://github.com/wailsapp/wails/pull/6406) by @taliesin-ai
- Avoid duplicating expanded URI targets in file explorer in [PR](https://github.com/wailsapp/wails/pull/6405) by @taliesin-ai
- Commit successful empty non-seekable asset responses instead of leaving them pending in [PR](https://github.com/wailsapp/wails/pull/6403) by @taliesin-ai
- Fill MIME sniffing prefix correctly across short reads in the asset server in [PR](https://github.com/wailsapp/wails/pull/6402) by @taliesin-ai
- Avoid reading full response bodies for non-seekable HEAD requests in the asset server in [PR](https://github.com/wailsapp/wails/pull/6401) by @taliesin-ai
- Only serve exact index files when resolving directory requests in the asset server in [PR](https://github.com/wailsapp/wails/pull/6400) by @taliesin-ai
- Fix dismissed Linux notifications being reported as clicks on servers without action support in [PR](https://github.com/wailsapp/wails/pull/6367) by @bercedev
- Stop dev processes when the terminal closes in [PR](https://github.com/wailsapp/wails/pull/6371) by @xmasdev
- Hold Windows taskbar restart debounce for one second before clearing the flag in [PR](https://github.com/wailsapp/wails/pull/6386) by @taliesin-ai
- Fix fallback response headers being lost in the asset server in [PR](https://github.com/wailsapp/wails/pull/6399) by @taliesin-ai
- Honor header deletions before fallback selection in the asset server in [PR](https://github.com/wailsapp/wails/pull/6398) by @taliesin-ai
- Release removed ring queue references in the asset server to avoid memory leaks in [PR](https://github.com/wailsapp/wails/pull/6397) by @taliesin-ai
- Release consumed content-sniffer buffers in the asset server to reduce memory use in [PR](https://github.com/wailsapp/wails/pull/6396) by @taliesin-ai
- Recognize mixed-case asset file extensions when serving assets in [PR](https://github.com/wailsapp/wails/pull/6395) by @taliesin-ai
- Correctly detect content type in asset server without retaining filenames in [PR](https://github.com/wailsapp/wails/pull/6394) by @taliesin-ai
- Fall back to a default when a desktop entry has no command in [PR](https://github.com/wailsapp/wails/pull/6392) by @taliesin-ai
- Stop the liquid-glass example from showing a dialog before the app starts on non-macOS in [PR](https://github.com/wailsapp/wails/pull/6389) by @taliesin-ai
- Keep window size constraints when fullscreen starts from a constrained size on Windows in [PR](https://github.com/wailsapp/wails/pull/6387) by @taliesin-ai
- Fix dialogs called before Run blocking forever instead of returning in [PR](https://github.com/wailsapp/wails/pull/6390) by @taliesin-ai

## Deprecated
<!-- Soon-to-be removed features -->

## Removed
<!-- Features removed in this release -->

## Security
<!-- Security-related changes -->
- Validate generated icon sizes before image allocation to prevent excessive memory use in [PR](https://github.com/wailsapp/wails/pull/6478) by @taliesin-ai
- Validate release notes API responses before use in [PR](https://github.com/wailsapp/wails/pull/6461) by @taliesin-ai
- Require complete hexadecimal statement handles in SQLite plugin in [PR](https://github.com/wailsapp/wails/pull/6452) by @taliesin-ai
- Constrain filesystem library fallback matches to prevent loading unintended libraries in [PR](https://github.com/wailsapp/wails/pull/6417) by @taliesin-ai

---
