# Native AppImage packaging

`CreateNative` collects ELF dependencies, adds an explicit GTK3 or GTK4 WebKit
profile for GUI applications (including purego's dynamically loaded libraries),
and writes a SquashFS image using `internal/readonlyfs`. It executes no target
programs or packaging tools. The older `Create` adapter retains its linuxdeploy
boundary for callers outside Wake.

`RuntimeRoot` is a Linux root filesystem tree, such as a directory containing
`usr/lib64` and `usr/share`. An explicit tree is resolved without searching the
packaging host. Absolute Linux symlinks stay inside that tree. Native Linux can
use the installed runtime; other hosts need a supplied tree for GUI applications.
The tree must include compiled GLib schemas and pre-generated GDK/GTK module
caches when it contains external loader modules. Target glibc and its ELF
interpreter remain operating-system prerequisites. Packaging a newer glibc ABI
does not make it compatible with older Linux distributions.

Staging materializes regular files, so Windows needs no elevated symlink
privileges. Archive modes are explicit Linux modes. `RuntimeIdentity` hashes the
selected destination/source mappings, source modes and contents, launcher, target
profile, and runtime artifact. Call it after binary inputs become available.

The external AppImage type-2 runtime is downloaded by immutable release-asset ID
and verified against pinned SHA256 values in `runtime_asset.go`. A local
`RuntimePath` override is target-validated and included in the content identity.
This runtime remains an external dependency; the implementation is not wholly
first-party. Runtime release:
https://github.com/AppImage/type2-runtime/commit/fedfad51acc2860b2916a1eef32ba4a500ff31ef

The new implementation was written from format documentation and existing
Wails-owned deployment code, without copying third-party implementations.
Update metadata and zsync sidecars follow the published formats:
https://github.com/AppImage/AppImageSpec/blob/master/draft.md
https://github.com/cph6/zsync/blob/master/SPEC.md
The existing `golang.org/x/crypto/md4` implementation supplies the legacy zsync
block hash; no Go module or version was added. Tests retain official zsyncmake
reference bytes and use the standard library. Benchmarks retain measurements of
runtime dependency collection and update-control generation.
